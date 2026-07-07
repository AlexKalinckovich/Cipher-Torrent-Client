package torrent

import (
	"path/filepath"
	"time"

	"github.com/anacrolix/torrent/metainfo"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type MetainfoMapper struct{}

func (m *MetainfoMapper) ToModel(mi *metainfo.MetaInfo) (torrentModel.Model, error) {
	info, err := mi.UnmarshalInfo()
	return m.handleInfoUnmarshal(mi, &info, err)
}

func (m *MetainfoMapper) handleInfoUnmarshal(mi *metainfo.MetaInfo, info *metainfo.Info, err error) (torrentModel.Model, error) {
	if err != nil {
		return torrentModel.Model{}, err
	}
	return m.buildModel(mi, info), nil
}

func (m *MetainfoMapper) buildModel(mi *metainfo.MetaInfo, info *metainfo.Info) torrentModel.Model {
	return torrentModel.Model{
		InfoHash:    mi.HashInfoBytes().HexString(),
		Name:        info.BestName(),
		SizeBytes:   m.totalSize(info),
		PieceLength: info.PieceLength,
		Files:       m.mapFiles(info),
		Status:      torrentModel.StatusPending,
		AddedAt:     time.Now(),
	}
}

func (m *MetainfoMapper) totalSize(info *metainfo.Info) int64 {
	var total int64
	for _, f := range info.UpvertedFiles() {
		total += f.Length
	}
	return total
}

func (m *MetainfoMapper) mapFiles(info *metainfo.Info) []torrentModel.File {
	upverted := info.UpvertedFiles()
	files := make([]torrentModel.File, len(upverted))
	for i, f := range upverted {
		files[i] = torrentModel.File{
			Path:      filepath.Join(f.Path...),
			SizeBytes: f.Length,
		}
	}
	return files
}
