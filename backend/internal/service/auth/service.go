package auth

import (
	"context"
	auth "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/auth/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/auth_errors"
	modelAuth "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/auth"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/service_errors"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

const RefreshTokenTTL = 7 * 24 * time.Hour

type AuthService struct {
	userDomain   auth.UserDomainPort
	hasher       auth.PasswordHasherPort
	tokenIssuer  auth.TokenIssuerPort
	refreshStore auth.RefreshTokenStorePort
}

func NewAuthService(
	userDomain auth.UserDomainPort,
	hasher auth.PasswordHasherPort,
	tokenIssuer auth.TokenIssuerPort,
	refreshStore auth.RefreshTokenStorePort,
) *AuthService {
	return &AuthService{
		userDomain:   userDomain,
		hasher:       hasher,
		tokenIssuer:  tokenIssuer,
		refreshStore: refreshStore,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (modelAuth.AuthResponse, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return modelAuth.AuthResponse{}, service_errors.NewEncryptionError(err)
	}
	return s.createUserAndIssueTokens(ctx, email, hash)
}

func (s *AuthService) createUserAndIssueTokens(ctx context.Context, email, hash string) (modelAuth.AuthResponse, error) {
	params := user.CreateUserInput{
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
	}
	userFull, err := s.userDomain.Create(ctx, params)
	if err != nil {
		return modelAuth.AuthResponse{}, err
	}
	return s.issueTokensAndBuildResult(ctx, userFull)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (modelAuth.AuthResponse, error) {
	userFull, err := s.verifyCredentials(ctx, email, password)
	if err != nil {
		return modelAuth.AuthResponse{}, err
	}
	return s.issueTokensAndBuildResult(ctx, userFull)
}

func (s *AuthService) verifyCredentials(ctx context.Context, email, password string) (userModel.UserFull, error) {
	hash, err := s.userDomain.GetPasswordHash(ctx, email)
	if err != nil {
		return userModel.UserFull{}, auth_errors.NewInvalidCredentialsError()
	}
	return s.compareAndFetchUser(hash, password, ctx, email)
}

func (s *AuthService) compareAndFetchUser(hash, password string, ctx context.Context, email string) (userModel.UserFull, error) {
	err := s.hasher.Compare(hash, password)
	if err != nil {
		return userModel.UserFull{}, auth_errors.NewInvalidCredentialsError()
	}
	return s.userDomain.GetByEmail(ctx, email)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (modelAuth.TokensResponse, error) {
	userID, err := s.refreshStore.GetUserID(ctx, refreshToken)
	if err != nil {
		return modelAuth.TokensResponse{}, auth_errors.NewInvalidRefreshTokenError()
	}
	return s.rotateTokens(ctx, userID, refreshToken)
}

func (s *AuthService) rotateTokens(ctx context.Context, userID int64, oldToken string) (modelAuth.TokensResponse, error) {
	user, err := s.userDomain.Get(ctx, userID)
	if err != nil {
		return modelAuth.TokensResponse{}, err
	}
	return s.deleteOldAndIssueNew(ctx, user, oldToken)
}

func (s *AuthService) deleteOldAndIssueNew(ctx context.Context, user userModel.UserFull, oldToken string) (modelAuth.TokensResponse, error) {
	err := s.refreshStore.Delete(ctx, oldToken)
	if err != nil {
		return modelAuth.TokensResponse{}, auth_errors.NewRefreshTokenSaveError(err)
	}
	return s.issueTokensAndBuildTokenResult(ctx, user)
}

func (s *AuthService) issueTokensAndBuildResult(ctx context.Context, user userModel.UserFull) (modelAuth.AuthResponse, error) {
	accessToken, err := s.tokenIssuer.IssueAccessToken(user.Id, user.PublicKey)
	if err != nil {
		return modelAuth.AuthResponse{}, auth_errors.NewTokenGenerationError(err)
	}
	return s.issueRefreshAndSave(ctx, user, accessToken)
}

func (s *AuthService) issueTokensAndBuildTokenResult(ctx context.Context, user userModel.UserFull) (modelAuth.TokensResponse, error) {
	accessToken, err := s.tokenIssuer.IssueAccessToken(user.Id, user.PublicKey)
	if err != nil {
		return modelAuth.TokensResponse{}, auth_errors.NewTokenGenerationError(err)
	}
	return s.issueRefreshAndSaveTokenResult(ctx, user.Id, accessToken)
}

func (s *AuthService) issueRefreshAndSave(ctx context.Context, user userModel.UserFull, accessToken string) (modelAuth.AuthResponse, error) {
	refreshToken, err := s.tokenIssuer.IssueRefreshToken()
	if err != nil {
		return modelAuth.AuthResponse{}, auth_errors.NewTokenGenerationError(err)
	}
	return s.saveRefreshAndBuildResult(ctx, user, accessToken, refreshToken)
}

func (s *AuthService) issueRefreshAndSaveTokenResult(ctx context.Context, userID int64, accessToken string) (modelAuth.TokensResponse, error) {
	refreshToken, err := s.tokenIssuer.IssueRefreshToken()
	if err != nil {
		return modelAuth.TokensResponse{}, auth_errors.NewTokenGenerationError(err)
	}
	return s.saveRefreshAndBuildTokenResult(ctx, userID, accessToken, refreshToken)
}

func (s *AuthService) saveRefreshAndBuildResult(ctx context.Context, user userModel.UserFull, accessToken, refreshToken string) (modelAuth.AuthResponse, error) {
	err := s.refreshStore.Save(ctx, refreshToken, user.Id, RefreshTokenTTL)
	if err != nil {
		return modelAuth.AuthResponse{}, auth_errors.NewRefreshTokenSaveError(err)
	}
	return s.buildAuthResponse(user, accessToken, refreshToken), nil
}

func (s *AuthService) saveRefreshAndBuildTokenResult(ctx context.Context, userID int64, accessToken, refreshToken string) (modelAuth.TokensResponse, error) {
	err := s.refreshStore.Save(ctx, refreshToken, userID, RefreshTokenTTL)
	if err != nil {
		return modelAuth.TokensResponse{}, auth_errors.NewRefreshTokenSaveError(err)
	}
	return s.buildTokenResponse(accessToken, refreshToken), nil
}

func (s *AuthService) buildAuthResponse(user userModel.UserFull, accessToken, refreshToken string) modelAuth.AuthResponse {
	return modelAuth.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.User,
	}
}

func (s *AuthService) buildTokenResponse(accessToken, refreshToken string) modelAuth.TokensResponse {
	return modelAuth.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
