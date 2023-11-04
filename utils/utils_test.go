package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUpdatePlaceHolders(t *testing.T) {
	t.Run("should generate update placeholders correctly", func(t *testing.T) {
		fields := []string{"name", "age", "phone"}
		expected := `name=$1, age=$2, phone=$3`
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
