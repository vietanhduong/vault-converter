package hcl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHcl_ToJSON(t *testing.T) {
	t.Run("With success case: return a map", func(tc *testing.T) {
		h := New()
		content := `
content = true
list = [1,2,3]
advance_list = [
{
	username = "admin"
	password = "password"
	sensitive = true
},
{
	username = "tester"
	password = "password"
	sensitive = true
}
]
map = {
content = "1"
data = 2
}
value = "string"
number = -1
pass = null
`
		values, err := h.ToJSON([]byte(content))
		assert.NoError(tc, err)
		assert.Equal(tc, values["content"].(bool), true)

		list, ok := values["list"].([]interface{})
		assert.True(tc, ok)
		assert.Equal(tc, 3, len(list))

		advList, ok := values["advance_list"].([]interface{})
		assert.True(tc, ok)
		assert.Equal(tc, 2, len(advList))

		advListFirstElement, ok := advList[0].(map[string]interface{})
		assert.True(tc, ok)
		assert.True(tc, advListFirstElement["sensitive"].(bool))

		assert.Equal(tc, -1, values["number"].(int))
	})

	t.Run("With error case: invalid expression", func(tc *testing.T) {
		h := New()
		_, err := h.ToJSON([]byte(`content = `))
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Invalid expression")
	})

	t.Run("With error case", func(tc *testing.T) {
		h := New()
		_, err := h.ToJSON([]byte(`content = tru`))
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Parse attributes form HCL failed")
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
