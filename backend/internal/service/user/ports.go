package user

import (
	"context"
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

type UserRepositoryPort interface {
	Create(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int64) (db.User, error)
	GetStatsByUserID(ctx context.Context, userID int64) (db.UserStat, error)
	Update(ctx context.Context, arg db.UpdateUserParams) error
	Delete(ctx context.Context, id int64) error
}

type UserServicePort interface {
	Create(ctx context.Context, params CreateUserInput) (userModel.UserFull, error)
	Get(ctx context.Context, id int64) (userModel.UserFull, error)
	Update(ctx context.Context, id int64, params UpdateUserInput) (userModel.UserFull, error)
	Patch(ctx context.Context, id int64, fields PatchUserFields) (userModel.UserFull, error)
	Delete(ctx context.Context, id int64) error
}

type UserMapperPort interface {
	ToUserFullDTO(u db.User, s db.UserStat) userModel.UserFull
}

type UserValidatorPort interface {
	ValidateCreate(params CreateUserInput) error
	ValidateUpdate(params UpdateUserInput) error
	ValidatePatch(fields PatchUserFields) error
}
