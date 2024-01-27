package task_service

import (
	"context"
	"strings"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"

	"github.com/jackc/pgx/v5"
)

func (s *TaskService) CreateTask(ctx context.Context, task task_entity.Task) (int, error) {
	if err := s.ValidateTask(ctx, task); err != nil {
		return 0, err
	}
	var taskID int
	if err := lib.ExecTX(ctx, s.DB, func(ctx context.Context, tx pgx.Tx) error {
		id, err := s.TaskRepo.CreateTask(ctx, tx, task)
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

func (s *TaskService) ValidateTask(ctx context.Context, task task_entity.Task) error {
	switch {
	case strings.TrimSpace(task.Title.String()) == "":
		return common.NewInvalidRequestError(task_entity.ErrorTitleIsEmpty, task_entity.ErrorTitleIsEmpty.Error(), "ValidateTask")
	case task.UserID.Int() == 0:
		return common.NewInvalidRequestError(task_entity.ErrorUserIsEmpty, task_entity.ErrorUserIsEmpty.Error(), "ValidateTask")
	case task.Status > 2 || task.Status < 0:
		return common.NewInvalidRequestError(task_entity.ErrorInvalidStatus, task_entity.ErrorInvalidStatus.Error(), "ValidateTask")
	}
	users, err := s.UserRepo.GetUsersByUserIds(ctx, s.DB, []int{task.UserID.Int()})
	if err != nil {
		return common.NewInternalError(err, common.InternalErrorMessage, "ValidateTask.UserRepo.GetUsersByUserIds")
	}
	if len(users) == 0 {
		return common.NewInvalidRequestError(task_entity.ErrorUserNotFound, task_entity.ErrorUserNotFound.Error(), "ValidateTask")
	}
	return nil
}
