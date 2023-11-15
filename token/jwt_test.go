package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTTokenMaker(t *testing.T) {
	t.Run("should create and verify token correctly", func(t *testing.T) {
		key := utils.RandomString(32)
		duration := time.Minute
		now := time.Now().UTC()
		expiredAt := now.Add(duration)
		userID := 1

		jwtMaker, err := NewJWTMaker(key)
		assert.NoError(t, err)
		assert.NotZero(t, jwtMaker)

		token, tokenPayload, err := jwtMaker.CreateToken(userID, duration, AccessToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, tokenPayload)

		payload, err := jwtMaker.VerifyToken(token)
		assert.NoError(t, err)
		assert.NotZero(t, payload.ID)
		assert.Equal(t, userID, payload.UserID)
		assert.WithinDuration(t, now, payload.IssuedAt, time.Second)
		assert.WithinDuration(t, expiredAt, payload.ExpiresAt, time.Second)
	})

	t.Run("should verify expired token correctly", func(t *testing.T) {
		key := utils.RandomString(32)
		duration := -time.Minute
		userID := 1

		jwtMaker, err := NewJWTMaker(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, jwtMaker)

		token, tokenPayload, err := jwtMaker.CreateToken(userID, duration, AccessToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, tokenPayload)

		_, err = jwtMaker.VerifyToken(token)
		assert.ErrorContains(t, err, jwt.ErrTokenExpired.Error())
	})

	t.Run("should throw error when key size length is not equal 32", func(t *testing.T) {
		key := utils.RandomString(30)

		_, err := NewPasetoMaker(key)
		assert.EqualError(t, err, common.InvalidKeySizeErrorMessage.Error())
	})
}
