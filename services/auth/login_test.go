package auth_service

import (
	"context"
	"testing"

	"github.com/huynhtruongson/simple-todo/common"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/auth"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_token "github.com/huynhtruongson/simple-todo/mocks/token"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceProp struct {
	DB          *mock_db.DB
	TX          *mock_db.Tx
	UserRepo    *mock_repo.UserRepo
	SessionRepo *mock_repo.SessionRepo
	TokenMaker  *mock_token.TokenMaker
}

func NewMockAuthService(t *testing.T) (*AuthService, *MockServiceProp) {
	userRepo := mock_repo.NewUserRepo(t)
	sessionRepo := mock_repo.NewSessionRepo(t)
	tokenMaker := mock_token.NewTokenMaker(t)
	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
	return &AuthService{
			DB:          db,
			UserRepo:    userRepo,
			SessionRepo: sessionRepo,
			TokenMaker:  tokenMaker,
		}, &MockServiceProp{
			DB:          db,
			TX:          tx,
			UserRepo:    userRepo,
			SessionRepo: sessionRepo,
			TokenMaker:  tokenMaker,
		}
}
func TestLoginBiz_Login(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	hashedPwd, err := utils.HashPassword("password")
	assert.NoError(t, err)
	mockACTokenPayload, err := token.NewTokenPayload(1, token.AccessTokenDuration, token.AccessToken)
	assert.NoError(t, err)
	mockRFTokenPayload, err := token.NewTokenPayload(1, token.RefreshTokenDuration, token.RefreshToken)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		cred      auth_entity.Credential
		mock      func(prop *MockServiceProp)
		expectErr *common.AppError
	}{
		{
			name: "should validate empty username correctly",
			cred: auth_entity.Credential{
				Username: "",
				Password: "password",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorEmptyCredential, "Login.UserRepo.GetUsersByUsername"),
		},
		{
			name: "should validate empty password correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "",
			},
			mock:      func(prop *MockServiceProp) {},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorEmptyCredential, "Login.UserRepo.GetUsersByUsername"),
		},
		{
			name: "should validate user not found correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "password",
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUsername(ctx, prop.DB, "username").Once().Return([]user_entity.User{}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, ""),
		},
		{
			name: "should validate wrong password correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "wrongPassword",
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUsername(ctx, prop.DB, "username").Once().Return([]user_entity.User{{UserID: 1, Password: hashedPwd}}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, ""),
		},
		{
			name: "should return token correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "password",
			},
			mock: func(prop *MockServiceProp) {
				prop.UserRepo.EXPECT().GetUsersByUsername(ctx, prop.DB, "username").Once().Return([]user_entity.User{{UserID: 1, Password: hashedPwd}}, nil)
				prop.TokenMaker.EXPECT().CreateToken(1, token.AccessTokenDuration, token.AccessToken).Once().Return("accessToken", mockACTokenPayload, nil)
				prop.TokenMaker.EXPECT().CreateToken(1, token.RefreshTokenDuration, token.RefreshToken).Once().Return("refreshToken", mockRFTokenPayload, nil)
				prop.DB.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(prop.TX, nil)
				prop.SessionRepo.EXPECT().CreateSession(ctx, prop.TX, auth_entity.NewSession(mockRFTokenPayload.ID, mockRFTokenPayload.UserID, "refreshToken", mockRFTokenPayload.ExpiresAt)).Once().Return(nil)
				prop.TX.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockAuthService(t)
			tt.mock(prop)
			acToken, rfToken, err := sv.Login(ctx, tt.cred.Username, tt.cred.Password)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
				return
			}
			assert.NotEmpty(t, acToken)
			assert.NotEmpty(t, rfToken)
		})
	}
}
