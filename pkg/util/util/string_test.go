package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrim(t *testing.T) {
	t.Run("With success case: return trimmed string", func(tc *testing.T) {
		inputs := []string{
			"  test   ",
			"test \n",
			"test\n    ",
			"test   \n",
		}
		for _, input := range inputs {

			assert.Equal(tc, "test", Trim(input))
		}

	})
}
