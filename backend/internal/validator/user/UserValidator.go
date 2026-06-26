package user

import (
	"encoding/base64"
	"fmt"
	serviceUser "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	"regexp"
)

const neededPublicKeyLength = 32

var (
	emailRegex    = regexp.MustCompile(`^([\w._]{2,10})@(\w+)\.([a-z]{2,4})$`)
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,16}$`)
	allowedRoles  = map[string]struct{}{
		"user":  {},
		"admin": {},
	}
)

type UserValidator struct{}

func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

func (v *UserValidator) ValidateCreate(params serviceUser.CreateUserInput) error {
	agg := validation.NewAggregateError()
	v.checkEmailField(agg, params.Email)
	v.checkPublicKeyField(agg, params.PublicKey)
	v.checkNicknameField(agg, params.Nickname)
	v.checkRoleField(agg, params.Role)
	return v.evaluateAggregate(agg)
}

func (v *UserValidator) ValidateUpdate(params serviceUser.UpdateUserInput) error {
	agg := validation.NewAggregateError()
	v.checkEmailField(agg, params.Email)
	v.checkPublicKeyField(agg, params.PublicKey)
	v.checkNicknameField(agg, params.Nickname)
	v.checkRoleField(agg, params.Role)
	return v.evaluateAggregate(agg)
}

func (v *UserValidator) ValidatePatch(fields serviceUser.PatchUserFields) error {
	agg := validation.NewAggregateError()
	v.checkEmailPtrField(agg, fields.Email)
	v.checkPublicKeyPtrField(agg, fields.PublicKey)
	v.checkNicknamePtrField(agg, fields.Nickname)
	v.checkRolePtrField(agg, fields.Role)
	return v.evaluateAggregate(agg)
}

func (v *UserValidator) evaluateAggregate(agg *validation.AggregateError) error {
	if agg.HasErrors() {
		return agg
	}
	return nil
}

func (v *UserValidator) checkEmailField(agg *validation.AggregateError, email string) {
	if email == "" {
		agg.Add("email", email, "email cannot be empty")
		return
	}
	if !emailRegex.MatchString(email) {
		agg.Add("email", email, "invalid email format")
	}
}

func (v *UserValidator) checkPublicKeyField(agg *validation.AggregateError, publicKey string) {
	if publicKey == "" {
		agg.Add("public_key", publicKey, "public key cannot be empty")
		return
	}
	decoded, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		agg.Add("public_key", publicKey, "public key must be a valid base64 string")
		return
	}

	decodedLen := len(decoded)

	if decodedLen != neededPublicKeyLength {
		agg.Add("public_key", publicKey, fmt.Sprintf("invalid Ed25519 public key length, need %d got %d", neededPublicKeyLength, decodedLen))
	}
}

func (v *UserValidator) checkNicknameField(agg *validation.AggregateError, nickname string) {
	if !nicknameRegex.MatchString(nickname) {
		agg.Add("nickname", nickname, "nickname must be alphanumeric and 3-16 characters long")
	}
}

func (v *UserValidator) checkRoleField(agg *validation.AggregateError, role string) {
	if _, ok := allowedRoles[role]; !ok {
		agg.Add("role", role, "authorization node role assignment must be user or admin")
	}
}

func (v *UserValidator) checkEmailPtrField(agg *validation.AggregateError, email *string) {
	if email != nil {
		v.checkEmailField(agg, *email)
	}
}

func (v *UserValidator) checkPublicKeyPtrField(agg *validation.AggregateError, publicKey *string) {
	if publicKey != nil {
		v.checkPublicKeyField(agg, *publicKey)
	}
}

func (v *UserValidator) checkNicknamePtrField(agg *validation.AggregateError, nickname *string) {
	if nickname != nil {
		v.checkNicknameField(agg, *nickname)
	}
}

func (v *UserValidator) checkRolePtrField(agg *validation.AggregateError, role *string) {
	if role != nil {
		v.checkRoleField(agg, *role)
	}
}
