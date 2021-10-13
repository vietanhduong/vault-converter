package hcl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHcl_ToJSON(t *testing.T) {
	t.Run("With success case: return a map", func(tc *testing.T) {
		h := New()
		values, err := h.ToJSON([]byte(`content = true`))
		assert.NoError(tc, err)
		assert.Equal(tc, values["content"].(bool), true)
	})

	t.Run("With error case", func(tc *testing.T) {
		h := New()
		_, err := h.ToJSON([]byte(`content = tru`))
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Unknown token")
	})
}

func TestHcl_ToHCL(t *testing.T) {
	t.Run("With success case: return byte array", func(tc *testing.T) {
		h := New()
		input := map[string]interface{}{
			"content": "string",
		}
		content, err := h.ToHCL(input)
		assert.NoError(tc, err)
		assert.Equal(tc, `content = "string"

`, string(content))
	})
}
