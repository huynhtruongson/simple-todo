package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUpdatePlaceHolders(t *testing.T) {
	t.Run("should generate update placeholders correctly", func(t *testing.T) {
		fields := []string{"name", "age", "phone"}
		expected := `name=$1, age=$2, phone=$4`
		placeholderStr := GenerateUpdatePlaceHolders(fields)
		assert.Equal(t, expected, placeholderStr)
	})
}

func TestGeneratePlaceHolders(t *testing.T) {
	t.Run("should generate placeholders correctly", func(t *testing.T) {
		expected := `$1,$2,$3`
		placeholderStr := GeneratePlaceHolders(3)
		assert.Equal(t, expected, placeholderStr)
	})
}

func TestRandomInt(t *testing.T) {
	t.Run("should generate random int correctly", func(t *testing.T) {
		randInt := RandomInt(0, 10)
		assert.GreaterOrEqual(t, randInt, 0)
		assert.LessOrEqual(t, randInt, 10)
	})
}

func TestRandomString(t *testing.T) {
	t.Run("should generate random string correctly", func(t *testing.T) {
		randStr := RandomString(10)
		assert.Equal(t, len(randStr), 10)
	})
}
