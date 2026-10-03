package torrent

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/service_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/stats"
	"io"
	"log"
	"mime/multipart"

	"github.com/anacrolix/torrent/metainfo"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/meta_info"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	torrentErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/torrent_errors"
	storage_ports "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/storage/ports"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type Service struct {
	repository   repository_ports.TorrentRepositoryPort
	engine       service_ports.TorrentEngine
	storage      storage_ports.TorrentStoragePort
	userRepo     user.UserRepositoryPort
	mapper       *meta_info.MetainfoMapper
	statsTracker *stats.StatsTracker
}

func NewService(
	repository repository_ports.TorrentRepositoryPort,
	engine service_ports.TorrentEngine,
	storage storage_ports.TorrentStoragePort,
	userRepo user.UserRepositoryPort,
	mapper *meta_info.MetainfoMapper,
	statsTracker *stats.StatsTracker,
) *Service {
	return &Service{
		repository:   repository,
		engine:       engine,
		storage:      storage,
		userRepo:     userRepo,
		mapper:       mapper,
		statsTracker: statsTracker,
	}
}

func (s *Service) Inspect(file multipart.File) (torrentModel.TorrentEntity, error) {
	mi, err := metainfo.Load(file)
	if err != nil {
		return torrentModel.TorrentEntity{}, err
	}
	return s.mapper.ToEntity(mi)
}

func (s *Service) Create(ctx context.Context, req service_ports.CreateTorrentServiceRequest) (torrentModel.TorrentDTO, error) {

	fullFileBytes, err := io.ReadAll(req.File)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	mi, err := metainfo.Load(bytes.NewReader(fullFileBytes))
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	return s.processCreate(ctx, mi, req.UserID, fullFileBytes)
}

func (s *Service) processCreate(ctx context.Context, mi *metainfo.MetaInfo, userID int64, fullFileBytes []byte) (torrentModel.TorrentDTO, error) {
	entity, err := s.mapper.ToEntity(mi)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.uploadAndSave(ctx, mi, entity, userID, fullFileBytes)
}

func (s *Service) uploadAndSave(ctx context.Context, mi *metainfo.MetaInfo, entity torrentModel.TorrentEntity, userID int64, fullFileBytes []byte) (torrentModel.TorrentDTO, error) {
	pubKey, err := s.getUserPublicKey(ctx, userID)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.uploadToMinioAndPersist(ctx, mi, entity, userID, pubKey, fullFileBytes)
}

func (s *Service) uploadToMinioAndPersist(ctx context.Context, mi *metainfo.MetaInfo, entity torrentModel.TorrentEntity, userID int64, pubKey []byte, fullFileBytes []byte) (torrentModel.TorrentDTO, error) {
	storageUploadRequest := storage_ports.StorageUploadRequest{
		InfoHash:      entity.InfoHash,
		CreatorPubKey: pubKey,
		FileBytes:     fullFileBytes,
	}
	if err := s.storage.UploadBaseTorrent(ctx, storageUploadRequest); err != nil {
		log.Printf("Error uploading torrent to minio: %v", err)
		return torrentModel.TorrentDTO{}, err
	}
	return s.persistAndStart(ctx, mi, entity, userID, pubKey)
}

func (s *Service) persistAndStart(ctx context.Context, mi *metainfo.MetaInfo, entity torrentModel.TorrentEntity, userID int64, pubKey []byte) (torrentModel.TorrentDTO, error) {
	repoReq := repository_ports.CreateTorrentRepositoryRequest{
		Entity:        entity,
		CreatorPubKey: pubKey,
		CreatorUserID: userID,
	}
	if err := s.repository.CreateTorrent(ctx, repoReq); err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.startAndBuildDTO(mi, entity, userID, pubKey)
}

func (s *Service) startAndBuildDTO(mi *metainfo.MetaInfo, entity torrentModel.TorrentEntity, userID int64, pubKey []byte) (torrentModel.TorrentDTO, error) {
	if _, err := s.engine.StartDownload(mi, pubKey); err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	s.statsTracker.RegisterTorrent(hex.EncodeToString(entity.InfoHash), base64.RawURLEncoding.EncodeToString(pubKey), userID)

	return s.mapper.ToDTOWithCreator(entity, pubKey, torrentModel.StatusDownloading, 0.0), nil
}

func (s *Service) getUserPublicKey(ctx context.Context, userID int64) ([]byte, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.PublicKey, nil
}

func (s *Service) Add(ctx context.Context, req service_ports.AddTorrentServiceRequest) error {
	// The store DTO carries the ORIGINAL creator's key. The .torrent lives in
	// MinIO under that creator's key (uploaded when the creator did /create).
	storageIdentityRequest := storage_ports.StorageIdentityRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}

	fullFileBytes, err := s.storage.DownloadBaseTorrent(ctx, storageIdentityRequest)
	if err != nil {
		log.Printf("Error downloading torrent file from storage: %v", err)
		return err
	}

	mi, err := metainfo.Load(bytes.NewReader(fullFileBytes))
	if err != nil {
		log.Printf("Error parsing downloaded torrent file: %v", err)
		return err
	}

	// Start the engine download from peers using the original creator's key.
	if _, err := s.engine.StartDownload(mi, req.CreatorPubKey); err != nil {
		return err
	}

	s.statsTracker.RegisterTorrent(
		hex.EncodeToString(req.InfoHash),
		base64.RawURLEncoding.EncodeToString(req.CreatorPubKey),
		req.UserID,
	)

	// Re-upload the .torrent under the CURRENT user's key so that pause/resume
	// and delete can locate it in MinIO for this user.
	currentPubKey, err := s.getUserPublicKey(ctx, req.UserID)
	if err != nil {
		return err
	}
	uploadReq := storage_ports.StorageUploadRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: currentPubKey,
		FileBytes:     fullFileBytes,
	}
	if err := s.storage.UploadBaseTorrent(ctx, uploadReq); err != nil {
		log.Printf("Error uploading torrent to minio for user: %v", err)
		return err
	}

	repoReq := repository_ports.CreateUserTorrentRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
		Status:        torrentModel.StatusDownloading,
	}
	err = s.repository.CreateUserTorrent(ctx, repoReq)
	if s.isDuplicateUserTorrent(err) {
		// Already in the user's library -> idempotent no-op.
		return nil
	}
	return err
}

