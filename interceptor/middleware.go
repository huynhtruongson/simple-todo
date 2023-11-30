package interceptor

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/token"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationPayloadKey = "authorization_payload"
	AuthorizationTypeBearer = "bearer"
)

var (
	UnAuthorizedMessage        = errors.New("Unauthorized")
	EmptyAuthHeaderMessage     = errors.New("Missing authorization header")
	InvalidAuthHeaderMessage   = errors.New("Invalid authorization header format")
	UnsupportedAuthTypeMessage = errors.New("Unsupported authorization type")
)

func AuthMiddleware(tokenMaker token.TokenMaker) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			err := common.NewUnAuthorizedRequestError(EmptyAuthHeaderMessage, EmptyAuthHeaderMessage.Error(), "")
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		tokenFields := strings.Fields(authHeader)
		if len(tokenFields) < 2 {
			err := common.NewUnAuthorizedRequestError(InvalidAuthHeaderMessage, InvalidAuthHeaderMessage.Error(), "")
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		if strings.ToLower(tokenFields[0]) != AuthorizationTypeBearer {
			err := common.NewUnAuthorizedRequestError(nil, fmt.Sprintf(UnsupportedAuthTypeMessage.Error()+" %s", tokenFields[0]), "")
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		accessToken := tokenFields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, common.NewUnAuthorizedRequestError(err, UnAuthorizedMessage.Error(), "Verify token failed"))
			ctx.Abort()
			return
		}
		if payload.Type != token.AccessToken {
			err := common.NewUnAuthorizedRequestError(UnAuthorizedMessage, UnAuthorizedMessage.Error(), "invalid access token")
			ctx.Error(err)
			ctx.JSON(http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}
		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func LoggingMiddleware(ctx *gin.Context) {
	now := time.Now()
	ctx.Next()
	logger := log.Info()
	if ctx.Writer.Status() != http.StatusOK {
		errs := []error{}
		for _, err := range ctx.Errors {
			errs = append(errs, err)
		}
		logger = log.Error().Errs("errors", errs)

	}
	logger.
		Str("protocol", ctx.Request.Proto).
		Str("method", ctx.Request.Method).
		Str("path", ctx.Request.URL.Path).
		Int("status_code", ctx.Writer.Status()).
		Str("status_text", http.StatusText(ctx.Writer.Status())).
		Dur("duration", time.Since(now)).
		Msg("receive a http request")
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func GatewayLoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, r)
		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)

		}
		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", time.Since(now)).
			Msg("receive a http request")
	})
}
