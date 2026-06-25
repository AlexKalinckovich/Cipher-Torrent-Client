package repository_errors

import (
	"database/sql"
	"errors"
	notFound "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
)

type UserErrorTranslator interface {
	TranslateUserError(err error) error
}

type repoErrorTranslator struct{}

var Translator UserErrorTranslator = &repoErrorTranslator{}

func (t *repoErrorTranslator) TranslateUserError(err error) error {
	if err == nil {
		return nil
	}
	return t.resolveUserError(err)
}

func (t *repoErrorTranslator) resolveUserError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFound.NewNotFoundError[UserEntity]()
	}
	return err
}
