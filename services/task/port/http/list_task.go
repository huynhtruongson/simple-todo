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
		payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(token.TokenPayload)

		tasks, err := biz.ListTask(ctx, common.Paging{Page: page, Limit: limit}, common.Filter{UserID: payload.UserID})
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
}
