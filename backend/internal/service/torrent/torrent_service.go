package torrent

import (
	"context"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/ports"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"mime/multipart"

	"github.com/anacrolix/torrent/metainfo"
)

type Service struct {
	engine ports.TorrentEngine
	mapper *torrent.MetainfoMapper
}

func NewService(engine ports.TorrentEngine, mapper *torrent.MetainfoMapper) *Service {
	return &Service{engine: engine, mapper: mapper}
}

func (s *Service) Inspect(_ context.Context, file multipart.File) (torrentModel.Model, error) {
	mi, err := metainfo.Load(file)
	if err != nil {
		return torrentModel.Model{}, err
	}
	return s.mapper.ToModel(mi)
}

func (s *Service) Download(ctx context.Context, file multipart.File) (ports.DownloadResponse, error) {
	mi, err := metainfo.Load(file)
	if err != nil {
		return ports.DownloadResponse{}, err
	}

	return s.startEngine(ctx, mi.InfoBytes)
}

func (s *Service) startEngine(ctx context.Context, infoBytes []byte) (ports.DownloadResponse, error) {
	infoHash, err := s.engine.StartDownload(ctx, infoBytes)
	if err != nil {
		return ports.DownloadResponse{}, err
	}
	return ports.DownloadResponse{InfoHash: infoHash}, nil
}
