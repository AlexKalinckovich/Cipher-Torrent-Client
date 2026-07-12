package repository_errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/service_errors"
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"

	notFound "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
)

const (
	mysqlDuplicateEntryCode uint16 = 1062
	emailKeyName                   = "users.uq_users_email"
	nicknameKeyName                = "users.uq_users_nickname"
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
	return t.resolveMySQLError(err)
}

func (t *repoErrorTranslator) resolveMySQLError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return t.classifyMySQLError(mysqlErr)
	}
	return t.resolveNetworkError(err)
}

func (t *repoErrorTranslator) classifyMySQLError(err *mysql.MySQLError) error {
	if err.Number == mysqlDuplicateEntryCode {
		return t.classifyDuplicateEntry(err.Message)
	}
	return err
}

func (t *repoErrorTranslator) classifyDuplicateEntry(message string) error {
	fmt.Println(message)
	if strings.Contains(message, emailKeyName) {
		return &DuplicateEmailError{}
	}
	if strings.Contains(message, nicknameKeyName) {
		return &DuplicateNicknameError{}
	}
	return &DuplicateEmailError{}
}

func (t *repoErrorTranslator) resolveNetworkError(err error) error {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return &service_errors.DatabaseUnavailableError{}
	}
	return err
}
