package interceptor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/huynhtruongson/simple-todo/token"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"github.com/rs/zerolog/log"
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

func UnaryServerAuthInterceptor(ai *AuthInterceptor) grpc.UnaryServerInterceptor {
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
func UnaryServerLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		now := time.Now()
		result, err := handler(ctx, req)
		duration := time.Since(now)
		statusCode := codes.Unknown
		if st, ok := status.FromError(err); ok {
			statusCode = st.Code()
		}
		logger := log.Info()
		if err != nil {
			logger = log.Error().Err(err)
		}
		logger.
			Str("protocol", "grpc").
			Str("method", info.FullMethod).
			Int("status_code", int(statusCode)).
			Str("status_text", statusCode.String()).
			Dur("duration", duration).
			Msg("receive a grpc request")
		return result, err
	}
}

func (ai *AuthInterceptor) AuthHandler(ctx context.Context) (context.Context, error) {
	authHeader := metadata.ExtractIncoming(ctx).Get(AuthorizationHeaderKey)
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, EmptyAuthHeaderMessage.Error())
	}
	tokenFields := strings.Fields(authHeader)
	if len(tokenFields) < 2 {
		return nil, status.Error(codes.Unauthenticated, InvalidAuthHeaderMessage.Error())
	}
	if strings.ToLower(tokenFields[0]) != AuthorizationTypeBearer {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf(UnsupportedAuthTypeMessage.Error()+" %s", tokenFields[0]))
	}
	accessToken := tokenFields[1]
	payload, err := ai.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, UnAuthorizedMessage.Error())
	}
	if payload.Type != token.AccessToken {
		return nil, status.Error(codes.Unauthenticated, UnAuthorizedMessage.Error())
	}
	return context.WithValue(ctx, AuthTokenPayload, payload), nil
}
