package torrent_signature

import (
	"context"
	"database/sql"
	"errors"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/torrent_signature_repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/tx"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	repository "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	"github.com/go-sql-driver/mysql"

	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
)

type TorrentSignatureRepository struct {
	db      *sql.DB
	queries *repository.Queries
}

func NewTorrentSignatureRepository(db *sql.DB, queries *generated.Queries) *TorrentSignatureRepository {
	return &TorrentSignatureRepository{db: db, queries: queries}
}

func (r *TorrentSignatureRepository) Create(ctx context.Context, entity torrentSignatureModel.TorrentSignatureEntity) error {
	params := r.mapToCreateParams(entity)
	_, err := r.queries.CreateTorrentSignature(ctx, params)
	return r.translateCreateError(err)
}

func (r *TorrentSignatureRepository) CreateInTransaction(ctx context.Context, entity torrentSignatureModel.TorrentSignatureEntity) error {
	sqlTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	rollbacker := tx.NewRollbacker(sqlTx)
	queries := r.queries.WithTx(sqlTx)
	return r.executeCreateInTx(ctx, queries, rollbacker, entity)
}

func (r *TorrentSignatureRepository) executeCreateInTx(ctx context.Context, queries *generated.Queries, rollbacker *tx.Rollbacker, entity torrentSignatureModel.TorrentSignatureEntity) error {
	params := r.mapToCreateParams(entity)
	_, err := queries.CreateTorrentSignature(ctx, params)
	if err != nil {
		return rollbacker.Rollback(r.translateCreateError(err))
	}
	return r.commitTx(rollbacker)
}

func (r *TorrentSignatureRepository) commitTx(rollbacker *tx.Rollbacker) error {
	err := rollbacker.Tx.Commit()
	if err != nil {
		return rollbacker.Rollback(err)
	}
	return nil
}

func (r *TorrentSignatureRepository) GetByTorrentHash(ctx context.Context, torrentHash []byte) ([]torrentSignatureModel.TorrentSignatureEntity, error) {
	rows, err := r.queries.GetTorrentSignaturesByHash(ctx, torrentHash)
	return r.mapRowsToEntities(rows), err
}

func (r *TorrentSignatureRepository) GetByUserAndTorrentHash(ctx context.Context, userID int64, torrentHash []byte) (*torrentSignatureModel.TorrentSignatureEntity, error) {
	params := generated.GetTorrentSignatureByUserAndHashParams{TorrentHash: torrentHash, UserID: userID}
	row, err := r.queries.GetTorrentSignatureByUserAndHash(ctx, params)
	return r.mapRowToEntityPtr(row), err
}

func (r *TorrentSignatureRepository) DeleteByTorrentHash(ctx context.Context, torrentHash []byte) error {
	_, err := r.queries.DeleteTorrentSignaturesByHash(ctx, torrentHash)
	return err
}

func (r *TorrentSignatureRepository) mapToCreateParams(entity torrentSignatureModel.TorrentSignatureEntity) generated.CreateTorrentSignatureParams {
	return generated.CreateTorrentSignatureParams{
		TorrentHash:   entity.TorrentHash,
		UserID:        entity.UserID,
		SignatureBlob: entity.SignatureBlob,
		PayloadHash:   entity.PayloadHash,
	}
}

func (r *TorrentSignatureRepository) translateCreateError(err error) error {
	if r.isDuplicateKeyError(err) {
		return repositoryErrors.NewDuplicateSignatureError()
	}
	return err
}

func (r *TorrentSignatureRepository) isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return false
}

func (r *TorrentSignatureRepository) mapRowsToEntities(rows []generated.TorrentSignature) []torrentSignatureModel.TorrentSignatureEntity {
	entities := make([]torrentSignatureModel.TorrentSignatureEntity, len(rows))
	for i, row := range rows {
		entities[i] = r.mapRowToEntity(row)
	}
	return entities
}

func (r *TorrentSignatureRepository) mapRowToEntity(row generated.TorrentSignature) torrentSignatureModel.TorrentSignatureEntity {
	return torrentSignatureModel.TorrentSignatureEntity{
		ID:            row.ID,
		TorrentHash:   row.TorrentHash,
		UserID:        row.UserID,
		SignatureBlob: row.SignatureBlob,
		PayloadHash:   row.PayloadHash,
		CreatedAt:     row.CreatedAt,
	}
}

func (r *TorrentSignatureRepository) mapRowToEntityPtr(row generated.TorrentSignature) *torrentSignatureModel.TorrentSignatureEntity {
	entity := r.mapRowToEntity(row)
	return &entity
}
