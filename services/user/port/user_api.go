package user_port

import (
	"context"
	"net/http"

	"github.com/huynhtruongson/simple-todo/common"
	user_entity "github.com/huynhtruongson/simple-todo/services/user/entity"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	CreateUser(ctx context.Context, user user_entity.User) (int, error)
}

type UserAPIService struct {
	UserService
}

func NewUserAPIService(userService UserService) *UserAPIService {
	return &UserAPIService{
		UserService: userService,
	}
}

func (sv *UserAPIService) CreateUser(ctx *gin.Context) {
	var user user_entity.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "CreateUser Bind Json"))
		return
	}

	userID, err := sv.UserService.CreateUser(ctx, user)
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
