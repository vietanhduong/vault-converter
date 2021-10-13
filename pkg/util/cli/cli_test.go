package cli

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlagsRequired(t *testing.T) {
	c := &cobra.Command{}
	flags := c.Flags()
	flags.StringP("username", "u", "", "")
	flags.StringP("password", "p", "", "")
	flags.StringP("address", "a", "", "")

	_ = flags.Set("username", "")
	_ = flags.Set("password", "")
	_ = flags.Set("address", "true")

	t.Run("With success case: return no error", func(tc *testing.T) {
		err := FlagsRequired(c, []string{"address"})
		assert.NoError(tc, err)
	})

	t.Run("With error case: failure because missing value", func(tc *testing.T) {
		err := FlagsRequired(c, []string{"username", "password", "address"})
		assert.Error(tc, err)
		assert.Contains(tc, err.Error(), "Error: required flag(s) \"username\", \"password\" not set")
	})

}
