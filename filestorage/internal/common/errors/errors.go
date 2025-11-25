package errors

import (
	stdErrors "errors"
)

type Code string

const (
	CodeValidation Code = "validation_error"
	CodeNotFound   Code = "not_found"
	CodeDatabase   Code = "database_error"
	CodeStorage    Code = "storage_error"
	CodeInternal   Code = "internal_error"
)

type Error struct {
	code    Code
	message string
	err     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	if e.err != nil {
		return e.message + ": " + e.err.Error()
	}
	return e.message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) Code() Code {
	if e == nil {
		return CodeInternal
	}
	return e.code
}

func (e *Error) Message() string {
	if e == nil {
		return ""
	}
	return e.message
}

func New(code Code, message string) error {
	return &Error{
		code:    code,
		message: message,
	}
}

func Wrap(err error, code Code, message string) error {
	if err == nil {
		return New(code, message)
	}
	return &Error{
		code:    code,
		message: message,
		err:     err,
	}
}

func CodeOf(err error) Code {
	var appErr *Error
	if stdErrors.As(err, &appErr) {
		return appErr.code
	}
	return CodeInternal
}

func IsCode(err error, code Code) bool {
	return CodeOf(err) == code
}

func Message(err error) string {
	var appErr *Error
	if stdErrors.As(err, &appErr) {
		return appErr.message
	}
	return ""
}
