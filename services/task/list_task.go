package task_service

import (
	"context"
	"math"

	"github.com/huynhtruongson/simple-todo/common"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
)

func (s *TaskService) ListTask(ctx context.Context, paging common.Paging, filter common.Filter) (task_entity.TasksWithPaging, error) {
	tasksWithPaging := task_entity.TasksWithPaging{}
	totalTasks, err := s.TaskRepo.CountTask(ctx, s.DB, filter.UserID)
	paging.Total = totalTasks
	if err != nil {
		return tasksWithPaging, common.NewInternalError(err, common.InternalErrorMessage, "ListTask.TaskRepo.CountTask")
	}
	if paging.Limit < 1 || paging.Limit > 20 {
		paging.Limit = 10
	}
	if paging.Page < 1 {
		paging.Page = 1
	}
	totalPage := math.Ceil(float64(totalTasks) / float64(paging.Limit))
	if paging.Page > int(totalPage) {
		paging.Page = int(totalPage)
	}
	offset := paging.Limit * (paging.Page - 1)

	tasks, err := s.TaskRepo.GetTasksWithFilter(ctx, s.DB, filter.UserID, paging.Limit, offset)
	if err != nil {
		return tasksWithPaging, common.NewInternalError(err, common.InternalErrorMessage, "ListTask.TaskRepo.GetTasksWithFilter")
	}
	tasksWithPaging.Tasks = tasks
	tasksWithPaging.Paging = paging

	return tasksWithPaging, nil
}
