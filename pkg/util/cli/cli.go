package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"strings"
)

func ArgumentsRequired(cmd *cobra.Command, requiredArgs []string) error {
	flags := make(map[string]string)
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		flags[flag.Name] = flag.Value.String()
	})

	var missing []string

	for _, arg := range requiredArgs {
		v, found := flags[arg]
		if !found || util.IsNullOrEmpty(v) {
			missing = append(missing, fmt.Sprintf("\"%s\"", arg))
		}
	}
	if len(missing) > 0 {
		miss := strings.Join(missing, ", ")
		return errors.New(fmt.Sprintf("Error: required flag(s) %s not set", miss))
	}
	return nil
}
