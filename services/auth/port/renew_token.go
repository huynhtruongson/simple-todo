package auth_port

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huynhtruongson/simple-todo/common"
)

type RenewTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RenewTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary		Renew token
// @Description	renew token
// @Tags         auth
// @Accept			json
// @Produce		json
// @Param			refresh_token	body	RenewTokenRequest	true	"user's refresh token"
// @Success		200		{object}	RenewTokenResponse
// @Failure		500	{object}	common.AppError
// @Failure		400	{object}	common.AppError
// @Router			/auth/renew-token [post]
func (api *AuthAPI) RenewToken(ctx *gin.Context) {
	var req RenewTokenRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, common.NewInvalidRequestError(err, common.InvalidRequestErrorMessage, "Login Bind Json"))
		return
	}
	acToken, err := api.AuthService.RenewToken(ctx, req.RefreshToken)
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
