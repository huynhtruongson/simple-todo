package auth_service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/huynhtruongson/simple-todo/interceptor"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/stretchr/testify/assert"

	"github.com/huynhtruongson/simple-todo/common"
)

func TestRenewTokenBiz_RenewToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockPayload, err := token.NewTokenPayload(1, token.AccessTokenDuration, token.AccessToken)
	assert.NoError(t, err)
	tests := []struct {
		name         string
		token        string
		mock         func(prop *MockServiceProp)
		expectErr    *common.AppError
		expectOutput string
	}{
		{
			name:  "should validate not found session correctly",
			token: "token",
			mock: func(prop *MockServiceProp) {
				prop.TokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				prop.SessionRepo.EXPECT().GetSessionByIds(ctx, prop.DB, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(interceptor.UnAuthorizedMessage, interceptor.UnAuthorizedMessage.Error(), ""),
		},
		{
			name:  "should validate mismatch refresh token correctly",
			token: "wrong_token",
			mock: func(prop *MockServiceProp) {
				prop.TokenMaker.EXPECT().VerifyToken("wrong_token").Once().Return(mockPayload, nil)
				prop.SessionRepo.EXPECT().GetSessionByIds(ctx, prop.DB, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token"}}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(interceptor.UnAuthorizedMessage, interceptor.UnAuthorizedMessage.Error(), ""),
		},
		{
			name:  "should validate mismatch refresh token is blocked",
			token: "token",
			mock: func(prop *MockServiceProp) {
				prop.TokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				prop.SessionRepo.EXPECT().GetSessionByIds(ctx, prop.DB, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token", IsBlocked: true}}, nil)
			},
			expectErr: common.NewUnAuthorizedRequestError(interceptor.UnAuthorizedMessage, interceptor.UnAuthorizedMessage.Error(), ""),
		},
		{
			name:  "should validate mismatch refresh token is blocked",
			token: "token",
			mock: func(prop *MockServiceProp) {
				prop.TokenMaker.EXPECT().VerifyToken("token").Once().Return(mockPayload, nil)
				prop.SessionRepo.EXPECT().GetSessionByIds(ctx, prop.DB, uuid.UUIDs{mockPayload.ID}).Once().Return([]auth_entity.Session{{RefreshToken: "token", UserID: 1}}, nil)
				prop.TokenMaker.EXPECT().CreateToken(1, token.AccessTokenDuration, token.AccessToken).Once().Return("new_access_token", token.TokenPayload{}, nil)
			},
			expectErr:    nil,
			expectOutput: "new_access_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sv, prop := NewMockAuthService(t)
			tt.mock(prop)
			acToken, err := sv.RenewToken(ctx, tt.token)
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
