package env

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEnv_ToENV(t *testing.T) {
	t.Run("With success case: return byte array", func(tc *testing.T) {
		e := NewEnv()
		input := map[string]interface{}{
			"TEST":  "TEST",
			"DATA":  []int{1, 2, 3},
			"STR":   []string{"a", "b", "c"},
			"EMPTY": map[string]interface{}{"test": true},
			"FLOAT": 0.23,
			"STR_WITH_SPACE": "string with space",
		}

		expected := map[string]bool{
			"export TEST=TEST": true,
			"export DATA=": true,
			"export STR=": true,
			"export EMPTY=": true,
			"export FLOAT=0.23": true,
			"export STR_WITH_SPACE=\"string with space\"": true,
		}

		actual := strings.Split(string(e.ToENV(input)), "\n")

		assert.Equal(tc, len(expected), len(actual))
		for _, a := range actual {
			assert.Contains(tc, expected, a)
		}
	})
}

func TestEnv_ToJSON(t *testing.T) {
	t.Run("With success case: return a map", func(tc *testing.T) {
		e := NewEnv()
		input := []string{
			"export TEST=true",
			"export DATA=1,2,3",
			"#export PASS=no",
			"CI=true=test",
			"T=",
			"export T1='test data'",
			"export T2=\"test data\"",
		}

		expected := map[string]interface{}{
			"TEST": "true",
			"DATA": "1,2,3",
			"CI":   "true=test",
			"T":    "",
			"T1":   "test data",
			"T2":   "test data",
		}

		content, err := e.ToJSON(input)
		assert.NoError(tc, err)
		assert.Equal(tc, len(expected), len(content))

		for k, v := range expected {
			assert.Contains(tc, content, k)
			assert.Equal(tc, v, content[k])
		}
	})

	t.Run("With error case: return an error key invalid", func(tc *testing.T) {
		e := NewEnv()
		input := []string{
			"export test = true",
		}
		_, err := e.ToJSON(input)
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Env: Key test  is invalid format")
	})
}
