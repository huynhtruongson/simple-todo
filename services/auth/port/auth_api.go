package auth_port

import "context"

type AuthService interface {
	Login(ctx context.Context, username, password string) (acToken string, rfToken string, e error)
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
