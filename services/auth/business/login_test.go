package auth_biz

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

func TestLoginBiz_Login(t *testing.T) {
	ctx := context.Background()
	userRepo := mock_repo.NewUserRepo(t)
	sessionRepo := mock_repo.NewSessionRepo(t)
	tokenMaker := mock_token.NewTokenMaker(t)

	db := mock_db.NewDB(t)
	tx := mock_db.NewTx(t)
	hashedPwd, err := utils.HashPassword("password")
	assert.NoError(t, err)
	mockACTokenPayload, err := token.NewTokenPayload(1, token.AccessTokenDuration, token.AccessToken)
	assert.NoError(t, err)
	mockRFTokenPayload, err := token.NewTokenPayload(1, token.RefreshTokenDuration, token.RefreshToken)
	assert.NoError(t, err)
	tests := []struct {
		name      string
		cred      auth_entity.Credential
		mock      func()
		expectErr *common.AppError
	}{
		{
			name: "should validate empty username correctly",
			cred: auth_entity.Credential{
				Username: "",
				Password: "password",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorEmptyCredential, "Login.UserRepo.GetUsersByUsername"),
		},
		{
			name: "should validate empty password correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "",
			},
			mock:      func() {},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorEmptyCredential, "Login.UserRepo.GetUsersByUsername"),
		},
		{
			name: "should validate user not found correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "password",
			},
			mock: func() {
				userRepo.EXPECT().GetUsersByUsername(ctx, db, "username").Once().Return([]user_entity.User{}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, ""),
		},
		{
			name: "should validate wrong password correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "wrongPassword",
			},
			mock: func() {
				userRepo.EXPECT().GetUsersByUsername(ctx, db, "username").Once().Return([]user_entity.User{{UserID: 1, Password: hashedPwd}}, nil)
			},
			expectErr: common.NewInvalidRequestError(nil, auth_entity.ErrorInvalidCredential, ""),
		},
		{
			name: "should return token correctly",
			cred: auth_entity.Credential{
				Username: "username",
				Password: "password",
			},
			mock: func() {
				userRepo.EXPECT().GetUsersByUsername(ctx, db, "username").Once().Return([]user_entity.User{{UserID: 1, Password: hashedPwd}}, nil)
				tokenMaker.EXPECT().CreateToken(1, token.AccessTokenDuration, token.AccessToken).Once().Return("accessToken", mockACTokenPayload, nil)
				tokenMaker.EXPECT().CreateToken(1, token.RefreshTokenDuration, token.RefreshToken).Once().Return("refreshToken", mockRFTokenPayload, nil)
				db.EXPECT().BeginTx(ctx, mock.Anything).Once().Return(tx, nil)
				sessionRepo.EXPECT().CreateSession(ctx, tx, auth_entity.NewSession(mockRFTokenPayload.ID, mockRFTokenPayload.UserID, "refreshToken", mockRFTokenPayload.ExpiresAt)).Once().Return(nil)
				tx.EXPECT().Commit(ctx).Once().Return(nil)
			},
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewLoginBiz(db, tokenMaker, userRepo, sessionRepo)
			tt.mock()
			acToken, rfToken, err := biz.Login(ctx, tt.cred.Username, tt.cred.Password)
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
