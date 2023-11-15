package auth_port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/lib"
	auth_biz "github.com/huynhtruongson/simple-todo/services/auth/business"
	auth_repo "github.com/huynhtruongson/simple-todo/services/auth/repository"
	"github.com/huynhtruongson/simple-todo/token"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func RenewToken(db lib.DB, tokenMaker token.TokenMaker) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req RenewTokenRequest

		if err := ctx.ShouldBind(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "Login Bind Json"))
			return
		}
		biz := auth_biz.NewRenewTokenBiz(db, tokenMaker, auth_repo.NewSessionRepo())
		acToken, err := biz.RenewToken(ctx, req.RefreshToken)
		if err != nil {
			code := http.StatusBadRequest
			appErr, ok := err.(*common.AppError)
			if ok {
				code = appErr.Code
			}
			ctx.JSON(code, err)
			return
		}
		ctx.JSON(http.StatusOK, &RenewTokenResponse{AccessToken: acToken})
	}
}
