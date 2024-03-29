package auth_service

import (
	"context"
	"errors"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Login(ctx context.Context, cred auth_entity.Credential, info auth_entity.LoginInfo) (acToken string, rfToken string, e error) {
	if cred.Username.String() == "" || cred.Password.String() == "" {
		e = common.NewInvalidRequestError(auth_entity.ErrorEmptyCredential, auth_entity.ErrorEmptyCredential.Error(), "Login.UserRepo.GetUsersByUsername")
		return
	}
	users, err := s.UserRepo.GetUsersByUsername(ctx, s.DB, cred.Username.String())
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.UserRepo.GetUsersByUsername")
		return
	}
	if len(users) != 1 {
		e = common.NewInvalidRequestError(auth_entity.ErrorInvalidCredential, auth_entity.ErrorInvalidCredential.Error(), "")
		return
	}
	if err := utils.CheckPassword(cred.Password.String(), users[0].Password.String()); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			e = common.NewInvalidRequestError(auth_entity.ErrorInvalidCredential, auth_entity.ErrorInvalidCredential.Error(), "")
			return
		}
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CheckPassword")
		return
	}
	acToken, _, err = s.TokenMaker.CreateToken(users[0].UserID.Int(), token.AccessTokenDuration, token.AccessToken)
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CreateToken")
		return
	}

	rfToken, payload, err := s.TokenMaker.CreateToken(users[0].UserID.Int(), token.RefreshTokenDuration, token.RefreshToken)
	if err != nil {
		e = common.NewInternalError(err, common.InternalErrorMessage, "Login.CreateRefreshToken")
		return
	}
	session := auth_entity.NewSession(payload.ID, payload.UserID, rfToken, payload.ExpiresAt, info.UserAgent.String(), info.ClientIP.String())

	if err := lib.ExecTX(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		return s.SessionRepo.CreateSession(ctx, tx, session)
	}); err != nil {
		e = err
		return
	}

	return
}
