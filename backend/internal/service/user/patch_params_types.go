package user

import (
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
)

type PatchUserFields struct {
	Email        *string
	Nickname     *string
	PasswordHash *string
	Role         *string
}

type PatchParamsBuilder struct {
	current db.User
	fields  PatchUserFields
}
