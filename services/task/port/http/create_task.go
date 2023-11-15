package task_port

import (
	"net/http"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	"github.com/huynhtruongson/simple-todo/middleware"
	task_biz "github.com/huynhtruongson/simple-todo/services/task/business"
	task_entity "github.com/huynhtruongson/simple-todo/services/task/entity"
	task_repo "github.com/huynhtruongson/simple-todo/services/task/repository"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
)

func CreateTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var task task_entity.Task
		if err := ctx.ShouldBind(&task); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "CreateTask Bind Json"))
			return
		}
		payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(token.TokenPayload)
		task.UserID = payload.UserID
		biz := task_biz.NewCreateTaskBiz(db, task_repo.NewTaskRepo(), user_repo.NewUserRepo())

		taskID, err := biz.CreateTask(ctx, task)
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
}
