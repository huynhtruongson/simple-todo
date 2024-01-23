package auth_entity

import (
	"errors"

	"github.com/huynhtruongson/simple-todo/field"
)

type Credential struct {
	Username field.String `json:"username"`
	Password field.String `json:"password"`
}

type LoginInfo struct {
	UserAgent field.String `json:"-"`
	ClientIP  field.String `json:"-"`
}

type LoginResponse struct {
	AccessToken  field.String `json:"access_token"`
	RefreshToken field.String `json:"refresh_token"`
}

func NewLoginResponse(acToken, rfToken string) LoginResponse {
	return LoginResponse{
		AccessToken:  field.NewString(acToken),
		RefreshToken: field.NewString(rfToken),
	}
}

var (
	ErrorEmptyCredential   = errors.New("Username or Password is required")
	ErrorInvalidCredential = errors.New("Username or Password incorrect")
)
