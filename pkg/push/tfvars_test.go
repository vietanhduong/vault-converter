package push

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTfvars_Convert(t *testing.T) {
	t.Run("With success case: return a map", func(tc *testing.T) {
		tf := NewTfvars()
		content, err := tf.Convert([]byte(`name = "test"
arr = [1,2,3]

dict = {

value = true

pass = true

content = "content"
}
`))
		assert.NoError(tc, err)
		dict := map[string]interface{}{
			"value":   true,
			"pass":    true,
			"content": "content",
		}
		expected := map[string]interface{}{
			"name": "test",
			"arr":  []interface{}{1, 2, 3},
			"dict": dict,
		}
		actualDict, ok := content["dict"].(map[string]interface{})

		assert.Equal(tc, true, ok)
		assert.Equal(tc, len(expected), len(content))
		assert.Equal(tc, len(expected["arr"].([]interface{})), len(content["arr"].([]interface{})))
		assert.Equal(tc,dict["value"], actualDict["value"])
	})

	t.Run("With error case: src invalid", func(tc *testing.T) {
		tf := NewTfvars()
		_, err := tf.Convert([]byte(`name = "test"
arr = [1,2,3]

dict = {

value = true

pass = true

content = "content"
`))
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Push: Parse content to JSON failed")
	})
}
