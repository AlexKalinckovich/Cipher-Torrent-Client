package stats

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"log"
	"sync"
	"time"

	repository_port "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
)

// ProgressRepository persists a torrent's download progress for a user so it
// survives pause/resume and reloads.
type ProgressRepository interface {
	UpdateUserTorrentProgress(ctx context.Context, req repository_port.UpdateProgressRepositoryRequest) error
}

type torrentState struct {
	infoHashHex       string
	creatorPubKeyB64  string
	userID            int64
	lastDownloadedBytes int64
	lastUploadedBytes   int64
	lastUpdateTime      time.Time
	lastProgress        float32
}

type userDelta struct {
	downloaded int64
	uploaded   int64
}

type StatsTracker struct {
	broker            ports.EventBroker
	userRepo          *user.UserRepository
	progressRepo      ProgressRepository
	torrentMap        map[string]int64
	torrentStates     map[string]*torrentState
	mu                sync.Mutex
}

func NewStatsTracker(broker ports.EventBroker, userRepo *user.UserRepository, progressRepo ProgressRepository) *StatsTracker {
	tracker := &StatsTracker{
		broker:            broker,
		userRepo:          userRepo,
		progressRepo:      progressRepo,
		torrentMap:        make(map[string]int64),
		torrentStates:     make(map[string]*torrentState),
	}
	go tracker.listenToProgress()
	go tracker.periodicStatsUpdate()
	return tracker
}

func (t *StatsTracker) listenToProgress() {
	events, err := t.broker.Subscribe(context.Background(), "global_progress")
	if err != nil {
		log.Printf("[STATS] Failed to subscribe to global_progress: %v", err)
		return
	}
	log.Printf("[STATS] 🎧 Subscribed to global_progress")
	for evt := range events {
		log.Printf("[STATS] 📥 Received raw event id=%s size=%d", evt.ID, len(evt.Payload))
		t.processProgressEvent(string(evt.Payload))
	}
	log.Printf("[STATS] ⚠️ global_progress event channel closed")
}

func (t *StatsTracker) RegisterTorrent(infoHash, creatorPubKey string, userID int64) {
	key := t.makeKey(infoHash, creatorPubKey)
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.torrentStates[key]; exists {
		log.Printf("[STATS] Torrent %s already registered, skipping re-init", key)
		return
	}
	t.torrentMap[key] = userID
	t.torrentStates[key] = &torrentState{
		infoHashHex:      infoHash,
		creatorPubKeyB64: creatorPubKey,
		userID:           userID,
		lastUpdateTime:   time.Now(),
	}
	log.Printf("[STATS] Registered torrent: %s for user %d", key, userID)
}

func (t *StatsTracker) UnregisterTorrent(infoHash, creatorPubKey string) {
	key := t.makeKey(infoHash, creatorPubKey)
	state := t.extractAndRemoveState(key)
	if state == nil {
		return
	}
	t.flushSingleTorrentIfHasDelta(state)
}

func (t *StatsTracker) extractAndRemoveState(key string) *torrentState {
	t.mu.Lock()
	defer t.mu.Unlock()
	state, exists := t.torrentStates[key]
	if !exists {
		return nil
	}
	delete(t.torrentStates, key)
	delete(t.torrentMap, key)
	log.Printf("[STATS] Unregistered torrent: %s", key)
	return state
}

func (t *StatsTracker) flushSingleTorrentIfHasDelta(state *torrentState) {
	if state.lastDownloadedBytes == 0 && state.lastUploadedBytes == 0 {
		return
	}
	t.flushSingleTorrent(state.userID, state.lastDownloadedBytes, state.lastUploadedBytes)
}

func (t *StatsTracker) processProgressEvent(payload string) {
	event, err := t.parseEvent(payload)
	if err != nil || event.Progress == nil {
		return
	}
	t.applyProgressDelta(event.Progress)
}

func (t *StatsTracker) parseEvent(payload string) (*packet.Event, error) {
	var event packet.Event
	if err := json.Unmarshal([]byte(payload), &event); err != nil {
		log.Printf("[STATS] Failed to unmarshal event: %v", err)
		return nil, err
	}
	return &event, nil
}

func (t *StatsTracker) applyProgressDelta(progress *packet.ProgressLog) {
	key := t.makeKey(progress.InfoHash, progress.CreatorPublicKey)
	now := time.Now()

	t.mu.Lock()
	defer t.mu.Unlock()

	state, exists := t.torrentStates[key]
	if !exists {
		log.Printf("[STATS] ⚠️ no tracked state for key=%q (known keys: %v)", key, t.knownKeysLocked())
		return
	}

	elapsed := now.Sub(state.lastUpdateTime).Seconds()
	if elapsed <= 0 {
		return
	}

	state.lastDownloadedBytes += progress.DownloadSpeedBps * int64(elapsed)
	state.lastUploadedBytes += progress.UploadSpeedBps * int64(elapsed)
	state.lastUpdateTime = now
	// Track the latest completed fraction so we can persist it on the periodic flush.
	state.lastProgress = progress.Progress
	log.Printf("[STATS] ✅ applied delta key=%q downloaded+=%d uploaded+=%d progress=%.4f", key, progress.DownloadSpeedBps*int64(elapsed), progress.UploadSpeedBps*int64(elapsed), progress.Progress)
}

func (t *StatsTracker) periodicStatsUpdate() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		t.flushStatsToDB()
	}
}

func (t *StatsTracker) knownKeysLocked() []string {
	keys := make([]string, 0, len(t.torrentStates))
	for k := range t.torrentStates {
		keys = append(keys, k)
	}
	return keys
}

