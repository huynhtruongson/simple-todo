package auth_biz

import (
	"context"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	"github.com/huynhtruongson/simple-todo/middleware"
	"github.com/huynhtruongson/simple-todo/token"
)

type RenewTokenBiz struct {
	DB         lib.DB
	TokenMaker token.TokenMaker
	SessionRepo
}

func NewRenewTokenBiz(db lib.DB, tokenMaker token.TokenMaker, sessionRepo SessionRepo) *RenewTokenBiz {
	return &RenewTokenBiz{
		DB:          db,
		TokenMaker:  tokenMaker,
		SessionRepo: sessionRepo,
	}
}

func (b RenewTokenBiz) RenewToken(ctx context.Context, rfToken string) (string, error) {
	payload, err := b.TokenMaker.VerifyToken(rfToken)
	if err != nil {
		return "", common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, "")
	}
	sessions, err := b.SessionRepo.GetSessionByIds(ctx, b.DB, uuid.UUIDs{payload.ID})
	if err != nil {
		return "", common.NewInternalError(err, common.InternalErrorMessage, "RenewToken.GetSessionByIds")
	}
	switch {
	case len(sessions) != 1:
		return "", common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, "")
	case sessions[0].RefreshToken != rfToken:
		return "", common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, "")
	case sessions[0].IsBlocked:
		return "", common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, "")
	}
	token, _, err := b.TokenMaker.CreateToken(sessions[0].UserID, token.AccessTokenDuration, token.AccessToken)
	if err != nil {
		return "", common.NewInternalError(err, common.InternalErrorMessage, "RenewToken.CreateToken")
	}
	return token, nil
}
