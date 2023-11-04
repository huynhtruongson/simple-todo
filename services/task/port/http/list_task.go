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

func ListTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		page, err := strconv.Atoi(ctx.Query("page"))
		if err != nil {
			common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "ListTask Get URL Query Param")
		}
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "ListTask Get URL Query Param")
		}

		biz := task_biz.NewListTaskBiz(db, task_repo.NewTaskRepo())

		tasks, err := biz.ListTask(ctx, common.Paging{Page: page, Limit: limit})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(tasks))
	}
}
