package common

import (
	"fmt"
	"net/http"
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
	errMsg := ""
	if rootErr != nil {
		errMsg = rootErr.Error()
	}
	return &AppError{
		Code:         code,
		Message:      msg,
		RootErr:      rootErr,
		DebugMessage: fmt.Sprintf(debugMsg+"->%s", errMsg),
	}
}

func NewCustomError(err error, code int, message string, debugMsg string) *AppError {
	appErr, ok := err.(*AppError)
	if ok {
		return NewAppError(
			err,
			http.StatusInternalServerError,
			message,
			fmt.Sprintf(debugMsg+"->%s", appErr.DebugMessage),
		)
	}
	return NewAppError(err, code, message, debugMsg)
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
