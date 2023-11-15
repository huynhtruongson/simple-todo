package common

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
)

const (
	InternalErrorMessage       = "Something went wrong"
	InvalidRequestErrorMessage = "Invalid request params"
)

var (
	// token
	ExpiredTokenErrorMessage   = errors.New("Token is expired")
	InvalidTokenErrorMessage   = errors.New("Invalid token")
	InvalidKeySizeErrorMessage = fmt.Errorf("invalid key size, must be exactly %d characters", chacha20poly1305.KeySize)
)
