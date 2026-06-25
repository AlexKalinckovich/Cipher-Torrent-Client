package repository_errors

import "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"

type UserEntity struct{}

func (UserEntity) EntityCode() abstract_error_code.ErrorCode {
	return "USER_NOT_FOUND"
}

func (UserEntity) EntityName() string {
	return "user"
}
