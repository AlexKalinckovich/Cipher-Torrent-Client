package tx

import (
	"database/sql"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/repository_errors"
)

type Rollbacker struct {
	Tx *sql.Tx
}

func NewRollbacker(tx *sql.Tx) *Rollbacker {
	return &Rollbacker{Tx: tx}
}

func (r *Rollbacker) Rollback(originalErr error) error {

	rollbackErr := r.Tx.Rollback()

	return r.resolveRollbackError(originalErr, rollbackErr)
}

func (r *Rollbacker) resolveRollbackError(originalErr error, rollbackErr error) error {
	if rollbackErr == nil {
		return originalErr
	}
	return repository_errors.NewRollbackError(originalErr, rollbackErr)
}
