package user_port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	user_biz "github.com/huynhtruongson/simple-todo/services/user/business"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
)

func CreateUser(db lib.DB) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var user user_entity.User

		if err := ctx.ShouldBind(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "CreateUser Bind Json"))
			return
		}

		biz := user_biz.NewCreateUserBiz(db, user_repo.NewUserRepo())

		userID, err := biz.CreateUser(ctx, user)
		if err != nil {
			code := http.StatusBadRequest
			appErr, ok := err.(*common.AppError)
			if ok {
				code = appErr.Code
			}
			ctx.JSON(code, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(userID))
	}
}
