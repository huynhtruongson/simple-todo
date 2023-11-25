package auth_entity

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

const (
	ErrorEmptyCredential   = "Username or Password is required"
	ErrorInvalidCredential = "Username or Password incorrect"
)
