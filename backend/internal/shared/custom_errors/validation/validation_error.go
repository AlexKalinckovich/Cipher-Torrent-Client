package validation

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"strings"
)

const (
	ErrorCode          abstract_error_code.ErrorCode = "VALIDATION_ERROR"
	AggregateErrorCode abstract_error_code.ErrorCode = "VALIDATION_AGGREGATE_ERROR"
)

type FieldError struct {
	Field   string `json:"field"`
	Value   any    `json:"value,omitempty"`
	Message string `json:"message"`
}

type Error struct {
	ValidationError FieldError
}

type AggregateError struct {
	Errors []FieldError
}

func NewError(field string, value any, message string) *Error {
	return &Error{ValidationError: createFieldError(field, value, message)}
}

func NewAggregateError() *AggregateError {
	return &AggregateError{}
}

func (e *AggregateError) Add(field string, value any, message string) {
	e.Errors = append(e.Errors, createFieldError(field, value, message))
}

func createFieldError(field string, value any, message string) FieldError {
	return FieldError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

func (e *AggregateError) HasErrors() bool {
	return len(e.Errors) > 0
}

func (e *AggregateError) Error() string {
	if len(e.Errors) == 0 {
		return "validation passed"
	}
	msgs := make([]string, len(e.Errors))
	for i, fe := range e.Errors {
		msgs[i] = fe.Field + ": " + fe.Message
	}
	return strings.Join(msgs, "; ")
}

func (e *AggregateError) Code() abstract_error_code.ErrorCode {
	return AggregateErrorCode
}
