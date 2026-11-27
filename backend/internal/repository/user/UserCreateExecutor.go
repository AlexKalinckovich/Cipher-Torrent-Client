package user

import (
	"context"
	"database/sql"
	"time"

	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
)

type UserCreateExecutor struct {
	database *sql.DB
	queries  *generated.Queries
}

func NewUserCreateExecutor(
	database *sql.DB,
	queries *generated.Queries,
) *UserCreateExecutor {
	return &UserCreateExecutor{
		database: database,
		queries:  queries,
	}
}

func (e *UserCreateExecutor) Execute(ctx context.Context, arg generated.CreateUserParams) (generated.User, error) {
	sqlTx, err := e.database.BeginTx(ctx, nil)
	if err != nil {
		return generated.User{}, err
	}
	defer func() { _ = sqlTx.Rollback() }()

	qtx := e.queries.WithTx(sqlTx)
	res, err := qtx.CreateUser(ctx, arg)
	if err != nil {
		return generated.User{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return generated.User{}, err
	}
	if err := qtx.CreateUserStats(ctx, generated.CreateUserStatsParams{UserID: id}); err != nil {
		return generated.User{}, err
	}
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
