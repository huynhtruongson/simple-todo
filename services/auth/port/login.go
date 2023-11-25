package auth_port

import (
	"net/http"

	"github.com/huynhtruongson/simple-todo/common"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"

	"github.com/gin-gonic/gin"
)

func (api *AuthAPI) Login(ctx *gin.Context) {
	var credential auth_entity.Credential

	if err := ctx.ShouldBind(&credential); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "Login Bind Json"))
		return
	}

	acToken, rfToken, err := api.AuthService.Login(ctx, credential, auth_entity.LoginInfo{
		UserAgent: ctx.Request.UserAgent(),
		ClientIP:  ctx.ClientIP(),
	})
	if err != nil {
		code := http.StatusBadRequest
		appErr, ok := err.(*common.AppError)
		if ok {
			code = appErr.Code
		}
		ctx.JSON(code, err)
		return
	}
	ctx.JSON(http.StatusOK, auth_entity.NewLoginResponse(acToken, rfToken))

}
