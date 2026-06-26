package user

import (
	"context"
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
	"time"
)

type UserService struct {
	repository UserRepositoryPort
	mapper     UserMapperPort
	validator  UserValidatorPort
}

func NewUserService(
	repository UserRepositoryPort,
	mapper UserMapperPort,
	validator UserValidatorPort,
) *UserService {
	return &UserService{
		repository: repository,
		mapper:     mapper,
		validator:  validator,
	}
}

func (s *UserService) Create(ctx context.Context, params CreateUserInput) (userModel.UserFull, error) {
	err := s.validator.ValidateCreate(params)
	return s.handleCreateValidation(ctx, params, err)
}

func (s *UserService) handleCreateValidation(
	ctx context.Context,
	params CreateUserInput,
	err error,
) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	userRow, repoErr := s.repository.Create(ctx, db.CreateUserParams{
		Email:     params.Email,
		PublicKey: params.PublicKey,
		Nickname:  params.Nickname,
		Role:      db.UsersRole(params.Role),
		CreatedAt: time.Now(),
	})
	return s.handleCreateResult(ctx, userRow, repoErr)
}

func (s *UserService) handleCreateResult(ctx context.Context, userRow db.User, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) Get(ctx context.Context, id int64) (userModel.UserFull, error) {
	userRow, err := s.repository.GetByID(ctx, id)
	return s.handleGetResult(ctx, userRow, err)
}

func (s *UserService) handleGetResult(ctx context.Context, userRow db.User, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) Update(ctx context.Context, id int64, params UpdateUserInput) (userModel.UserFull, error) {
	err := s.validator.ValidateUpdate(params)
	return s.handleUpdateValidation(ctx, id, params, err)
}

func (s *UserService) handleUpdateValidation(ctx context.Context, id int64, params UpdateUserInput, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	updateErr := s.repository.Update(ctx, db.UpdateUserParams{
		ID:        id,
		Email:     params.Email,
		PublicKey: params.PublicKey,
		Nickname:  params.Nickname,
		Role:      db.UsersRole(params.Role),
	})
	return s.handleUpdateResult(ctx, id, updateErr)
}

func (s *UserService) handleUpdateResult(ctx context.Context, id int64, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.Get(ctx, id)
}

func (s *UserService) Patch(ctx context.Context, id int64, fields PatchUserFields) (userModel.UserFull, error) {
	err := s.validator.ValidatePatch(fields)
	return s.handlePatchValidation(ctx, id, fields, err)
}

func (s *UserService) handlePatchValidation(ctx context.Context, id int64, fields PatchUserFields, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	userRow, fetchErr := s.repository.GetByID(ctx, id)
	return s.handlePatchFetch(ctx, userRow, fields, fetchErr)
}

func (s *UserService) handlePatchFetch(
	ctx context.Context,
	userRow db.User,
	fields PatchUserFields,
	err error,
) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	params := NewPatchParamsBuilder(userRow, fields).Build()
	updateErr := s.repository.Update(ctx, params)
	return s.handlePatchResult(ctx, userRow.ID, updateErr)
}

func (s *UserService) handlePatchResult(ctx context.Context, id int64, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.Get(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *UserService) fetchUserFull(ctx context.Context, userRow db.User) (userModel.UserFull, error) {
	statsRow, err := s.repository.GetStatsByUserID(ctx, userRow.ID)
	return s.handleFetchStats(userRow, statsRow, err)
}

func (s *UserService) handleFetchStats(userRow db.User, statsRow db.UserStat, err error) (userModel.UserFull, error) {
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.mapper.ToUserFullDTO(userRow, statsRow), nil
}
