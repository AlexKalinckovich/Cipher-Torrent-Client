package user

import (
	"encoding/base64"

	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

type UserDTOMapper struct{}

func NewUserDTOMapper() *UserDTOMapper {
	return &UserDTOMapper{}
}

func (m *UserDTOMapper) ToUserFullDTO(u db.User, s db.UserStat) userModel.UserFull {
	base := m.toUserDTO(u)
	stats := m.toUserStatsDTO(s)
	return userModel.UserFull{
		User:  base,
		Role:  string(u.Role),
		Stats: stats,
	}
}

func (m *UserDTOMapper) toUserDTO(src db.User) userModel.User {
	encodedToStringPublicKey := base64.StdEncoding.EncodeToString(src.PublicKey)
	return userModel.User{
		Id:        src.ID,
		Email:     src.Email,
		PublicKey: encodedToStringPublicKey,
		Nickname:  src.Nickname,
		CreatedAt: src.CreatedAt,
	}
}

func (m *UserDTOMapper) toUserStatsDTO(src db.UserStat) userModel.UserStats {
	return userModel.UserStats{
		TotalUploadedBytes:   src.TotalUploadedBytes,
		TotalDownloadedBytes: src.TotalDownloadedBytes,
		ReputationScore:      float32(src.ReputationScore),
		SignedTorrentsCount:  src.SignedTorrentsCount,
		ActiveTorrentsCount:  src.ActiveTorrentsCount,
		PeersTrustedCount:    src.PeersTrustedCount,
	}
}
