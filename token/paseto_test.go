package token

import (
	"testing"
	"time"

	"github.com/huynhtruongson/simple-todo/common"
	"github.com/huynhtruongson/simple-todo/utils"

	"github.com/stretchr/testify/assert"
)

func TestPasetoTokenMaker(t *testing.T) {
	t.Run("should create and verify token correctly", func(t *testing.T) {
		key := utils.RandomString(32)
		duration := time.Minute
		now := time.Now().UTC()
		expiredAt := now.Add(duration)
		userID := 1

		pasetoMaker, err := NewPasetoMaker(key)
		assert.NoError(t, err)
		assert.NotZero(t, pasetoMaker)

		token, tokenPayload, err := pasetoMaker.CreateToken(userID, duration, AccessToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, tokenPayload)

		payload, err := pasetoMaker.VerifyToken(token)
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

		pasetoMaker, err := NewPasetoMaker(key)
		assert.NoError(t, err)
		assert.NotEmpty(t, pasetoMaker)

		token, tokenPayload, err := pasetoMaker.CreateToken(userID, duration, AccessToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, tokenPayload)

		_, err = pasetoMaker.VerifyToken(token)
		assert.EqualError(t, err, common.ExpiredTokenErrorMessage.Error())
	})

	t.Run("should throw error when key size length is not equal 32", func(t *testing.T) {
		key := utils.RandomString(30)

		_, err := NewPasetoMaker(key)
		assert.EqualError(t, err, common.InvalidKeySizeErrorMessage.Error())
	})
}
