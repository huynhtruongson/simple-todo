package user_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, db lib.QueryExecer, user user_entity.User) (int, error)
	GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error)
}
type UserService struct {
	DB lib.DB
	UserRepo
}

func NewUserService(db lib.DB, userRepo UserRepo) *UserService {
	return &UserService{
		DB:       db,
		UserRepo: userRepo,
	}
}
