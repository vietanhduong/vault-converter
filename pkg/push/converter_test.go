package push

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConverter(t *testing.T) {
	t.Run("With success case", func(tc *testing.T) {
		_, err := NewConverter(".tfvars")
		assert.NoError(tc, err)
	})

	t.Run("With error case: failure with unsupported format", func(tc *testing.T) {
		_, err := NewConverter(".exe")
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "`.exe` is not yet supported.")
	})
}