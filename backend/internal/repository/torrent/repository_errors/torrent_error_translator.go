package repository_errors

import (
	"database/sql"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	notFound "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlDuplicateEntryCode   uint16 = 1062
	torrentPrimaryKeyName            = "torrents.PRIMARY"
	userTorrentPrimaryKeyName        = "user_torrents.PRIMARY"
)

type TorrentErrorTranslator interface {
	TranslateTorrentError(err error) error
}

type torrentErrorTranslator struct{}

var TorrentTranslator TorrentErrorTranslator = &torrentErrorTranslator{}

func (t *torrentErrorTranslator) TranslateTorrentError(err error) error {
	if err == nil {
		return nil
	}
	return t.resolveTorrentError(err)
}

func (t *torrentErrorTranslator) resolveTorrentError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return notFound.NewNotFoundError[torrentModel.TorrentEntity]()
	}
	return t.resolveMySQLError(err)
}

func (t *torrentErrorTranslator) resolveMySQLError(err error) error {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return t.classifyMySQLError(mysqlErr)
	}
	return t.resolveNetworkError(err)
}

func (t *torrentErrorTranslator) classifyMySQLError(err *mysql.MySQLError) error {
	if err.Number == mysqlDuplicateEntryCode {
		return t.classifyDuplicateEntry(err.Message)
	}
	return err
}

func (t *torrentErrorTranslator) classifyDuplicateEntry(message string) error {
	if strings.Contains(message, torrentPrimaryKeyName) {
		return &DuplicateTorrentError{}
	}
	if strings.Contains(message, userTorrentPrimaryKeyName) {
		return &DuplicateUserTorrentError{}
	}
	return &DuplicateTorrentError{}
}

func (t *torrentErrorTranslator) resolveNetworkError(err error) error {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return &repository_errors.DatabaseUnavailableError{}
	}
	return err
}
