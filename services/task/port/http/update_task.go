package task_port

import (
	"net/http"
	"strconv"

	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_biz "github.com/sondev/todo-list/services/task/business"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_repo "github.com/sondev/todo-list/services/task/repository"
	user_repo "github.com/sondev/todo-list/services/user/repository"

	"github.com/gin-gonic/gin"
)

func UpdateTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var task task_entity.Task

		taskId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "UpdateTask Get URL Param")
		}

		if err := ctx.ShouldBind(&task); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "UpdateTask Bind Json"))
			return
		}
		task.TaskID = taskId

		biz := task_biz.NewUpdateTaskBiz(db, task_repo.NewTaskRepo(), user_repo.NewUserRepo())

		err = biz.UpdateTask(ctx, task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(true))
	}
}
