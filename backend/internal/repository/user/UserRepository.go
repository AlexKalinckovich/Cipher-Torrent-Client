package user

import (
	"context"
	"database/sql"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	generated "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
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
	res, err := executor.Execute(arg)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (generated.User, error) {
	res, err := r.queries.GetUserByID(ctx, id)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (generated.User, error) {
	res, err := r.queries.GetUserByEmail(ctx, email)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) GetByNickname(ctx context.Context, nickname string) (generated.User, error) {
	res, err := r.queries.GetUserByNickname(ctx, nickname)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) GetByPublicKey(ctx context.Context, publicKey []byte) (generated.User, error) {
	res, err := r.queries.GetUserByPublicKey(ctx, publicKey)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) GetStatsByUserID(ctx context.Context, userID int64) (generated.UserStat, error) {
	res, err := r.queries.GetUserStatsByUserID(ctx, userID)
	return translateResult(res, err, r.translator)
}

func (r *UserRepository) Update(ctx context.Context, arg generated.UpdateUserParams) error {
	err := r.queries.UpdateUser(ctx, arg)
	return r.translator.TranslateUserError(err)
}

func (r *UserRepository) Patch(ctx context.Context, id int64, fields user.PatchUserFields) error {
	params := r.buildPatchParams(id, fields)
	err := r.queries.PatchUser(ctx, params)
	return r.translator.TranslateUserError(err)
}

func (r *UserRepository) UpdateStats(ctx context.Context, arg generated.UpdateUserStatsParams) error {
	err := r.queries.UpdateUserStats(ctx, arg)
	return r.translator.TranslateUserError(err)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	err := r.queries.DeleteUser(ctx, id)
	return r.translator.TranslateUserError(err)
}

func (r *UserRepository) buildPatchParams(id int64, fields user.PatchUserFields) generated.PatchUserParams {
	return generated.PatchUserParams{
		ID:           id,
		Email:        toNullString(fields.Email),
		Nickname:     toNullString(fields.Nickname),
		PasswordHash: toNullString(fields.PasswordHash),
		Role:         toNullRole(fields.Role),
	}
}

func translateResult[T any](res T, err error, translator repositoryErrors.UserErrorTranslator) (T, error) {
	if translatedErr := translator.TranslateUserError(err); translatedErr != nil {
		var zero T
		return zero, translatedErr
	}
	return res, nil
}

func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func toNullRole(s *string) generated.NullUsersRole {
	if s == nil {
		return generated.NullUsersRole{}
	}
	return generated.NullUsersRole{
		UsersRole: generated.UsersRole(*s),
		Valid:     true,
	}
}
