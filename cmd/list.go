package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/vietanhduong/vault-converter/pkg/util/env"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"github.com/vietanhduong/vault-converter/pkg/vault"
)

var lsCmd = &cobra.Command{
	Use:   "ls [SECRET_PATH]",
	Short: "List secrets in path",
	Long: `List all secret in input [SECRET_PATH]. If the input [SECRET_PATH] is empty.
It will return all secret from 'root'. It will try to get all secret that current user
can 'read'.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var secretPath string
		if len(args) > 0 {
			secretPath = args[0]
		}

		address := cmd.Flag("address").Value.String()
		recursive, _ := cmd.Flags().GetBool("recursive")

		token, err := os.Cat(vault.DefaultTokenPath)
		if err != nil {
			return err
		}

		if util.IsNullOrEmpty(token) {
			return errors.New("vault: Unauthorized")
		}
		v := vault.New(address, util.Trim(string(token)))
		secrets, err := v.List(secretPath, recursive)
		if err != nil {
			return err
		}

		for _, secret := range secrets {
			output.Printf("%s\n", secret)
		}
		return nil
	},
}

func init() {
	flags := lsCmd.Flags()
	flags.StringP("address", "a", env.GetEnvAsStringOrFallback("VAULT_ADDR", "https://127.0.0.1:8200"), "addr of the Auth server. This can also be specified via the VAULT_ADDR environment variable.")
	flags.BoolP("recursive", "r", false, "List secret recursive or not.")
	rootCmd.AddCommand(lsCmd)
}
