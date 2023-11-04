package task_biz

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_entity "github.com/sondev/todo-list/services/task/entity"
)

type DeleteTaskBiz struct {
	DB lib.DB
	TaskRepo
}

func NewDeleteTaskBiz(db lib.DB, taskRepo TaskRepo) *CreateTaskBiz {
	return &CreateTaskBiz{
		DB:       db,
		TaskRepo: taskRepo,
	}
}

func (biz CreateTaskBiz) DeleteTask(ctx context.Context, id int) error {
	tasks, err := biz.TaskRepo.GetTasksByIds(ctx, biz.DB, []int{id})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "DeleteTask.TaskRepo.GetTasksByIds")
	}
	if len(tasks) != 1 {
		return common.NewInvalidRequestError(err, task_entity.ErrorTaskNotFound, "DeleteTask")
	}
	if err := lib.ExecTX(ctx, biz.DB, func(ctx context.Context, tx pgx.Tx) error {
		err := biz.TaskRepo.DeleteTask(ctx, tx, id)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "DeleteTask.TaskRepo.DeleteTask")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
