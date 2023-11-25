package auth_port

import (
	"context"

	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
)

type AuthService interface {
	Login(ctx context.Context, credential auth_entity.Credential, loginInfo auth_entity.LoginInfo) (acToken string, rfToken string, e error)
	RenewToken(ctx context.Context, rfToken string) (string, error)
}

type AuthAPI struct {
	AuthService
}

func NewAuthAPIService(authService AuthService) *AuthAPI {
	return &AuthAPI{
		AuthService: authService,
	}
}
