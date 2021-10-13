package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetEnvAsStringOrFallback(t *testing.T) {
	t.Run("With success case: return value", func(tc *testing.T) {
		_ = os.Setenv("TEST_GET_ENV", "PASS")
		value := GetEnvAsStringOrFallback("TEST_GET_ENV", "")
		assert.Equal(tc, "PASS", value)
	})

	t.Run("With success case: return blank value", func(tc *testing.T) {
		value := GetEnvAsStringOrFallback("TEST_GET_ENV", "")
		assert.Equal(tc, "", value)
	})
}