func (s *Service) isDuplicateUserTorrent(err error) bool {
	var userTorrentErr *repository_errors.DuplicateUserTorrentError
	if errors.As(err, &userTorrentErr) {
		return true
	}
	// The MySQL duplicate-entry fallthrough in the translator surfaces as
	// DuplicateTorrentError; treat that as an idempotent no-op too.
	var torrentErr *repository_errors.DuplicateTorrentError
	return errors.As(err, &torrentErr)
}

func (s *Service) GetByInfoHash(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) (torrentModel.TorrentEntity, error) {
	repoReq := repository_ports.TorrentIdentityRepositoryRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}
	return s.repository.GetTorrentByIdentity(ctx, repoReq)
}

func (s *Service) GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error) {
	return s.repository.GetUserTorrents(ctx, userID)
}

func (s *Service) GetStoreTorrents(ctx context.Context) ([]torrentModel.StoreTorrentDTO, error) {
	return s.repository.GetAllTorrents(ctx)
}

func (s *Service) PauseTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error {
	if err := s.engine.PauseTorrent(req.InfoHash); err != nil {
		return err
	}

	repoReq := repository_ports.UpdateStatusRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
		Status:        torrentModel.StatusPaused,
	}
	return s.repository.UpdateUserTorrentStatus(ctx, repoReq)
}

func (s *Service) ResumeTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error {
	storageIdentityRequest := storage_ports.StorageIdentityRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}

	fullFileBytes, err := s.storage.DownloadBaseTorrent(ctx, storageIdentityRequest)
	if err != nil {
		log.Printf("Error downloading torrent file from storage: %v", err)
		return err
	}

	mi, err := metainfo.Load(bytes.NewReader(fullFileBytes))
	if err != nil {
		log.Printf("Error parsing downloaded torrent file: %v", err)
		return err
	}

	if err := s.engine.ResumeTorrent(mi, req.CreatorPubKey); err != nil {
		return err
	}

	s.statsTracker.RegisterTorrent(
		hex.EncodeToString(req.InfoHash),
		base64.RawURLEncoding.EncodeToString(req.CreatorPubKey),
		req.UserID,
	)

	repoReq := repository_ports.UpdateStatusRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
		Status:        torrentModel.StatusDownloading,
	}
	return s.repository.UpdateUserTorrentStatus(ctx, repoReq)
}
func (s *Service) UpdateProgress(ctx context.Context, req service_ports.UpdateProgressServiceRequest) error {

	repoReq := repository_ports.UpdateProgressRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
		Progress:      req.Progress,
	}
	return s.repository.UpdateUserTorrentProgress(ctx, repoReq)
}

// DeleteTorrent removes a torrent from the CURRENT user's library only.
// It intentionally does NOT touch MinIO: the shared .torrent object is owned
// by the store and must remain available for other users' /add requests.
func (s *Service) DeleteTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error {
	if err := s.pauseTorrentIfRunning(req.InfoHash); err != nil {
		return err
	}

	s.statsTracker.UnregisterTorrent(hex.EncodeToString(req.InfoHash), base64.RawURLEncoding.EncodeToString(req.CreatorPubKey))

	repoReq := repository_ports.UserTorrentIdentityRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}
	return s.repository.DeleteUserTorrent(ctx, repoReq)
}

// DeletePublishedTorrent removes a torrent from the store entirely: the MinIO
// .torrent object, the torrents row (FK cascades to user_torrents and
// torrent_signature_map), and stops any active engine download.
func (s *Service) DeletePublishedTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error {
	if err := s.pauseTorrentIfRunning(req.InfoHash); err != nil {
		return err
	}

	s.statsTracker.UnregisterTorrent(hex.EncodeToString(req.InfoHash), base64.RawURLEncoding.EncodeToString(req.CreatorPubKey))

	storageIdentityRequest := storage_ports.StorageIdentityRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}
	if err := s.storage.DeleteTorrent(ctx, storageIdentityRequest); err != nil {
		return err
	}

	repoReq := repository_ports.TorrentIdentityRepositoryRequest{
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
	}
	return s.repository.DeleteTorrentByIdentity(ctx, repoReq)
}

func (s *Service) pauseTorrentIfRunning(infoHash []byte) error {
	err := s.engine.PauseTorrent(infoHash)
	if s.isTorrentNotFoundInClient(err) {
		return nil
	}
	return err
}

func (s *Service) isTorrentNotFoundInClient(err error) bool {
	var target *torrentErrors.TorrentNotFoundInClientError
	return errors.As(err, &target)
}
