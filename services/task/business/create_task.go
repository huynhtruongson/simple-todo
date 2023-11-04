package task_biz

import (
	"context"
	"strings"

	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	user_entity "github.com/sondev/todo-list/services/user/entity"

	"github.com/jackc/pgx/v5"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) (int, error)
	GetTasksByIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]task_entity.Task, error)
	UpdateTask(ctx context.Context, db lib.QueryExecer, task task_entity.Task) error
	DeleteTask(ctx context.Context, db lib.QueryExecer, id int) error
	CountTask(ctx context.Context, db lib.QueryExecer) (int, error)
	GetTasksWithFilter(ctx context.Context, db lib.QueryExecer, limit, offset int) ([]task_entity.Task, error)
}
type UserRepo interface {
	GetUsersByUserIds(ctx context.Context, db lib.QueryExecer, ids []int) ([]user_entity.User, error)
}
type CreateTaskBiz struct {
	DB lib.DB
	TaskRepo
	UserRepo
}

func NewCreateTaskBiz(db lib.DB, taskRepo TaskRepo, userRepo UserRepo) *CreateTaskBiz {
	return &CreateTaskBiz{
		DB:       db,
		TaskRepo: taskRepo,
		UserRepo: userRepo,
	}
}

func (biz CreateTaskBiz) CreateTask(ctx context.Context, task task_entity.Task) (int, error) {
	if err := biz.ValidateTask(ctx, task); err != nil {
		return 0, err
	}
	var taskID int
	if err := lib.ExecTX(ctx, biz.DB, func(ctx context.Context, tx pgx.Tx) error {
		id, err := biz.TaskRepo.CreateTask(ctx, tx, task)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "TaskRepo.CreateTask")
		}
		taskID = id
		return nil
	}); err != nil {
		return taskID, err
	}
	return taskID, nil
}

func (biz CreateTaskBiz) ValidateTask(ctx context.Context, task task_entity.Task) error {
	switch {
	case strings.TrimSpace(task.Title) == "":
		return common.NewInvalidRequestError(nil, task_entity.ErrorTitleIsEmpty, "ValidateTask")
	case task.UserID == 0:
		return common.NewInvalidRequestError(nil, task_entity.ErrorUserIsEmpty, "ValidateTask")
	case task.Status > 2:
		return common.NewInvalidRequestError(nil, task_entity.ErrorInvalidStatus, "ValidateTask")
	}
	users, err := biz.UserRepo.GetUsersByUserIds(ctx, biz.DB, []int{task.UserID})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateTask.UserRepo.GetUsersByUserIds")
	}
	if len(users) == 0 {
		return common.NewInvalidRequestError(nil, task_entity.ErrorUserNotFound, "ValidateTask")
	}
	return nil
}
