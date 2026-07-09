package user

import (
	"context"
	"time"

	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	serviceErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/service_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/dpki"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

type UserService struct {
	repository UserRepositoryPort
	mapper     UserMapperPort
	validator  UserValidatorPort
	crypto     security.CryptoServicePort
}

func NewUserService(
	repository UserRepositoryPort,
	mapper UserMapperPort,
	validator UserValidatorPort,
	crypto security.CryptoServicePort,
) *UserService {
	return &UserService{
		repository: repository,
		mapper:     mapper,
		validator:  validator,
		crypto:     crypto,
	}
}

func (s *UserService) Create(ctx context.Context, params CreateUserInput) (userModel.UserFull, error) {
	if err := s.validator.ValidateCreate(params); err != nil {
		return userModel.UserFull{}, err
	}
	keyPair, err := dpki.GenerateIdentity()
	if err != nil {
		return userModel.UserFull{}, serviceErrors.NewDpkiError(err)
	}
	encryptedKey, err := s.crypto.EncryptPrivateKey(keyPair.PrivateKey)
	if err != nil {
		return userModel.UserFull{}, serviceErrors.NewEncryptionError(err)
	}

	userRow, err := s.repository.Create(ctx, s.buildCreateParams(params, keyPair, encryptedKey))
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) buildCreateParams(params CreateUserInput, keyPair *dpki.KeyPair, encryptedKey []byte) db.CreateUserParams {
	return db.CreateUserParams{
		Email:         params.Email,
		PublicKey:     keyPair.PublicKey,
		PrivateKeyEnc: encryptedKey,
		Nickname:      params.Nickname,
		PasswordHash:  params.PasswordHash,
		Role:          db.UsersRole(params.Role),
		CreatedAt:     time.Now(),
	}
}

func (s *UserService) Get(ctx context.Context, id int64) (userModel.UserFull, error) {
	userRow, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (userModel.UserFull, error) {
	userRow, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) GetByNickname(ctx context.Context, nickname string) (userModel.UserFull, error) {
	userRow, err := s.repository.GetByNickname(ctx, nickname)
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) GetByPublicKey(ctx context.Context, publicKey []byte) (userModel.UserFull, error) {
	userRow, err := s.repository.GetByPublicKey(ctx, publicKey)
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.fetchUserFull(ctx, userRow)
}

func (s *UserService) Update(ctx context.Context, id int64, params UpdateUserInput) (userModel.UserFull, error) {
	if err := s.validator.ValidateUpdate(params); err != nil {
		return userModel.UserFull{}, err
	}
	updateErr := s.repository.Update(ctx, db.UpdateUserParams{
		ID:           id,
		Email:        params.Email,
		Nickname:     params.Nickname,
		PasswordHash: params.PasswordHash,
		Role:         db.UsersRole(params.Role),
	})
	if updateErr != nil {
		return userModel.UserFull{}, updateErr
	}
	return s.Get(ctx, id)
}

func (s *UserService) Patch(ctx context.Context, id int64, fields PatchUserFields) (userModel.UserFull, error) {
	if err := s.validator.ValidatePatch(fields); err != nil {
		return userModel.UserFull{}, err
	}
	if err := s.repository.Patch(ctx, id, fields); err != nil {
		return userModel.UserFull{}, err
	}
	return s.Get(ctx, id)
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *UserService) fetchUserFull(ctx context.Context, userRow db.User) (userModel.UserFull, error) {
	statsRow, err := s.repository.GetStatsByUserID(ctx, userRow.ID)
	if err != nil {
		return userModel.UserFull{}, err
	}
	return s.mapper.ToUserFullDTO(userRow, statsRow), nil
}

func (s *UserService) GetPasswordHash(ctx context.Context, email string) (string, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	return user.PasswordHash, nil
}
