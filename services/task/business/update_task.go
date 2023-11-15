package task_biz

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/jackc/pgx/v5"
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
	tasks, err := biz.TaskRepo.GetTasksByIds(ctx, biz.DB, task.UserID, []int{task.TaskID})
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
