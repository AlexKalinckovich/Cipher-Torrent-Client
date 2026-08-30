package stats

import (
	"context"
	"encoding/json"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"log"
	"sync"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
)

type torrentState struct {
	userID              int64
	lastDownloadedBytes int64
	lastUploadedBytes   int64
	lastUpdateTime      time.Time
}

type userDelta struct {
	downloaded int64
	uploaded   int64
}

type StatsTracker struct {
	broker        ports.EventBroker
	userRepo      *user.UserRepository
	torrentMap    map[string]int64
	torrentStates map[string]*torrentState
	mu            sync.Mutex
}

func NewStatsTracker(broker ports.EventBroker, userRepo *user.UserRepository) *StatsTracker {
	tracker := &StatsTracker{
		broker:        broker,
		userRepo:      userRepo,
		torrentMap:    make(map[string]int64),
		torrentStates: make(map[string]*torrentState),
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
		userID:         userID,
		lastUpdateTime: time.Now(),
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
	log.Printf("[STATS] ✅ applied delta key=%q downloaded+=%d uploaded+=%d", key, progress.DownloadSpeedBps*int64(elapsed), progress.UploadSpeedBps*int64(elapsed))
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
	userDeltas := t.extractAndResetDeltas()
	t.persistDeltasToDB(userDeltas)
}

func (t *StatsTracker) extractAndResetDeltas() map[int64]userDelta {
	t.mu.Lock()
	defer t.mu.Unlock()
	userDeltas := make(map[int64]userDelta)
	for _, state := range t.torrentStates {
		t.accumulateDelta(userDeltas, state)
	}
	return userDeltas
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
