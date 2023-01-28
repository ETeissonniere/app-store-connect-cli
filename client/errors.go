package client

import (
	"errors"
	"fmt"
	"net/http"
)

// Shared error between methods go here
var (
	ErrUnauthorized   = errors.New("unauthorized, check your credentials")
	ErrInternalServer = errors.New("internal server error")
	ErrUnkown         = errors.New("unknown error")
)

func handleGeneralError(statusCode int, action string) error {
	var err error
	switch statusCode {
	case http.StatusUnauthorized:
		err = ErrUnauthorized
	case http.StatusInternalServerError:
		err = ErrInternalServer
	default:
		err = ErrUnkown
	}

	return fmt.Errorf("%s: %w", action, err)
}
