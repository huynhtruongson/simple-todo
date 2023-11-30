package auth_entity

import "errors"

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginInfo struct {
	UserAgent string `json:"-"`
	ClientIP  string `json:"-"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewLoginResponse(acToken, rfToken string) LoginResponse {
	return LoginResponse{
		AccessToken:  acToken,
		RefreshToken: rfToken,
	}
}

var (
	ErrorEmptyCredential   = errors.New("Username or Password is required")
	ErrorInvalidCredential = errors.New("Username or Password incorrect")
)
