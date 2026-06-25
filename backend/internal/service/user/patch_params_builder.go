package user

import (
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
)

type PatchUserFields struct {
	Email     *string
	PublicKey *string
	Nickname  *string
	Role      *string
}

type PatchParamsBuilder struct {
	current db.User
	fields  PatchUserFields
}

func NewPatchParamsBuilder(current db.User, fields PatchUserFields) *PatchParamsBuilder {
	return &PatchParamsBuilder{
		current: current,
		fields:  fields,
	}
}

func (b *PatchParamsBuilder) Build() db.UpdateUserParams {
	return db.UpdateUserParams{
		ID:        b.current.ID,
		Email:     b.resolveString(b.fields.Email, b.current.Email),
		PublicKey: b.resolveString(b.fields.PublicKey, b.current.PublicKey),
		Nickname:  b.resolveString(b.fields.Nickname, b.current.Nickname),
		Role:      b.resolveRole(b.fields.Role, b.current.Role),
	}
}

func (b *PatchParamsBuilder) resolveString(val *string, fallback string) string {
	if val != nil {
		return *val
	}
	return fallback
}

func (b *PatchParamsBuilder) resolveRole(val *string, fallback db.UsersRole) db.UsersRole {
	if val != nil {
		return db.UsersRole(*val)
	}
	return fallback
}
