package push

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnv_Convert(t *testing.T) {
	t.Run("With success case: return a map", func(tc *testing.T) {
		e := NewEnv()
		content, err := e.Convert([]byte(`export test=true
data=1,2,3
export str='this is test string'
`))
		assert.NoError(tc, err)
		expected := map[string]interface{}{
			"test": "true",
			"data": "1,2,3",
			"str":  "this is test string",
		}
		assert.Equal(tc, len(expected), len(content))
		for k, v := range expected {
			assert.Contains(tc, content, k)
			assert.Equal(tc, v, content[k])
		}
	})

	t.Run("With error case: src invalid", func(tc *testing.T) {
		e := NewTfvars()
		_, err := e.Convert([]byte(`export name = "test"`))
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Push: Parse content to JSON failed")

	})
}
