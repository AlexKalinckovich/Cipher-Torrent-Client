package user

import (
	"context"
	"database/sql"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
)

type UserRepository struct {
	database   *sql.DB
	queries    *generated.Queries
	translator repositoryErrors.UserErrorTranslator
}

func NewUserRepository(database *sql.DB, queries *generated.Queries) *UserRepository {
	return &UserRepository{
		database:   database,
		queries:    queries,
		translator: repositoryErrors.Translator,
	}
}

func (r *UserRepository) Create(ctx context.Context, arg generated.CreateUserParams) (generated.User, error) {
	executor := NewUserCreateExecutor(ctx, r.database, r.queries)
	return executor.Execute(arg)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (generated.User, error) {
	res, err := r.queries.GetUserByID(ctx, id)
	return r.handleUserResult(res, err)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (generated.User, error) {
	res, err := r.queries.GetUserByEmail(ctx, email)
	return r.handleUserResult(res, err)
}

func (r *UserRepository) GetByPublicKey(ctx context.Context, publicKey string) (generated.User, error) {
	res, err := r.queries.GetUserByPublicKey(ctx, publicKey)
	return r.handleUserResult(res, err)
}

func (r *UserRepository) handleUserResult(res generated.User, err error) (generated.User, error) {
	translated := r.translator.TranslateUserError(err)
	if translated != nil {
		return generated.User{}, translated
	}
	return res, nil
}

func (r *UserRepository) GetStatsByUserID(ctx context.Context, userID int64) (generated.UserStat, error) {
	res, err := r.queries.GetUserStatsByUserID(ctx, userID)
	return r.handleStatsResult(res, err)
}

func (r *UserRepository) handleStatsResult(res generated.UserStat, err error) (generated.UserStat, error) {
	translated := r.translator.TranslateUserError(err)
	if translated != nil {
		return generated.UserStat{}, translated
	}
	return res, nil
}

func (r *UserRepository) Update(ctx context.Context, arg generated.UpdateUserParams) error {
	return r.queries.UpdateUser(ctx, arg)
}

func (r *UserRepository) UpdateStats(ctx context.Context, arg generated.UpdateUserStatsParams) error {
	return r.queries.UpdateUserStats(ctx, arg)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}
