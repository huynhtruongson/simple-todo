package common

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	RootErr      error  `json:"-"`
	DebugMessage string `json:"debug_message"`
}

func (e AppError) Error() string {
	return e.RootErr.Error()
}

func NewAppError(rootErr error, code int, msg string, debugMsg string) *AppError {
	err := errors.New("")
	if rootErr != nil {
		err = rootErr
	}
	return &AppError{
		Code:         code,
		Message:      msg,
		RootErr:      err,
		DebugMessage: fmt.Sprintf(debugMsg+"->%s", err.Error()),
	}
}

func NewInternalError(err error, message string, debugMsg string) *AppError {
	appErr, ok := err.(*AppError)
	if ok {
		return NewAppError(
			err,
			http.StatusInternalServerError,
			message,
			fmt.Sprintf(debugMsg+"->%s", appErr.DebugMessage),
		)
	}
	return NewAppError(err, http.StatusInternalServerError, message, debugMsg)
}

func NewInvalidRequestError(err error, message string, debugMsg string) *AppError {
	appErr, ok := err.(*AppError)
	if ok {
		return NewAppError(
			err,
			http.StatusBadRequest,
			message,
			fmt.Sprintf(debugMsg+"->%s", appErr.DebugMessage),
		)
	}
	return NewAppError(err, http.StatusBadRequest, message, debugMsg)
}

func NewUnAuthorizedRequestError(err error, message string, debugMsg string) *AppError {
	appErr, ok := err.(*AppError)
	if ok {
		return NewAppError(
			err,
			http.StatusUnauthorized,
			message,
			fmt.Sprintf(debugMsg+"->%s", appErr.DebugMessage),
		)
	}
	return NewAppError(err, http.StatusUnauthorized, message, debugMsg)
}

func MapAppErrorToGRPCError(err error, message string) error {
	code := codes.Internal
	msg := ""
	appErr, ok := err.(*AppError)
	if ok {
		switch appErr.Code {
		// define more
		case http.StatusBadRequest:
			code = codes.InvalidArgument
		case http.StatusUnauthorized:
			code = codes.PermissionDenied
		}
		msg = appErr.Message
	}
	return status.Errorf(code, "%s: %s", message, msg)
}
