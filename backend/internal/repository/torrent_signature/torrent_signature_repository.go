package torrent_signature

import (
	"context"
	"database/sql"
	"errors"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/torrent_signature_repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature/torrent_signature_repository_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/tx"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
	"github.com/go-sql-driver/mysql"
)

type TorrentSignatureRepository struct {
	db      *sql.DB
	queries *generated.Queries
}

func NewTorrentSignatureRepository(db *sql.DB, queries *generated.Queries) *TorrentSignatureRepository {
	return &TorrentSignatureRepository{db: db, queries: queries}
}

func (r *TorrentSignatureRepository) CreateSignature(
	ctx context.Context,
	entity torrentSignatureModel.TorrentSignatureEntity,
) error {
	params := r.mapToCreateSignatureParams(entity)
	_, err := r.queries.CreateSignature(ctx, params)
	return r.translateCreateError(err)
}

func (r *TorrentSignatureRepository) CreateSignatureMap(
	ctx context.Context,
	req torrent_signature_repository_ports.CreateSignatureMapRequest,
) error {
	params := r.mapToCreateSignatureMapParams(req)
	_, err := r.queries.CreateSignatureMap(ctx, params)
	return r.translateCreateError(err)
}

func (r *TorrentSignatureRepository) CreateInTransaction(
	ctx context.Context,
	sigEntity torrentSignatureModel.TorrentSignatureEntity,
	mapReq torrent_signature_repository_ports.CreateSignatureMapRequest,
) error {
	sqlTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	rollbacker := tx.NewRollbacker(sqlTx)
	queries := r.queries.WithTx(sqlTx)
	return r.executeCreateInTx(ctx, queries, rollbacker, sigEntity, mapReq)
}

func (r *TorrentSignatureRepository) executeCreateInTx(
	ctx context.Context,
	queries *generated.Queries,
	rollbacker *tx.Rollbacker,
	sigEntity torrentSignatureModel.TorrentSignatureEntity,
	mapReq torrent_signature_repository_ports.CreateSignatureMapRequest,
) error {
	sigParams := r.mapToCreateSignatureParams(sigEntity)

	result, err := queries.CreateSignature(ctx, sigParams)
	if err != nil {
		return rollbacker.Rollback(r.translateCreateError(err))
	}

	sigID, err := result.LastInsertId()
	if err != nil {
		return rollbacker.Rollback(err)
	}

	mapReq.SignatureID = sigID

	mapParams := r.mapToCreateSignatureMapParams(mapReq)
	_, err = queries.CreateSignatureMap(ctx, mapParams)
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

func (r *TorrentSignatureRepository) GetByTorrentIdentity(ctx context.Context, req torrent_signature_repository_ports.TorrentIdentityRepositoryRequest) ([]torrentSignatureModel.TorrentSignatureEntity, error) {
	params := r.mapIdentityParams(req)
	rows, err := r.queries.GetSignaturesWithSignerKey(ctx, params)
	return r.mapRowsToEntities(rows), err
}

func (r *TorrentSignatureRepository) GetByUserAndTorrentIdentity(ctx context.Context, req torrent_signature_repository_ports.UserTorrentIdentityRepositoryRequest) (*torrentSignatureModel.TorrentSignatureEntity, error) {
	params := r.mapUserAndIdentityParams(req)
	row, err := r.queries.GetSignatureMapByTorrentAndSigner(ctx, params)
	return r.mapRowToEntityPtrFromMap(row), err
}

func (r *TorrentSignatureRepository) DeleteByTorrentIdentity(ctx context.Context, req torrent_signature_repository_ports.TorrentIdentityRepositoryRequest) error {
	params := r.mapIdentityParamsToDeleteParams(req)
	_, err := r.queries.DeleteSignatureMapByTorrent(ctx, params)
	return err
}

func (r *TorrentSignatureRepository) mapToCreateSignatureParams(entity torrentSignatureModel.TorrentSignatureEntity) generated.CreateSignatureParams {
	return generated.CreateSignatureParams{
		TorrentHash:   entity.TorrentHash,
		UserID:        entity.UserID,
		SignatureBlob: entity.SignatureBlob,
		PayloadHash:   entity.PayloadHash,
	}
}

func (r *TorrentSignatureRepository) mapToCreateSignatureMapParams(req torrent_signature_repository_ports.CreateSignatureMapRequest) generated.CreateSignatureMapParams {
	return generated.CreateSignatureMapParams{
		TorrentInfoHash:  req.TorrentHash,
		CreatorPublicKey: req.CreatorPubKey,
		SignatureID:      req.SignatureID,
		SignerID:         req.SignerID,
		TrustLevel:       req.TrustLevel,
	}
}

func (r *TorrentSignatureRepository) mapIdentityParamsToDeleteParams(req torrent_signature_repository_ports.TorrentIdentityRepositoryRequest) generated.DeleteSignatureMapByTorrentParams {
	return generated.DeleteSignatureMapByTorrentParams{
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}

func (r *TorrentSignatureRepository) mapIdentityParams(req torrent_signature_repository_ports.TorrentIdentityRepositoryRequest) generated.GetSignaturesWithSignerKeyParams {
	return generated.GetSignaturesWithSignerKeyParams{
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}

func (r *TorrentSignatureRepository) mapUserAndIdentityParams(req torrent_signature_repository_ports.UserTorrentIdentityRepositoryRequest) generated.GetSignatureMapByTorrentAndSignerParams {
	return generated.GetSignatureMapByTorrentAndSignerParams{
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
		SignerID:         req.UserID,
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

func (r *TorrentSignatureRepository) mapRowsToEntities(rows []generated.GetSignaturesWithSignerKeyRow) []torrentSignatureModel.TorrentSignatureEntity {
	entities := make([]torrentSignatureModel.TorrentSignatureEntity, len(rows))
	for i, row := range rows {
		entities[i] = r.mapRowToEntityFromKey(row)
	}
	return entities
}

func (r *TorrentSignatureRepository) mapRowToEntityFromKey(row generated.GetSignaturesWithSignerKeyRow) torrentSignatureModel.TorrentSignatureEntity {
	return torrentSignatureModel.TorrentSignatureEntity{
		ID:            row.SignatureID,
		TorrentHash:   row.TorrentHash,
		UserID:        row.SignerID,
		SignatureBlob: row.SignatureBlob,
		PayloadHash:   row.PayloadHash,
		CreatedAt:     row.CreatedAt,
	}
}

func (r *TorrentSignatureRepository) mapRowToEntityPtrFromMap(row generated.TorrentSignatureMap) *torrentSignatureModel.TorrentSignatureEntity {
	entity := torrentSignatureModel.TorrentSignatureEntity{
		TorrentHash: row.TorrentInfoHash,
		UserID:      row.SignerID,
	}
	return &entity
}
