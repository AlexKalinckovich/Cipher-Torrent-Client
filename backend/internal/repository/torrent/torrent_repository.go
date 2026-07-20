package torrent

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	torrentMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/entity"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/meta_data_extractor"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_error_translator"
	repository_port "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent/repository_ports"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentRepository struct {
	database   *sql.DB
	queries    *generated.Queries
	translator repository_error_translator.TorrentErrorTranslator
	extractor  *meta_data_extractor.MetadataExtractor
}

func NewTorrentRepository(database *sql.DB, queries *generated.Queries) *TorrentRepository {
	return &TorrentRepository{
		database:   database,
		queries:    queries,
		translator: repository_error_translator.TorrentTranslator,
		extractor:  meta_data_extractor.NewMetadataExtractor(),
	}
}

func (r *TorrentRepository) CreateTorrent(ctx context.Context, req repository_port.CreateTorrentRepositoryRequest) error {
	params := r.mapCreateTorrentParams(req)
	_, err := r.queries.CreateTorrent(ctx, params)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) mapCreateTorrentParams(req repository_port.CreateTorrentRepositoryRequest) generated.CreateTorrentParams {
	return generated.CreateTorrentParams{
		InfoHash:         req.Entity.InfoHash,
		InfoBytes:        req.Entity.InfoBytes,
		CreatorPublicKey: req.CreatorPubKey,
		CreatorUserID:    req.CreatorUserID,
		Name:             req.Entity.Name,
		SizeBytes:        req.Entity.SizeBytes,
		PieceLength:      int32(req.Entity.PieceLength),
		IsPrivate:        req.Entity.IsPrivate,
	}
}

func (r *TorrentRepository) CreateUserTorrent(ctx context.Context, req repository_port.CreateUserTorrentRepositoryRequest) error {
	params := r.mapCreateUserTorrentParams(req)
	_, err := r.queries.CreateUserTorrent(ctx, params)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) mapCreateUserTorrentParams(req repository_port.CreateUserTorrentRepositoryRequest) generated.CreateUserTorrentParams {
	return generated.CreateUserTorrentParams{
		UserID:           req.UserID,
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
		Status:           string(req.Status),
	}
}

func (r *TorrentRepository) GetTorrentByIdentity(ctx context.Context, req repository_port.TorrentIdentityRepositoryRequest) (torrentModel.TorrentEntity, error) {
	params := r.mapTorrentIdentityParams(req)
	res, err := r.queries.GetTorrentByIdentity(ctx, params)
	if translatedErr := r.translator.TranslateTorrentError(err); translatedErr != nil {
		return torrentModel.TorrentEntity{}, translatedErr
	}
	return torrentMapper.ToEntity(res), nil
}

func (r *TorrentRepository) mapTorrentIdentityParams(req repository_port.TorrentIdentityRepositoryRequest) generated.GetTorrentByIdentityParams {
	return generated.GetTorrentByIdentityParams{
		InfoHash:         req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}

func (r *TorrentRepository) GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error) {
	rows, err := r.queries.GetUserTorrents(ctx, userID)
	if translatedErr := r.translator.TranslateTorrentError(err); translatedErr != nil {
		return nil, translatedErr
	}
	return r.aggregateRowsToDTOs(rows), nil
}

func (r *TorrentRepository) aggregateRowsToDTOs(rows []generated.GetUserTorrentsRow) []torrentModel.TorrentDTO {

	type torrentKey struct {
		infoHash      string
		creatorPubKey string
	}

	groupedTorrents := make(map[torrentKey]*torrentModel.TorrentDTO)
	var orderedKeys []torrentKey

	for _, row := range rows {
		key := torrentKey{
			infoHash:      string(row.InfoHash),
			creatorPubKey: string(row.CreatorPublicKey),
		}

		dto, exists := groupedTorrents[key]
		if !exists {

			files, _ := r.extractor.Extract(row.InfoBytes)

			newDTO := torrentMapper.ToDTO(row, files, []torrentModel.SignatureDTO{})
			dto = &newDTO

			groupedTorrents[key] = dto
			orderedKeys = append(orderedKeys, key)
		}

		if sig := r.mapRowToSignatureDTO(row); sig != nil {
			dto.Signatures = append(dto.Signatures, *sig)
		}
	}

	result := make([]torrentModel.TorrentDTO, 0, len(orderedKeys))
	for _, key := range orderedKeys {
		result = append(result, *groupedTorrents[key])
	}
	return result
}

func (r *TorrentRepository) mapRowToSignatureDTO(row generated.GetUserTorrentsRow) *torrentModel.SignatureDTO {
	if !row.SignatureID.Valid {
		return nil
	}

	var signerPubKey string
	if row.SignerPublicKey.Valid {
		signerPubKey = base64.RawURLEncoding.EncodeToString([]byte(row.SignerPublicKey.String))
	}

	var sigBytes string
	if row.SignatureBlob.Valid {
		sigBytes = hex.EncodeToString([]byte(row.SignatureBlob.String))
	}

	var timestamp int64
	if row.SignatureCreatedAt.Valid {
		timestamp = row.SignatureCreatedAt.Time.Unix()
	}
	return &torrentModel.SignatureDTO{
		SignerPublicKey: signerPubKey,
		SignatureBytes:  sigBytes,
		Timestamp:       timestamp,
	}
}

func (r *TorrentRepository) mapRowsToDTOs(rows []generated.GetUserTorrentsRow) []torrentModel.TorrentDTO {
	dtos := make([]torrentModel.TorrentDTO, len(rows))
	for i, row := range rows {
		dtos[i] = r.mapRowToDTO(row)
	}
	return dtos
}

func (r *TorrentRepository) mapRowToDTO(row generated.GetUserTorrentsRow) torrentModel.TorrentDTO {
	files, signatures := r.extractor.Extract(row.InfoBytes)
	return torrentMapper.ToDTO(row, files, signatures)
}

func (r *TorrentRepository) UpdateUserTorrentStatus(ctx context.Context, req repository_port.UpdateStatusRepositoryRequest) error {
	params := r.mapUpdateStatusParams(req)
	err := r.queries.UpdateUserTorrentStatus(ctx, params)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) mapUpdateStatusParams(req repository_port.UpdateStatusRepositoryRequest) generated.UpdateUserTorrentStatusParams {
	return generated.UpdateUserTorrentStatusParams{
		Status:           string(req.Status),
		UserID:           req.UserID,
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}

func (r *TorrentRepository) UpdateUserTorrentProgress(ctx context.Context, req repository_port.UpdateProgressRepositoryRequest) error {
	params := r.mapUpdateProgressParams(req)
	err := r.queries.UpdateUserTorrentProgress(ctx, params)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) mapUpdateProgressParams(req repository_port.UpdateProgressRepositoryRequest) generated.UpdateUserTorrentProgressParams {
	return generated.UpdateUserTorrentProgressParams{
		Progress:         float64(req.Progress),
		UserID:           req.UserID,
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}

func (r *TorrentRepository) DeleteUserTorrent(ctx context.Context, req repository_port.UserTorrentIdentityRepositoryRequest) error {
	params := r.mapUserTorrentIdentityParams(req)
	_, err := r.queries.DeleteUserTorrent(ctx, params)
	return r.translator.TranslateTorrentError(err)
}

func (r *TorrentRepository) mapUserTorrentIdentityParams(req repository_port.UserTorrentIdentityRepositoryRequest) generated.DeleteUserTorrentParams {
	return generated.DeleteUserTorrentParams{
		UserID:           req.UserID,
		TorrentInfoHash:  req.InfoHash,
		CreatorPublicKey: req.CreatorPubKey,
	}
}
