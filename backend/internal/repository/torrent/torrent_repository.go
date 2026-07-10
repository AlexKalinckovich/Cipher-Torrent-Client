package torrent

import (
	"context"
	"database/sql"
	torrentMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/entity"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_errors"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentRepository struct {
	database   *sql.DB
	queries    *generated.Queries
	translator repository_errors.TorrentErrorTranslator
}

func NewTorrentRepository(database *sql.DB, queries *generated.Queries) *TorrentRepository {
	return &TorrentRepository{
		database:   database,
		queries:    queries,
		translator: repository_errors.TorrentTranslator,
	}
}

func (r *TorrentRepository) AddTorrent(ctx context.Context, entity torrentModel.TorrentEntity, userID int64) error {
	executor := NewTorrentCreateExecutor(ctx, r.database, r.queries)
	err := executor.Execute(entity, userID)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) GetByInfoHash(ctx context.Context, infoHash []byte) (torrentModel.TorrentEntity, error) {
	res, err := r.queries.GetTorrentByInfoHash(ctx, infoHash)
	if translatedErr := r.translator.TranslateTorrentError(err); translatedErr != nil {
		return torrentModel.TorrentEntity{}, translatedErr
	}
	return torrentMapper.ToEntity(res), nil
}

func (r *TorrentRepository) GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error) {
	rows, err := r.queries.GetUserTorrents(ctx, userID)
	if translatedErr := r.translator.TranslateTorrentError(err); translatedErr != nil {
		return nil, translatedErr
	}
	dtos := make([]torrentModel.TorrentDTO, len(rows))
	for i, row := range rows {
		dtos[i] = torrentMapper.ToDTO(row)
	}
	return dtos, nil
}

func (r *TorrentRepository) UpdateStatus(ctx context.Context, userID int64, infoHash []byte, status torrentModel.TorrentStatus) error {
	err := r.queries.UpdateUserTorrentStatus(ctx, generated.UpdateUserTorrentStatusParams{
		UserID:          userID,
		TorrentInfoHash: infoHash,
		Status:          string(status),
	})
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) UpdateProgress(ctx context.Context, userID int64, infoHash []byte, progress float32) error {
	err := r.queries.UpdateUserTorrentProgress(ctx, generated.UpdateUserTorrentProgressParams{
		UserID:          userID,
		TorrentInfoHash: infoHash,
		Progress:        float64(progress),
	})
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) DeleteTorrent(ctx context.Context, infoHash []byte) error {
	_, err := r.queries.DeleteTorrent(ctx, infoHash)
	return err
}
