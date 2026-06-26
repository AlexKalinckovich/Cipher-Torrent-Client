package not_found

import (
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

type Entity interface {
	EntityCode() abstract_error_code.ErrorCode
	EntityName() string
}

type EntityNotFoundBase interface {
	isNotFound()
}

type NotFoundError[T Entity] struct {
	code    abstract_error_code.ErrorCode
	message string
}

func NewNotFoundError[T Entity]() *NotFoundError[T] {
	var zero T
	return &NotFoundError[T]{
		code:    zero.EntityCode(),
		message: fmt.Sprintf("%s not found", zero.EntityName()),
	}
}

func CodeFor[T Entity]() abstract_error_code.ErrorCode {
	var zero T
	return zero.EntityCode()
}

func (e *NotFoundError[T]) Code() abstract_error_code.ErrorCode {
	return e.code
}

func (e *NotFoundError[T]) Error() string {
	return e.message
}

func (e *NotFoundError[T]) isNotFound() {}
