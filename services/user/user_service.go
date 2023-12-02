package user_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/lib"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	"github.com/huynhtruongson/simple-todo/worker"

	"github.com/hibiken/asynq"
)

type UserRepo interface {
	CreateUser(ctx context.Context, db lib.QueryExecer, user user_entity.User) (int, error)
	GetUsersByUsername(ctx context.Context, db lib.QueryExecer, username string) ([]user_entity.User, error)
	GetUsersByEmail(ctx context.Context, db lib.QueryExecer, email string) ([]user_entity.User, error)
}
type WorkerClient interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *worker.TaskSendVerifyEmailPayload, opts ...asynq.Option) error
}
type UserService struct {
	DB lib.DB
	UserRepo
	WorkerClient
}

func NewUserService(db lib.DB, workerClient WorkerClient, userRepo UserRepo) *UserService {
	return &UserService{
		DB:           db,
		UserRepo:     userRepo,
		WorkerClient: workerClient,
	}
}
