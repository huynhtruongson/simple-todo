package task_port

import (
	"net/http"
	"strconv"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	"github.com/huynhtruongson/simple-todo/middleware"
	task_biz "github.com/huynhtruongson/simple-todo/services/task/business"
	task_repo "github.com/huynhtruongson/simple-todo/services/task/repository"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
)

func DeleteTask(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		taskId, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "DeleteTask Get URL Param")
		}

		biz := task_biz.NewDeleteTaskBiz(db, task_repo.NewTaskRepo())
		payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(token.TokenPayload)

		err = biz.DeleteTask(ctx, payload.UserID, taskId)
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
