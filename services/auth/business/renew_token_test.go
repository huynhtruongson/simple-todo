package auth_biz

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/middleware"
	mock_repo "github.com/huynhtruongson/simple-todo/mocks/auth"
	mock_db "github.com/huynhtruongson/simple-todo/mocks/lib"
	mock_token "github.com/huynhtruongson/simple-todo/mocks/token"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/stretchr/testify/assert"

	"github.com/huynhtruongson/simple-todo/common"
)

func TestRenewTokenBiz_RenewToken(t *testing.T) {
	ctx := context.Background()
	sessionRepo := mock_repo.NewSessionRepo(t)
	tokenMaker := mock_token.NewTokenMaker(t)
	db := mock_db.NewDB(t)
	mockPayload, err := token.NewTokenPayload(1, token.AccessTokenDuration, token.AccessToken)
	assert.NoError(t, err)
	tests := []struct {
		name         string
		token        string
		mock         func()
		expectErr    *common.AppError
		expectOutput string
	}{
		{
			name:  "should validate not found session correctly",
			token: "token",
			mock: func() {
				tokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				sessionRepo.EXPECT().GetSessionByIds(ctx, db, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, ""),
		},
		{
			name:  "should validate mismatch refresh token correctly",
			token: "wrong_token",
			mock: func() {
				tokenMaker.EXPECT().VerifyToken("wrong_token").Once().Return(mockPayload, nil)
				sessionRepo.EXPECT().GetSessionByIds(ctx, db, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token"}}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, ""),
		},
		{
			name:  "should validate mismatch refresh token is blocked",
			token: "token",
			mock: func() {
				tokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				sessionRepo.EXPECT().GetSessionByIds(ctx, db, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token", IsBlocked: true}}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(err, middleware.UnAuthorizedMessage, ""),
		},
		{
			name:  "should validate mismatch refresh token is blocked",
			token: "token",
			mock: func() {
				tokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				sessionRepo.EXPECT().GetSessionByIds(ctx, db, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token", UserID: 1}}, nil)
				tokenMaker.EXPECT().CreateToken(1, token.AccessTokenDuration, token.AccessToken).Once().Return("new_access_token", token.TokenPayload{}, nil)
			},
			expectErr:    nil,
			expectOutput: "new_access_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			biz := NewRenewTokenBiz(db, tokenMaker, sessionRepo)
			tt.mock()
			acToken, err := biz.RenewToken(ctx, tt.token)
			if tt.expectErr != nil {
				assert.Equal(t, tt.expectErr.Code, err.(*common.AppError).Code)
				assert.Equal(t, tt.expectErr.Message, err.(*common.AppError).Message)
				return
			}
			assert.NotEmpty(t, acToken)
			assert.Equal(t, tt.expectOutput, acToken)
		})
	}
}
