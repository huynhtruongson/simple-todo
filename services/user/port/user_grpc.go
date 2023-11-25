package user_port

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/pb"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
)

type UserGRPCService struct {
	pb.UnimplementedUserServiceServer
	UserService UserService
}

func NewUserGRPCService(userService UserService) *UserGRPCService {
	return &UserGRPCService{
		UserService: userService,
	}
}

func (sv *UserGRPCService) AuthHandlerOverride(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (sv *UserGRPCService) CreateUser(ctx context.Context, userReq *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := toUser(userReq)
	userID, err := sv.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, common.MapAppErrorToGRPCError(err, "Create user error")
	}
	return &pb.CreateUserResponse{
		Data: int64(userID),
	}, nil
}

func toUser(user *pb.CreateUserRequest) user_entity.User {
	return user_entity.User{
		FullName: user.GetFullname(),
		Username: user.GetUsername(),
		Password: user.GetPassword(),
	}
}
