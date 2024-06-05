package internal

import (
	"errors"
)

const (
	ENOTFOUND ErrorKind = "notfound"
	EINTERNAL ErrorKind = "internal"
	EBADINPUT ErrorKind = "badinput"
)

const (
	DefaultErrorKind    = EINTERNAL
	DefaultErrorMessage = "Internal error"
	DefaultErrorCode    = "internal_error"
)

type ErrorKind string

// appError represents application level error, it should not expose error
// verbosely. Any unknown error should be served as server error.
type appError interface {
	error

	// Kind returns ErrorKind string. For example this can be useful to
	// translate error to suitable status code in http layer.
	Kind() ErrorKind

	// Code is error code that can be used as identifier of the error.
	Code() string
}

// Error implements Error
type Error struct {
	kind    ErrorKind
	code    string
	message string
}

func NewError(kind ErrorKind, message, code string) *Error {
	return &Error{
		kind:    kind,
		message: message,
		code:    code,
	}
}

func (c *Error) Error() string {
	return c.message
}

func (c *Error) Kind() ErrorKind {
	return c.kind
}

func (c *Error) Code() string {
	return c.code
}

// ParseError parses err and hide internal level or implementation error.
// Any Error non-compatible will be ignored and returns EINTERNAL kind of error.
func ParseError(err error) (kind ErrorKind, message, code string) {
	var E appError
	if errors.As(err, &E) {
		kind = E.Kind()
		message = E.Error()
		code = E.Code()
		return
	}

	kind = DefaultErrorKind
	message = DefaultErrorMessage
	code = DefaultErrorCode
	return
}
