package pull

import (
	"github.com/stretchr/testify/assert"
	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
	"os"
	"testing"
)

func TestTfvars_Convert(t *testing.T) {
	dir := os.TempDir()
	defer os.Remove(dir)

	t.Run("With success case: return no error", func(tc *testing.T) {
		tf := NewTfvars()

		raw := map[string]interface{}{
			"test": "this is a str",
			"arr": []interface{}{
				map[string]interface{}{
					"node":  1,
					"ready": true,
				},
				map[string]interface{}{
					"node":  2,
					"ready": false,
				},
			},
			"bool_val": false,
		}
		outputPath := dir + "/out.tfvars"
		defer os.Remove(outputPath)
		err := tf.Convert(raw, outputPath)
		assert.NoError(tc, err)

		content, _ := osext.Cat(outputPath)
		assert.Contains(tc, string(content), `test = "this is a str"`)
		assert.Contains(tc, string(content), `bool_val = false`)
		assert.Contains(tc, string(content), `arr = [{
  node  = 1
  ready = true
  }, {
  node  = 2
  ready = false
}]`)
	})

}
