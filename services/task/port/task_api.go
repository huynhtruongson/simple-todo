package task_port

import (
	"context"
	"net/http"
	"strconv"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/interceptor"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
)

type TaskService interface {
	CreateTask(ctx context.Context, task task_entity.Task) (int, error)
	DeleteTask(ctx context.Context, userID, taskID int) error
	ListTask(ctx context.Context, paging common.Paging, filter common.Filter) (task_entity.TasksWithPaging, error)
	UpdateTask(ctx context.Context, task task_entity.Task) error
}

type TaskAPI struct {
	TaskService
}

func NewTaskAPIService(taskService TaskService) *TaskAPI {
	return &TaskAPI{
		TaskService: taskService,
	}
}

func (api *TaskAPI) ListTask(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "ListTask Get URL Query Param")
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "ListTask Get URL Query Param")
	}

	payload := ctx.MustGet(interceptor.AuthorizationPayloadKey).(token.TokenPayload)

	tasks, err := api.TaskService.ListTask(ctx, common.Paging{Page: page, Limit: limit}, common.Filter{UserID: payload.UserID})
	if err != nil {
		code := http.StatusBadRequest
		appErr, ok := err.(*common.AppError)
		if ok {
			code = appErr.Code
		}
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(tasks))
}

func (api *TaskAPI) CreateTask(ctx *gin.Context) {
	var task task_entity.Task
	if err := ctx.ShouldBind(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "CreateTask Bind Json"))
		return
	}
	payload := ctx.MustGet(interceptor.AuthorizationPayloadKey).(token.TokenPayload)
	task.UserID = payload.UserID
	taskID, err := api.TaskService.CreateTask(ctx, task)
	if err != nil {
		code := http.StatusBadRequest
		appErr, ok := err.(*common.AppError)
		if ok {
			code = appErr.Code
		}
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(taskID))

}

func (api *TaskAPI) DeleteTask(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "DeleteTask Get URL Param")
	}

	payload := ctx.MustGet(interceptor.AuthorizationPayloadKey).(token.TokenPayload)

	err = api.TaskService.DeleteTask(ctx, payload.UserID, taskId)
	if err != nil {
		code := http.StatusBadRequest
		appErr, ok := err.(*common.AppError)
		if ok {
			code = appErr.Code
		}
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(true))

}

func (api *TaskAPI) UpdateTask(ctx *gin.Context) {
	var task task_entity.Task

	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "UpdateTask Get URL Param")
	}

	if err := ctx.ShouldBind(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "UpdateTask Bind Json"))
		return
	}
	payload := ctx.MustGet(interceptor.AuthorizationPayloadKey).(token.TokenPayload)
	task.TaskID = taskId
	task.UserID = payload.UserID

	err = api.TaskService.UpdateTask(ctx, task)
	if err != nil {
		code := http.StatusBadRequest
		appErr, ok := err.(*common.AppError)
		if ok {
			code = appErr.Code
		}
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(true))

}
