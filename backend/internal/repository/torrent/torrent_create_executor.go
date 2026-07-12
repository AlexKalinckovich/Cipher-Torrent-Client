package torrent

import (
	"context"
	"database/sql"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/tx"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	repository "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentCreateExecutor struct {
	ctx      context.Context
	database *sql.DB
	queries  *repository.Queries
}

func NewTorrentCreateExecutor(
	ctx context.Context,
	database *sql.DB,
	queries *generated.Queries,
) *TorrentCreateExecutor {
	return &TorrentCreateExecutor{
		ctx:      ctx,
		database: database,
		queries:  queries,
	}
}

func (e *TorrentCreateExecutor) Execute(entity torrentModel.TorrentEntity, userID int64) error {
	sqlTx, err := e.database.BeginTx(e.ctx, nil)
	return e.tryCreateTorrent(sqlTx, entity, userID, err)
}

func (e *TorrentCreateExecutor) tryCreateTorrent(sqlTx *sql.Tx, entity torrentModel.TorrentEntity, userID int64, err error) error {
	if err != nil {
		return err
	}
	qtx := e.queries.WithTx(sqlTx)
	roller := tx.NewRollbacker(sqlTx)
	return e.tryCreateUserTorrent(sqlTx, qtx, roller, entity, userID)
}

func (e *TorrentCreateExecutor) tryCreateUserTorrent(sqlTx *sql.Tx, qtx *generated.Queries, roller *tx.Rollbacker, entity torrentModel.TorrentEntity, userID int64) error {
	err := e.createTorrentIfNeeded(qtx, entity)
	if err != nil {
		return roller.Rollback(err)
	}
	return e.tryCommit(sqlTx, qtx, roller, entity, userID)
}

func (e *TorrentCreateExecutor) tryCommit(sqlTx *sql.Tx, qtx *generated.Queries, roller *tx.Rollbacker, entity torrentModel.TorrentEntity, userID int64) error {
	err := e.createUserTorrent(qtx, userID, entity.InfoHash)
	if err != nil {
		return roller.Rollback(err)
	}
	return sqlTx.Commit()
}

func (e *TorrentCreateExecutor) createTorrentIfNeeded(qtx *generated.Queries, entity torrentModel.TorrentEntity) error {
	_, err := qtx.GetTorrentByInfoHash(e.ctx, entity.InfoHash)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return e.insertTorrent(qtx, entity)
}

func (e *TorrentCreateExecutor) insertTorrent(qtx *generated.Queries, entity torrentModel.TorrentEntity) error {
	_, err := qtx.CreateTorrent(e.ctx, generated.CreateTorrentParams{
		InfoHash:    entity.InfoHash,
		InfoBytes:   entity.InfoBytes,
		Name:        entity.Name,
		SizeBytes:   entity.SizeBytes,
		PieceLength: int32(entity.PieceLength),
		IsPrivate:   entity.IsPrivate,
		StoragePath: entity.StoragePath,
		AddedAt:     entity.AddedAt,
	})
	return err
}

func (e *TorrentCreateExecutor) createUserTorrent(qtx *generated.Queries, userID int64, infoHash []byte) error {
	_, err := qtx.CreateUserTorrent(e.ctx, generated.CreateUserTorrentParams{
		UserID:          userID,
		TorrentInfoHash: infoHash,
		Status:          string(torrentModel.StatusIdle),
		Progress:        0.00,
	})
	return err
}
