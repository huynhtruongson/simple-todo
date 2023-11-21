package task_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) (int, error)
	GetTasksByIds(ctx context.Context, db lib.QueryExecer, userID int, taskID []int) ([]task_entity.Task, error)
	UpdateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) error
	DeleteTask(ctx context.Context, db lib.QueryExecer, taskID int) error
	CountTask(ctx context.Context, db lib.QueryExecer, userID int) (int, error)
	GetTasksWithFilter(ctx context.Context, db lib.QueryExecer, userID, limit, offset int) ([]task_entity.Task, error)
}
type UserRepo interface {
	GetUsersByUserIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]user_entity.User, error)
}

type TaskService struct {
	DB lib.DB
	TaskRepo
	UserRepo
}

func NewTaskService(db lib.DB, taskRepo TaskRepo, userRepo UserRepo) *TaskService {
	return &TaskService{
		DB:       db,
		TaskRepo: taskRepo,
		UserRepo: userRepo,
	}
}
