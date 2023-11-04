package user_port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sondev/todo-list/common"
	"github.com/sondev/todo-list/lib"
	user_biz "github.com/sondev/todo-list/services/user/business"
	user_entity "github.com/sondev/todo-list/services/user/entity"
	user_repo "github.com/sondev/todo-list/services/user/repository"
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
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		ctx.JSON(http.StatusOK, common.NewSimpleSuccessResponse(userID))
	}
}