func (t *StatsTracker) flushStatsToDB() {
	userDeltas, progressUpdates := t.extractAndResetDeltas()
	t.persistDeltasToDB(userDeltas)
	t.persistProgressUpdates(progressUpdates)
}

// extractAndResetDeltas snapshots the byte deltas and the latest progress for
// every tracked torrent, then resets the byte counters. Progress is kept in
// state (not reset) and read from a separate, locked snapshot.
func (t *StatsTracker) extractAndResetDeltas() (map[int64]userDelta, []repository_port.UpdateProgressRepositoryRequest) {
	t.mu.Lock()
	defer t.mu.Unlock()
	userDeltas := make(map[int64]userDelta)
	progressUpdates := make([]repository_port.UpdateProgressRepositoryRequest, 0, len(t.torrentStates))
	for _, state := range t.torrentStates {
		t.accumulateDelta(userDeltas, state)
		if state.lastProgress > 0 {
			progressUpdates = append(progressUpdates, repository_port.UpdateProgressRepositoryRequest{
				UserID:        state.userID,
				InfoHash:      hexDecode(state.infoHashHex),
				CreatorPubKey: base64Decode(state.creatorPubKeyB64),
				// The engine reports progress as a 0..1 fraction; the DB and the
				// Dashboard render it as a percentage, so convert here.
				Progress: state.lastProgress * 100,
			})
		}
	}
	return userDeltas, progressUpdates
}

func (t *StatsTracker) accumulateDelta(userDeltas map[int64]userDelta, state *torrentState) {
	if state.lastDownloadedBytes == 0 && state.lastUploadedBytes == 0 {
		return
	}
	delta := userDeltas[state.userID]
	delta.downloaded += state.lastDownloadedBytes
	delta.uploaded += state.lastUploadedBytes
	userDeltas[state.userID] = delta
	state.lastDownloadedBytes = 0
	state.lastUploadedBytes = 0
}

// persistProgressUpdates writes the latest known progress fraction to the DB
// so a torrent keeps its progress across pause/resume and page reloads.
func (t *StatsTracker) persistProgressUpdates(updates []repository_port.UpdateProgressRepositoryRequest) {
	ctx := context.Background()
	for _, update := range updates {
		if err := t.progressRepo.UpdateUserTorrentProgress(ctx, update); err != nil {
			log.Printf("[STATS] Failed to persist progress for user %d torrent %s: %v", update.UserID, hex.EncodeToString(update.InfoHash), err)
			continue
		}
		log.Printf("[STATS] ✅ persisted progress user=%d torrent=%s progress=%.4f", update.UserID, hex.EncodeToString(update.InfoHash), update.Progress)
	}
}

func hexDecode(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}

func base64Decode(s string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}

func (t *StatsTracker) persistDeltasToDB(userDeltas map[int64]userDelta) {
	ctx := context.Background()
	for userID, delta := range userDeltas {
		t.updateUserStatsInDB(ctx, userID, delta)
	}
}

func (t *StatsTracker) updateUserStatsInDB(ctx context.Context, userID int64, delta userDelta) {
	currentStats, err := t.userRepo.GetStatsByUserID(ctx, userID)
	if t.handleStatsError(err, userID, "get") {
		return
	}
	params := t.buildUpdateStatsParams(userID, currentStats, delta)
	if err := t.userRepo.UpdateStats(ctx, params); err != nil {
		log.Printf("[STATS] Failed to update stats for user %d: %v", userID, err)
		return
	}
	log.Printf("[STATS] Successfully updated stats for user %d: downloaded=%d, uploaded=%d", userID, delta.downloaded, delta.uploaded)
}

func (t *StatsTracker) flushSingleTorrent(userID int64, downloaded, uploaded int64) {
	ctx := context.Background()
	currentStats, err := t.userRepo.GetStatsByUserID(ctx, userID)
	if t.handleStatsError(err, userID, "get") {
		return
	}
	params := t.buildUpdateStatsParams(userID, currentStats, userDelta{downloaded: downloaded, uploaded: uploaded})
	if err := t.userRepo.UpdateStats(ctx, params); err != nil {
		log.Printf("[STATS] Failed to update stats for user %d: %v", userID, err)
		return
	}
	log.Printf("[STATS] Successfully flushed final stats for user %d: downloaded=%d, uploaded=%d", userID, downloaded, uploaded)
}

func (t *StatsTracker) handleStatsError(err error, userID int64, action string) bool {
	if err != nil {
		log.Printf("[STATS] Failed to %s current stats for user %d: %v", action, userID, err)
		return true
	}
	return false
}

func (t *StatsTracker) buildUpdateStatsParams(userID int64, currentStats generated.UserStat, delta userDelta) generated.UpdateUserStatsParams {
	return generated.UpdateUserStatsParams{
		UserID:               userID,
		TotalUploadedBytes:   currentStats.TotalUploadedBytes + delta.uploaded,
		TotalDownloadedBytes: currentStats.TotalDownloadedBytes + delta.downloaded,
		ReputationScore:      currentStats.ReputationScore,
		SignedTorrentsCount:  currentStats.SignedTorrentsCount,
		ActiveTorrentsCount:  currentStats.ActiveTorrentsCount,
		PeersTrustedCount:    currentStats.PeersTrustedCount,
	}
}

func (t *StatsTracker) makeKey(infoHash, creatorPubKey string) string {
	return infoHash + "|" + creatorPubKey
}
