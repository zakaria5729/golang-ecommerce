package apperror

import (
	"errors"
	"fmt"
)

type ServerError struct {
	msg string
	err error
}

func (e *ServerError) Error() string {
	return e.msg
}

func (e *ServerError) Unwrap() error {
	return e.err
}

func NewServerError(msg string) error {
	return &ServerError{msg: msg}
}

func WrapServerError(msg string, err error) error {
	return &ServerError{msg: msg, err: err}
}

func ServerErrorf(format string, args ...any) error {
	baseErr := fmt.Errorf(format, args...)
	wrappedErr := errors.Unwrap(baseErr)

	return &ServerError{
		msg: baseErr.Error(),
		err: wrappedErr,
	}
}
