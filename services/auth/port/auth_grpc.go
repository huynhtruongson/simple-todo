package auth_port

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/pb"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc/peer"
)

type AuthGRPCService struct {
	pb.UnimplementedAuthServiceServer
	AuthService AuthService
}

const (
	UserAgentGRPCGatewayHeader = "grpcgateway-user-agent"
	UserAgentHeader            = "user-agent"
	xForwardedFor              = "x-forwarded-for"
)

func NewAuthGRPCService(authService AuthService) *AuthGRPCService {
	return &AuthGRPCService{
		AuthService: authService,
	}
}
func (sv *AuthGRPCService) AuthHandlerOverride(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (sv *AuthGRPCService) Login(ctx context.Context, cred *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInfo := auth_entity.LoginInfo{}
	md := metadata.ExtractIncoming(ctx)

	if userAgent := md.Get(UserAgentGRPCGatewayHeader); len(userAgent) > 0 {
		loginInfo.UserAgent = userAgent
	}
	if userAgent := md.Get(UserAgentHeader); len(userAgent) > 0 {
		loginInfo.UserAgent = userAgent
	}
	if clientIP := md.Get(xForwardedFor); len(clientIP) > 0 {
		loginInfo.ClientIP = clientIP
	}

	if p, ok := peer.FromContext(ctx); ok {
		loginInfo.ClientIP = p.Addr.String()
	}

	acToken, rfToken, err := sv.AuthService.Login(ctx, auth_entity.Credential{
		Username: cred.GetUsername(),
		Password: cred.GetPassword(),
	}, loginInfo)
	if err != nil {
		return nil, common.MapAppErrorToGRPCError(err, "Login error")
	}
	return &pb.LoginResponse{
		AccessToken:  acToken,
		RefreshToken: rfToken,
	}, nil
}
