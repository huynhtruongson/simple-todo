package interceptor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationPayloadKey = "authorization_payload"
	AuthorizationTypeBearer = "bearer"

	UnAuthorizedMessage        = "Unauthorized"
	EmptyAuthHeaderMessage     = "Missing authorization header"
	InvalidAuthHeaderMessage   = "Invalid authorization header format"
	UnsupportedAuthTypeMessage = "Unsupported authorization type"
)

func AuthMiddleware(tokenMaker token.TokenMaker) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(nil, EmptyAuthHeaderMessage, ""))
			ctx.Abort()
			return
		}
		tokenFields := strings.Fields(authHeader)
		if len(tokenFields) < 2 {
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(nil, InvalidAuthHeaderMessage, ""))
			ctx.Abort()
			return
		}
		if strings.ToLower(tokenFields[0]) != AuthorizationTypeBearer {
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(nil, fmt.Sprintf(UnsupportedAuthTypeMessage+" %s", tokenFields[0]), ""))
			ctx.Abort()
			return
		}

		accessToken := tokenFields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(err, UnAuthorizedMessage, "Verify token failed"))
			ctx.Abort()
			return
		}
		if payload.Type != token.AccessToken {
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(nil, UnAuthorizedMessage, "invalid access token"))
			ctx.Abort()
			return
		}
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
