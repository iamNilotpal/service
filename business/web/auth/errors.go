package auth

import (
	"errors"
	"fmt"
)

// AuthError is used to pass an error during the request through the
// application with auth specific context.
type AuthError struct {
	Code int
	Msg  string
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(code int, format string, args ...any) error {
	return &AuthError{
		Code: code,
		Msg:  fmt.Sprintf(format, args...),
	}
}

// Error implements the error interface. It uses the default message of the
// wrapped error. This is what will be shown in the services' logs.
func (ae *AuthError) Error() string {
	return ae.Msg
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetAuthError(err error) *AuthError {
	var re *AuthError

	if !errors.As(err, &re) {
		return nil
	}

	return re
}
