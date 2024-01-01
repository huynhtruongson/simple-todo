package task_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/jackc/pgx/v5"
)

func (s *TaskService) UpdateTask(ctx context.Context, task task_entity.Task) error {
	tasks, err := s.TaskRepo.GetTasksByIds(ctx, s.DB, task.UserID.Int(), []int{task.TaskID.Int()})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "UpdateTask.TaskRepo.GetTasksByIds")
	}
	if len(tasks) != 1 {
		return common.NewInvalidRequestError(task_entity.ErrorTaskNotFound, task_entity.ErrorTaskNotFound.Error(), "UpdateTask")
	}
	if err := s.ValidateTask(ctx, task); err != nil {
		return err
	}
	if err := lib.ExecTX(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		err := s.TaskRepo.UpdateTask(ctx, tx, task)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "UpdateTask.TaskRepo.UpdateTask")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
