package task_port

import (
	"net/http"
	"strconv"

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
		payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(token.TokenPayload)
		task.TaskID = taskId
		task.UserID = payload.UserID

		biz := task_biz.NewUpdateTaskBiz(db, task_repo.NewTaskRepo(), user_repo.NewUserRepo())

		err = biz.UpdateTask(ctx, task)
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
}
