package user

import (
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

type UserDTOMapper struct{}

func NewUserDTOMapper() *UserDTOMapper {
	return &UserDTOMapper{}
}

func (m *UserDTOMapper) ToUserFullDTO(u db.User, s db.UserStat) user.UserFull {
	base := m.toUserDTO(u)
	stats := m.toUserStatsDTO(s)
	return user.UserFull{
		Id:        base.Id,
		Email:     base.Email,
		PublicKey: base.PublicKey,
		Nickname:  base.Nickname,
		CreatedAt: base.CreatedAt,
		Role:      string(u.Role),
		Stats:     stats,
	}
}

func (m *UserDTOMapper) toUserDTO(src db.User) user.User {
	return user.User{
		Id:        src.ID,
		Email:     src.Email,
		PublicKey: src.PublicKey,
		Nickname:  src.Nickname,
		CreatedAt: src.CreatedAt,
	}
}

func (m *UserDTOMapper) toUserStatsDTO(src db.UserStat) user.UserStats {
	return user.UserStats{
		TotalUploadedBytes:   src.TotalUploadedBytes,
		TotalDownloadedBytes: src.TotalDownloadedBytes,
		ReputationScore:      float32(src.ReputationScore),
		SignedTorrentsCount:  src.SignedTorrentsCount,
		ActiveTorrentsCount:  src.ActiveTorrentsCount,
		PeersTrustedCount:    src.PeersTrustedCount,
	}
}
