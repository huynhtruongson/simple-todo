package interceptor

import (
	"context"
	"fmt"
	"strings"

	"github.com/huynhtruongson/simple-todo/token"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	TokenMaker token.TokenMaker
}

func NewAuthInterceptor(tokenMaker token.TokenMaker) *AuthInterceptor {
	return &AuthInterceptor{
		TokenMaker: tokenMaker,
	}
}

type ServiceAuthFuncOverride interface {
	AuthHandlerOverride(context.Context) (context.Context, error)
}

var AuthTokenPayload = struct{}{}

func UnaryServerInterceptor(ai *AuthInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthHandlerOverride(ctx)
		} else {
			newCtx, err = ai.AuthHandler(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (ai *AuthInterceptor) AuthHandler(ctx context.Context) (context.Context, error) {
	authHeader := metadata.ExtractIncoming(ctx).Get(AuthorizationHeaderKey)
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, EmptyAuthHeaderMessage)
	}
	tokenFields := strings.Fields(authHeader)
	if len(tokenFields) < 2 {
		return nil, status.Error(codes.Unauthenticated, InvalidAuthHeaderMessage)
	}
	if strings.ToLower(tokenFields[0]) != AuthorizationTypeBearer {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf(UnsupportedAuthTypeMessage+" %s", tokenFields[0]))
	}
	accessToken := tokenFields[1]
	payload, err := ai.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, UnAuthorizedMessage)
	}
	if payload.Type != token.AccessToken {
		return nil, status.Error(codes.Unauthenticated, UnAuthorizedMessage)
	}
	return context.WithValue(ctx, AuthTokenPayload, payload), nil
}
