package meta_info

import (
	"encoding/hex"
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent/metainfo"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type MetainfoMapper struct{}

func NewMetainfoMapper() *MetainfoMapper {
	return &MetainfoMapper{}
}

func (m *MetainfoMapper) ToEntity(mi *metainfo.MetaInfo) (torrentModel.TorrentEntity, error) {
	info, err := mi.UnmarshalInfo()
	if err != nil {
		return torrentModel.TorrentEntity{}, err
	}
	return m.buildEntity(mi, &info), nil
}

func (m *MetainfoMapper) ToDTO(entity torrentModel.TorrentEntity, status torrentModel.TorrentStatus, progress float32) torrentModel.TorrentDTO {
	return torrentModel.TorrentDTO{
		InfoHash:    hex.EncodeToString(entity.InfoHash),
		Name:        entity.Name,
		SizeBytes:   entity.SizeBytes,
		StoragePath: entity.StoragePath,
		Status:      status,
		Progress:    progress,
		AddedAt:     entity.AddedAt,
	}
}

func (m *MetainfoMapper) buildEntity(mi *metainfo.MetaInfo, info *metainfo.Info) torrentModel.TorrentEntity {
	return torrentModel.TorrentEntity{
		InfoHash:    mi.HashInfoBytes().Bytes(),
		InfoBytes:   mi.InfoBytes,
		Name:        info.BestName(),
		SizeBytes:   m.calculateTotalSize(info),
		PieceLength: int(info.PieceLength),
		IsPrivate:   info.Private != nil && *info.Private,
		StoragePath: "",
		AddedAt:     time.Now(),
	}
}

func (m *MetainfoMapper) calculateTotalSize(info *metainfo.Info) int64 {
	var total int64
	for _, f := range info.UpvertedFiles() {
		total += f.Length
	}
	return total
}

func (m *MetainfoMapper) ExtractFilePaths(info *metainfo.Info) []string {
	files := info.UpvertedFiles()
	paths := make([]string, len(files))
	for i, f := range files {
		paths[i] = filepath.Join(f.Path...)
	}
	return paths
}
