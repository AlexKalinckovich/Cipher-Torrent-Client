package service_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DatabaseUnavailableCode abstract_error_code.ErrorCode = "DATABASE_UNAVAILABLE"
)

type DatabaseUnavailableError struct{}

func (e *DatabaseUnavailableError) Code() abstract_error_code.ErrorCode {
	return DatabaseUnavailableCode
}
func (e *DatabaseUnavailableError) Error() string { return "database unavailable" }
func (e *DatabaseUnavailableError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateServiceUnavailableResponse(DatabaseUnavailableCode, e.Error())
}
