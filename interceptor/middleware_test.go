package interceptor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huynhtruongson/simple-todo/token"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	t.Parallel()
	r := gin.Default()
	key := utils.RandomString(32)
	tokenMaker, err := token.NewPasetoMaker(key)
	assert.NoError(t, err)
	r.GET("/ping", AuthMiddleware(tokenMaker), func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{}) })

	testCases := []struct {
		name       string
		setup      func(*http.Request)
		expectCode int
	}{
		{
			name: EmptyAuthHeaderMessage.Error(),
			setup: func(req *http.Request) {
				req.Header.Set(AuthorizationHeaderKey, "")
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: InvalidAuthHeaderMessage.Error(),
			setup: func(req *http.Request) {
				req.Header.Set(AuthorizationHeaderKey, "token")
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: UnsupportedAuthTypeMessage.Error(),
			setup: func(req *http.Request) {
				req.Header.Set(AuthorizationHeaderKey, "bearer123 token")
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "token expired",
			setup: func(req *http.Request) {
				token, _, err := tokenMaker.CreateToken(1, -time.Minute, token.AccessToken)
				assert.NoError(t, err)
				req.Header.Set(AuthorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "invalid access token",
			setup: func(req *http.Request) {
				token, _, err := tokenMaker.CreateToken(1, time.Minute, token.RefreshToken)
				assert.NoError(t, err)
				req.Header.Set(AuthorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
			},
			expectCode: http.StatusUnauthorized,
		},
		{
			name: "OK",
			setup: func(req *http.Request) {
				token, _, err := tokenMaker.CreateToken(1, time.Minute, token.AccessToken)
				assert.NoError(t, err)
				req.Header.Set(AuthorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
			},
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/ping", nil)
			assert.NoError(t, err)
			tt.setup(req)
			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectCode)
		})
	}
}
