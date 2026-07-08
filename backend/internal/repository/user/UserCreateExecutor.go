package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/tx"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
)

type UserCreateExecutor struct {
	ctx      context.Context
	database *sql.DB
	queries  *generated.Queries
}

func NewUserCreateExecutor(
	ctx context.Context,
	database *sql.DB,
	queries *generated.Queries,
) *UserCreateExecutor {
	return &UserCreateExecutor{
		ctx:      ctx,
		database: database,
		queries:  queries,
	}
}

func (e *UserCreateExecutor) Execute(arg generated.CreateUserParams) (generated.User, error) {
	beginTx, err := e.database.BeginTx(e.ctx, nil)
	return e.tryInsertUser(beginTx, arg, err)
}

func (e *UserCreateExecutor) tryInsertUser(
	sqlTx *sql.Tx,
	arg generated.CreateUserParams,
	err error,
) (generated.User, error) {
	if err != nil {
		return generated.User{}, err
	}
	qtx := e.queries.WithTx(sqlTx)
	roller := tx.NewRollbacker(sqlTx)
	res, execErr := qtx.CreateUser(e.ctx, arg)
	return e.tryResolveInsertID(roller, qtx, arg, res, execErr)
}

func (e *UserCreateExecutor) tryResolveInsertID(
	roller *tx.Rollbacker,
	qtx *generated.Queries,
	arg generated.CreateUserParams,
	res sql.Result,
	err error,
) (generated.User, error) {
	if err != nil {
		return generated.User{}, roller.Rollback(err)
	}
	id, idErr := res.LastInsertId()
	return e.tryInsertStats(roller, qtx, arg, id, idErr)
}

func (e *UserCreateExecutor) tryInsertStats(
	roller *tx.Rollbacker,
	qtx *generated.Queries,
	arg generated.CreateUserParams,
	id int64,
	err error,
) (generated.User, error) {
	if err != nil {
		return generated.User{}, roller.Rollback(err)
	}
	statsErr := qtx.CreateUserStats(e.ctx, generated.CreateUserStatsParams{UserID: id})
	return e.tryCommit(roller, arg, id, statsErr)
}

func (e *UserCreateExecutor) tryCommit(
	roller *tx.Rollbacker,
	arg generated.CreateUserParams,
	id int64,
	err error,
) (generated.User, error) {
	if err != nil {
		return generated.User{}, roller.Rollback(err)
	}
	return e.commit(roller.Tx, arg, id)
}

func (e *UserCreateExecutor) commit(
	sqlTx *sql.Tx,
	arg generated.CreateUserParams,
	id int64,
) (generated.User, error) {
	if err := sqlTx.Commit(); err != nil {
		return generated.User{}, err
	}
	return buildUser(id, arg), nil
}

func buildUser(id int64, arg generated.CreateUserParams) generated.User {
	return generated.User{
		ID:            id,
		Email:         arg.Email,
		PublicKey:     arg.PublicKey,
		PrivateKeyEnc: arg.PrivateKeyEnc,
		Nickname:      arg.Nickname,
		Role:          arg.Role,
		CreatedAt:     time.Now(),
	}
}
