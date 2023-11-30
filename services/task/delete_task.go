package task_service

import (
	"context"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/jackc/pgx/v5"
)

func (s *TaskService) DeleteTask(ctx context.Context, userID, taskID int) error {
	tasks, err := s.TaskRepo.GetTasksByIds(ctx, s.DB, userID, []int{taskID})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "DeleteTask.TaskRepo.GetTasksByIds")
	}
	if len(tasks) != 1 {
		return common.NewInvalidRequestError(task_entity.ErrorTaskNotFound, task_entity.ErrorTaskNotFound.Error(), "DeleteTask")
	}
	if err := lib.ExecTX(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		err := s.TaskRepo.DeleteTask(ctx, tx, taskID)
		if err != nil {
			return common.NewInternalError(err, common.InternalErrorMessage, "DeleteTask.TaskRepo.DeleteTask")
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
