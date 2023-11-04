package task_port

import (
	"net/http"
	"strconv"

	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_biz "github.com/sondev/todo-list/services/task/business"
	task_repo "github.com/sondev/todo-list/services/task/repository"

	"github.com/gin-gonic/gin"
)

func DeleteTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		taskId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "DeleteTask Get URL Param")
		}

		biz := task_biz.NewDeleteTaskBiz(db, task_repo.NewTaskRepo())

		err = biz.DeleteTask(ctx, taskId)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(true))
	}
}
