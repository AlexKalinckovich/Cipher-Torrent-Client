package transport

import (
	"errors"
	customErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/transport"
	"net/http"
	"sync"
)

type ErrorHandler func(err error) HTTPResponse

type ErrorCodeProvider interface {
	Code() customErrors.ErrorCode
}

type ErrorTranslator interface {
	Translate(err error) HTTPResponse
}

type ErrorRegistry struct {
	mu       sync.RWMutex
	handlers map[customErrors.ErrorCode]ErrorHandler
	fallback ErrorHandler
}

func NewErrorRegistry() *ErrorRegistry {
	return &ErrorRegistry{
		handlers: make(map[customErrors.ErrorCode]ErrorHandler),
		fallback: DefaultFallbackHandler,
	}
}

func (r *ErrorRegistry) Register(code customErrors.ErrorCode, handler ErrorHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[code] = handler
}

func DefaultFallbackHandler(_ error) HTTPResponse {
	return NewHTTPResponse(http.StatusInternalServerError, map[string]any{
		"code":    transport.InternalErrorCode,
		"message": "unexpected internal server error",
	})
}

func (r *ErrorRegistry) Translate(err error) HTTPResponse {
	code := ExtractErrorCode(err)
	handler := r.ResolveHandler(code)
	return handler(err)
}

func ExtractErrorCode(err error) customErrors.ErrorCode {
	var provider ErrorCodeProvider

	if errors.As(err, &provider) {
		return provider.Code()
	}

	return transport.InternalErrorCode
}

func (r *ErrorRegistry) ResolveHandler(code customErrors.ErrorCode) ErrorHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[code]
	if exists {
		return handler
	}
	return r.fallback
}
