package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	t.Run("should hash password correctly", func(t *testing.T) {
		pwd := "password123"
		hashedPwd, _ := HashPassword(pwd)
		assert.NotEmpty(t, hashedPwd)
	})
	t.Run("should generate different hashed password when 2 raw password are the same", func(t *testing.T) {
		pwd1 := "password123"
		pwd2 := "password123"
		hashedPwd1, _ := HashPassword(pwd1)
		hashedPwd2, _ := HashPassword(pwd2)
		assert.NotEmpty(t, hashedPwd1)
		assert.NotEmpty(t, hashedPwd2)
		assert.NotEqual(t, hashedPwd1, hashedPwd2)
	})
	t.Run("should check raw password and hashed password are the same", func(t *testing.T) {
		pwd := "password123"
		hashedPwd, _ := HashPassword(pwd)
		assert.Nil(t, CheckPassword(pwd, hashedPwd))
	})
	t.Run("should check raw password and hashed password are not the same", func(t *testing.T) {
		pwd1 := "password123"
		hashedPwd1, _ := HashPassword(pwd1)
		pwd2 := "password456"
		assert.EqualError(t, CheckPassword(pwd2, hashedPwd1), bcrypt.ErrMismatchedHashAndPassword.Error())
	})
}
