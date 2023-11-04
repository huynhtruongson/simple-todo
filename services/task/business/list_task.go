package task_biz

import (
	"context"
	"math"

	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_entity "github.com/sondev/todo-list/services/task/entity"
)

type TasksWithPaging struct {
	Tasks  []task_entity.Task `json:"tasks"`
	Paging common.Paging      `json:"pagination"`
}

type ListTaskBiz struct {
	DB lib.DB
	TaskRepo
}

func NewListTaskBiz(db lib.DB, taskRepo TaskRepo) *ListTaskBiz {
	return &ListTaskBiz{
		DB:       db,
		TaskRepo: taskRepo,
	}
}

func (biz ListTaskBiz) ListTask(ctx context.Context, paging common.Paging) (TasksWithPaging, error) {
	tasksWithPaging := TasksWithPaging{}
	totalTasks, err := biz.TaskRepo.CountTask(ctx, biz.DB)
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

	tasks, err := biz.TaskRepo.GetTasksWithFilter(ctx, biz.DB, paging.Limit, offset)
	if err != nil {
		return tasksWithPaging, common.NewInternalError(err, common.InternalErrorMessage, "ListTask.TaskRepo.GetTasksWithFilter")
	}
	tasksWithPaging.Tasks = tasks
	tasksWithPaging.Paging = paging

	return tasksWithPaging, nil
}
