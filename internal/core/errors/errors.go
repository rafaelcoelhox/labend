package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error types
var (
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInternal      = errors.New("internal error")
)

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e AppError) Unwrap() error {
	return e.Err
}

func NotFound(entity string, id interface{}) error {
	return AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s with id %v not found", entity, id),
		Err:     ErrNotFound,
	}
}

func AlreadyExists(entity string, field string, value interface{}) error {
	return AppError{
		Code:    "ALREADY_EXISTS",
		Message: fmt.Sprintf("%s with %s %v already exists", entity, field, value),
		Err:     ErrAlreadyExists,
	}
}

func InvalidInput(msg string) error {
	return AppError{
		Code:    "INVALID_INPUT",
		Message: msg,
		Err:     ErrInvalidInput,
	}
}

func Internal(err error) error {
	return AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Err:     errors.Wrap(err, "internal error"),
	}
}

// Wrap - para adicionar contexto
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Is - helper para checagem de tipo
func Is(err, target error) bool {
	return errors.Is(err, target)
}
