package auth_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/interceptor"
	"github.com/huynhtruongson/simple-todo/token"
)

func (s *AuthService) RenewToken(ctx context.Context, rfToken string) (string, error) {
	payload, err := s.TokenMaker.VerifyToken(rfToken)
	if err != nil {
		return "", common.NewUnAuthorizedRequestError(err, interceptor.UnAuthorizedMessage, "")
	}
	sessions, err := s.SessionRepo.GetSessionByIds(ctx, s.DB, uuid.UUIDs{payload.ID})
	if err != nil {
		return "", common.NewInternalError(err, common.InternalErrorMessage, "RenewToken.GetSessionByIds")
	}
	switch {
	case len(sessions) != 1:
		return "", common.NewUnAuthorizedRequestError(err, interceptor.UnAuthorizedMessage, "")
	case sessions[0].RefreshToken != rfToken:
		return "", common.NewUnAuthorizedRequestError(err, interceptor.UnAuthorizedMessage, "")
	case sessions[0].IsBlocked:
		return "", common.NewUnAuthorizedRequestError(err, interceptor.UnAuthorizedMessage, "")
	}
	token, _, err := s.TokenMaker.CreateToken(sessions[0].UserID, token.AccessTokenDuration, token.AccessToken)
	if err != nil {
		return "", common.NewInternalError(err, common.InternalErrorMessage, "RenewToken.CreateToken")
	}
	return token, nil
}
