package entity

import (
	"encoding/base64"
	"encoding/hex"

	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

func ToEntity(row generated.Torrent) torrentModel.TorrentEntity {
	return torrentModel.TorrentEntity{
		InfoHash:    row.InfoHash,
		InfoBytes:   row.InfoBytes,
		Name:        row.Name,
		SizeBytes:   row.SizeBytes,
		PieceLength: int(row.PieceLength),
		IsPrivate:   row.IsPrivate,
		AddedAt:     row.AddedAt,
	}
}

func ToDTO(
	row generated.GetUserTorrentsRow,
	files []torrentModel.FileDTO,
	signatures []torrentModel.SignatureDTO,
) torrentModel.TorrentDTO {
	return torrentModel.TorrentDTO{
		InfoHash:         hex.EncodeToString(row.InfoHash),
		CreatorPublicKey: base64.RawURLEncoding.EncodeToString(row.CreatorPublicKey),
		Name:             row.Name,
		SizeBytes:        row.SizeBytes,
		StoragePath:      "",
		Status:           torrentModel.TorrentStatus(row.Status),
		Progress:         float32(row.Progress),
		AddedAt:          row.AddedAt,
		Files:            files,
		Signatures:       signatures,
	}
}

func ToStoreDTO(row generated.ListAllTorrentsRow) torrentModel.StoreTorrentDTO {
	return torrentModel.StoreTorrentDTO{
		InfoHash:         hex.EncodeToString(row.InfoHash),
		CreatorPublicKey: base64.RawURLEncoding.EncodeToString(row.CreatorPublicKey),
		Name:             row.Name,
		SizeBytes:        row.SizeBytes,
		AddedAt:          row.AddedAt,
	}
}
