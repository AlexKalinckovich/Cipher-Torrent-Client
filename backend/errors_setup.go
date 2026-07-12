package main

import (
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
	"net/http"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
)

func buildErrorRegistry() *transport.ErrorRegistry {
	registry := transport.NewErrorRegistry()
	registerUserNotFoundHandler(registry)
	return registry
}

func registerUserNotFoundHandler(registry *transport.ErrorRegistry) {
	code := not_found.CodeFor[repositoryErrors.UserEntity]()
	registry.Register(code, createNotFoundHandler(code))
}

func createNotFoundHandler(code abstract_error_code.ErrorCode) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusNotFound, common.ApiError{
			Status:    http.StatusNotFound,
			ErrorCode: string(code),
			Message:   "the requested user does not exist",
		})
	}
}
