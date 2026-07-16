package torrent

import (
	"bytes"
	"context"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/service_ports"
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

type creationPayload struct {
	mi     *metainfo.MetaInfo
	entity torrentModel.TorrentEntity
	userID int64
	pubKey []byte
}

type Service struct {
	repository repository_ports.TorrentRepositoryPort
	engine     service_ports.TorrentEngine
	storage    storage_ports.TorrentStoragePort
	userRepo   user.UserRepositoryPort
	mapper     *meta_info.MetainfoMapper
}

func NewService(
	repository repository_ports.TorrentRepositoryPort,
	engine service_ports.TorrentEngine,
	storage storage_ports.TorrentStoragePort,
	userRepo user.UserRepositoryPort,
	mapper *meta_info.MetainfoMapper,
) *Service {
	return &Service{
		repository: repository,
		engine:     engine,
		storage:    storage,
		userRepo:   userRepo,
		mapper:     mapper,
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
	// 1. Read the FULL file bytes from the multipart upload
	fullFileBytes, err := io.ReadAll(req.File)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}

	// 2. Parse the metainfo from the full bytes
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
		FileBytes:     fullFileBytes, // <--- SAVE THE FULL FILE, NOT JUST InfoBytes!
	}
	if err := s.storage.UploadBaseTorrent(ctx, storageUploadRequest); err != nil {
		log.Printf("Error uploading torrent to minio: %v", err)
		return torrentModel.TorrentDTO{}, err
	}
	return s.persistAndStart(ctx, entity, userID, pubKey)
}

func (s *Service) persistAndStart(ctx context.Context, entity torrentModel.TorrentEntity, userID int64, pubKey []byte) (torrentModel.TorrentDTO, error) {
	repoReq := repository_ports.CreateTorrentRepositoryRequest{
		Entity:        entity,
		CreatorPubKey: pubKey,
		CreatorUserID: userID,
	}
	if err := s.repository.CreateTorrent(ctx, repoReq); err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.startAndBuildDTO(entity)
}

func (s *Service) startAndBuildDTO(entity torrentModel.TorrentEntity) (torrentModel.TorrentDTO, error) {
	if _, err := s.engine.StartDownload(entity.InfoBytes); err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.mapper.ToDTO(entity, torrentModel.StatusDownloading, 0.0), nil
}

func (s *Service) getUserPublicKey(ctx context.Context, userID int64) ([]byte, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user.PublicKey, nil
}

func (s *Service) Add(ctx context.Context, req service_ports.AddTorrentServiceRequest) error {
	repoReq := repository_ports.CreateUserTorrentRepositoryRequest{
		UserID:        req.UserID,
		InfoHash:      req.InfoHash,
		CreatorPubKey: req.CreatorPubKey,
		Status:        torrentModel.StatusIdle,
	}
	return s.repository.CreateUserTorrent(ctx, repoReq)
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
	infoBytes, err := s.storage.DownloadBaseTorrent(ctx, storageIdentityRequest)
	if err != nil {
		log.Printf("Error downloading torrent info from storage: %v", err)
		return err
	}
	return s.resumeEngineAndStatus(ctx, infoBytes, req.UserID, req.InfoHash, req.CreatorPubKey)
}

func (s *Service) resumeEngineAndStatus(ctx context.Context, infoBytes []byte, userID int64, infoHash []byte, pubKey []byte) error {
	if err := s.engine.ResumeTorrent(infoBytes); err != nil {
		return err
	}

	repoReq := repository_ports.UpdateStatusRepositoryRequest{
		UserID:        userID,
		InfoHash:      infoHash,
		CreatorPubKey: pubKey,
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

func (s *Service) DeleteTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error {
	if err := s.pauseTorrentIfRunning(req.InfoHash); err != nil {
		return err
	}
	return s.deleteStorageAndDatabase(ctx, req.UserID, req.InfoHash, req.CreatorPubKey)
}

func (s *Service) pauseTorrentIfRunning(infoHash []byte) error {
	err := s.engine.PauseTorrent(infoHash)
	if s.isTorrentNotFoundInClient(err) {
		return nil
	}
	return err
}

func (s *Service) deleteStorageAndDatabase(ctx context.Context, userID int64, infoHash []byte, pubKey []byte) error {

	storageIdentityRequest := storage_ports.StorageIdentityRequest{
		InfoHash:      infoHash,
		CreatorPubKey: pubKey,
	}

	if err := s.storage.DeleteTorrent(ctx, storageIdentityRequest); err != nil {
		return err
	}

	repoReq := repository_ports.UserTorrentIdentityRepositoryRequest{
		UserID:        userID,
		InfoHash:      infoHash,
		CreatorPubKey: pubKey,
	}
	return s.repository.DeleteUserTorrent(ctx, repoReq)
}

func (s *Service) isTorrentNotFoundInClient(err error) bool {
	var target *torrentErrors.TorrentNotFoundInClientError
	return errors.As(err, &target)
}
