package interceptor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	grpcMetadata "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthHandler(t *testing.T) {
	t.Parallel()
	key := utils.RandomString(32)
	tokenMaker, err := token.NewPasetoMaker(key)
	assert.NoError(t, err)
	testCases := []struct {
		name       string
		setup      func() context.Context
		expectCode codes.Code
	}{
		{
			name: EmptyAuthHeaderMessage,
			setup: func() context.Context {
				md := grpcMetadata.Pairs("authorization", "")
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: codes.Unauthenticated,
		},
		{
			name: InvalidAuthHeaderMessage,
			setup: func() context.Context {
				md := grpcMetadata.Pairs("authorization", "token")
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: codes.Unauthenticated,
		},
		{
			name: UnsupportedAuthTypeMessage,
			setup: func() context.Context {
				md := grpcMetadata.Pairs("authorization", "bearer123 token")
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: codes.Unauthenticated,
		},
		{
			name: "token expired",
			setup: func() context.Context {
				token, _, err := tokenMaker.CreateToken(1, -time.Minute, token.AccessToken)
				assert.NoError(t, err)
				md := grpcMetadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: codes.Unauthenticated,
		},
		{
			name: "invalid access token",
			setup: func() context.Context {
				token, _, err := tokenMaker.CreateToken(1, time.Minute, token.RefreshToken)
				assert.NoError(t, err)
				md := grpcMetadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: codes.Unauthenticated,
		},
		{
			name: "OK",
			setup: func() context.Context {
				token, _, err := tokenMaker.CreateToken(1, time.Minute, token.AccessToken)
				assert.NoError(t, err)
				md := grpcMetadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token))
				return metadata.MD(md).ToIncoming(context.Background())
			},
			expectCode: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := NewAuthInterceptor(tokenMaker)
			ctx := tt.setup()
			ctxs, err := interceptor.AuthHandler(ctx)
			if tt.expectCode != 0 {
				assert.Equal(t, tt.expectCode, status.Code(err))
				return
			}
			assert.NotEmpty(t, ctxs)
			payload, ok := ctxs.Value(AuthTokenPayload).(token.TokenPayload)
			assert.True(t, ok)
			assert.NotEmpty(t, payload)
		})
	}
}
