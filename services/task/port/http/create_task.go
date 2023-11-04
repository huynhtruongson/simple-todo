package task_port

import (
	"net/http"

	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	task_biz "github.com/sondev/todo-list/services/task/business"
	task_entity "github.com/sondev/todo-list/services/task/entity"
	task_repo "github.com/sondev/todo-list/services/task/repository"
	user_repo "github.com/sondev/todo-list/services/user/repository"

	"github.com/gin-gonic/gin"
)

func CreateTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var task task_entity.Task
		if err := ctx.ShouldBind(&task); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "CreateTask Bind Json"))
			return
		}

		biz := task_biz.NewCreateTaskBiz(db, task_repo.NewTaskRepo(), user_repo.NewUserRepo())

		taskID, err := biz.CreateTask(ctx, task)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(taskID))
	}
}
