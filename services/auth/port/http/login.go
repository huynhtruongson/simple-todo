package auth_port

import (
	"net/http"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	auth_biz "github.com/huynhtruongson/simple-todo/services/auth/business"
	auth_entity "github.com/huynhtruongson/simple-todo/services/auth/entity"
	auth_repo "github.com/huynhtruongson/simple-todo/services/auth/repository"
	user_repo "github.com/huynhtruongson/simple-todo/services/user/repository"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
)

func Login(db lib.DB, tokenMaker token.TokenMaker) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var credential auth_entity.Credential

		if err := ctx.ShouldBind(&credential); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "Login Bind Json"))
			return
		}
		biz := auth_biz.NewLoginBiz(db, tokenMaker, user_repo.NewUserRepo(), auth_repo.NewSessionRepo())
		acToken, rfToken, err := biz.Login(ctx, credential.Username, credential.Password)
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
}
