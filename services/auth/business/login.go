package auth_biz

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error)
}
type SessionRepo interface {
	CreateSession(ctx context.Context, db lib.QueryExecer, session auth_entity.Session) error
	GetSessionByIds(ctx context.Context, db lib.QueryExecer, ids uuid.UUIDs) ([]auth_entity.Session, error)
}
type LoginBiz struct {
	DB         lib.DB
	TokenMaker token.TokenMaker
	UserRepo
	SessionRepo
}

func NewLoginBiz(db lib.DB, tokenMaker token.TokenMaker, userRepo UserRepo, sessionRepo SessionRepo) *LoginBiz {
	return &LoginBiz{
		DB:          db,
		TokenMaker:  tokenMaker,
		UserRepo:    userRepo,
		SessionRepo: sessionRepo,
	}
}

func (b LoginBiz) Login(ctx context.Context, username, password string) (acToken string, rfToken string, e error) {
	if username == "" || password == "" {
		e = common.NewInvalidRequestError(nil, auth_entity.ErrorEmptyCredential, "Login.UserRepo.GetUsersByUsername")
		return
	}
	users, err := b.UserRepo.GetUsersByUsername(ctx, b.DB, username)
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.UserRepo.GetUsersByUsername")
		return
	}
	if len(users) != 1 {
		e = common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, "")
		return
	}
	if err := utils.CheckPassword(password, users[0].Password); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			e = common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, "")
			return
		}
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CheckPassword")
		return
	}
	acToken, _, err = b.TokenMaker.CreateToken(users[0].UserID, token.AccessTokenDuration, token.AccessToken)
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CreateToken")
		return
	}

	rfToken, payload, err := b.TokenMaker.CreateToken(users[0].UserID, token.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CreateRefreshToken")
		return
	}
	session := auth_entity.NewSession(payload.ID, payload.UserID, rfToken, payload.ExpiresAt)

	if err := lib.ExecTX(ctx, b.DB, func(ctx context.Context, tx pgx.Tx) error {
		return b.SessionRepo.CreateSession(ctx, tx, session)
	}); err != nil {
		e = err
		return
	}

	return
}
