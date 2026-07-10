package torrent

import (
	"context"
	"errors"
	"mime/multipart"

	"github.com/anacrolix/torrent/metainfo"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/meta_info"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/ports"
	torrentErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/torrent_errors"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type Service struct {
	repository ports.TorrentRepositoryPort
	engine     ports.TorrentEngine
	mapper     *meta_info.MetainfoMapper
	validator  ports.TorrentValidatorPort
}

func NewService(
	repository ports.TorrentRepositoryPort,
	engine ports.TorrentEngine,
	mapper *meta_info.MetainfoMapper,
	validator ports.TorrentValidatorPort,
) *Service {
	return &Service{
		repository: repository,
		engine:     engine,
		mapper:     mapper,
		validator:  validator,
	}
}

func (s *Service) Inspect(file multipart.File) (torrentModel.TorrentEntity, error) {
	mi, err := metainfo.Load(file)
	if err != nil {
		return torrentModel.TorrentEntity{}, err
	}
	return s.mapper.ToEntity(mi)
}

func (s *Service) Add(ctx context.Context, file multipart.File, userID int64, savePath string) (torrentModel.TorrentDTO, error) {
	mi, err := metainfo.Load(file)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	return s.processAndSave(ctx, mi, userID, savePath)
}

func (s *Service) processAndSave(ctx context.Context, mi *metainfo.MetaInfo, userID int64, savePath string) (torrentModel.TorrentDTO, error) {
	entity, err := s.mapper.ToEntity(mi)
	if err != nil {
		return torrentModel.TorrentDTO{}, err
	}
	entity.StoragePath = savePath
	return s.saveAndStart(ctx, entity, userID)
}

func (s *Service) saveAndStart(ctx context.Context, entity torrentModel.TorrentEntity, userID int64) (torrentModel.TorrentDTO, error) {
	if err := s.repository.AddTorrent(ctx, entity, userID); err != nil {
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

func (s *Service) GetByInfoHash(ctx context.Context, infoHash []byte) (torrentModel.TorrentEntity, error) {
	return s.repository.GetByInfoHash(ctx, infoHash)
}

func (s *Service) GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error) {
	return s.repository.GetUserTorrents(ctx, userID)
}

func (s *Service) PauseTorrent(ctx context.Context, userID int64, infoHash []byte) error {
	if err := s.engine.PauseTorrent(infoHash); err != nil {
		return err
	}
	return s.repository.UpdateStatus(ctx, userID, infoHash, torrentModel.StatusPaused)
}

func (s *Service) ResumeTorrent(ctx context.Context, userID int64, infoHash []byte) error {
	entity, err := s.repository.GetByInfoHash(ctx, infoHash)
	if err != nil {
		return err
	}
	return s.resumeEngineAndStatus(ctx, entity, userID)
}

func (s *Service) resumeEngineAndStatus(ctx context.Context, entity torrentModel.TorrentEntity, userID int64) error {
	if err := s.engine.ResumeTorrent(entity.InfoBytes); err != nil {
		return err
	}
	return s.repository.UpdateStatus(ctx, userID, entity.InfoHash, torrentModel.StatusDownloading)
}

func (s *Service) UpdateProgress(ctx context.Context, userID int64, infoHash []byte, progress float32) error {
	return s.repository.UpdateProgress(ctx, userID, infoHash, progress)
}

func (s *Service) DeleteTorrent(ctx context.Context, userID int64, infoHash []byte) error {
	if err := s.pauseTorrentIfRunning(infoHash); err != nil {
		return err
	}
	return s.repository.DeleteTorrent(ctx, infoHash)
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
