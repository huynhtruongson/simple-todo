package common

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/stretchr/testify/assert"
)

var (
	ResponseKey       = "ResponseKey"
	RequestPayloadKey = "RequestPayload"
	TokenKey          = "TokenKey"
	TokenPayloadKey   = "TokenPayloadKey"
)

func (s *Suite) AssertErrorResponse(ctx context.Context, msg string) (context.Context, error) {
	stepState := StepStateFromContext(ctx)
	respErr, ok := stepState.Get(ResponseKey).(common.AppError)
	if !ok {
		return ctx, fmt.Errorf("retrieve response error failed")
	}
	err := AssertExpectedAndActual(assert.Equal, http.StatusBadRequest, respErr.Code, "response status code mismatch")
	if err != nil {
		return ctx, err
	}
	err = AssertExpectedAndActual(assert.Equal, msg, respErr.Message, "response error message mismatch")
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (s *Suite) UserSignedIn(ctx context.Context) (context.Context, error) {
	stepState := StepStateFromContext(ctx)
	tokenMaker, err := token.NewPasetoMaker(s.Config.TokenKey)
	if err != nil {
		return ctx, fmt.Errorf("init token maker failed,%w\n", err)
	}
	acToken, payload, err := tokenMaker.CreateToken(1, time.Minute*15, token.AccessToken)
	if err != nil {
		return ctx, fmt.Errorf("create token failed,%w\n", err)
	}
	stepState.Set(TokenKey, acToken)
	stepState.Set(TokenPayloadKey, payload)

	return StepStateToContext(ctx, stepState), nil
}
