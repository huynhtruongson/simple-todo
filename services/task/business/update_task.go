package task_biz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_entity "github.com/sondev/todo-list/services/task/entity"
)

type UpdateTaskBiz struct {
	DB lib.DB
	TaskRepo
	UserRepo
}

func NewUpdateTaskBiz(db lib.DB, taskRepo TaskRepo, userRepo UserRepo) *CreateTaskBiz {
	return &CreateTaskBiz{
		DB:       db,
		TaskRepo: taskRepo,
		UserRepo: userRepo,
	}
}

func (biz CreateTaskBiz) UpdateTask(ctx context.Context, task task_entity.Task) error {
	tasks, err := biz.TaskRepo.GetTasksByIds(ctx, biz.DB, []int{task.TaskID})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "UpdateTask.TaskRepo.GetTasksByIds")
	}
	if len(tasks) != 1 {
		return common.NewInvalidRequestError(err, task_entity.ErrorTaskNotFound, "UpdateTask")
	}
	if err := biz.ValidateTask(ctx, task); err != nil {
		return err
	}
	if err := lib.ExecTX(ctx, biz.DB, func(ctx context.Context, tx pgx.Tx) error {
		err := biz.TaskRepo.UpdateTask(ctx, tx, task)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "UpdateTask.TaskRepo.UpdateTask")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
