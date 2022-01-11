package pull

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	osext "github.com/vietanhduong/vault-converter/pkg/util/os"
)

func TestEnv_Convert(t *testing.T) {
	dir := os.TempDir()
	defer os.Remove(dir)

	t.Run("With success case: return no error", func(tc *testing.T) {
		e := NewEnv()
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
		outputPath := dir + "/.env"
		defer os.Remove(outputPath)
		err := e.Convert(raw, outputPath)
		assert.NoError(tc, err)

		content, _ := osext.Cat(outputPath)
		assert.Contains(tc, string(content), `export test='this is a str'`)
		assert.Contains(tc, string(content), `export bool_val=false`)
		assert.Contains(tc, string(content), `export arr=`)
	})

}
