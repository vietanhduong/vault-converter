package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vietanhduong/vault-converter/pkg/converter"
	"github.com/vietanhduong/vault-converter/pkg/util/env"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"github.com/vietanhduong/vault-converter/pkg/vault"
)

var convertCmd = &cobra.Command{
	Use:   "pull SECRET_PATH",
	Short: "Pull secrets from Vault and convert to file",
	Long: `Pull secrets from Vault with specified secret path and convert to file.
SECRET_PATH should be a absolute path at Vault and the values should be in JSON format.
Supports the following formats: .tfvars (Terraform variables)
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("argument SECRET_PATH must be specified")
		}
		address := cmd.Flag("address").Value.String()
		output := cmd.Flag("output").Value.String()
		format := cmd.Flag("format").Value.String()
		secretPath := args[0]

		token, err := os.Cat(vault.DefaultTokenPath)
		if err != nil {
			return err
		}

		if util.IsNullOrEmpty(token) {
			return errors.New("Vault: Unauthorized")
		}

		v := vault.New(address, token)

		values, err := v.Read(secretPath)
		if err != nil {
			return err
		}

		c, err := converter.NewConverter(format)
		if err != nil {
			return err
		}

		return c.Convert(values, output)
	},
}

func init() {
	flags := convertCmd.Flags()
	flags.StringP("address", "a", env.GetEnvAsStringOrFallback("VAULT_ADDR", "https://127.0.0.1:8200"), "Address of the Auth server. This can also be specified via the VAULT_ADDR environment variable.")
	flags.StringP("output", "o", "variables.auto.tfvars", "Output path. E.g: ~/data/variables.auto.tfvars")
	flags.StringP("format", "f", "tfvars", "Output format")
	rootCmd.AddCommand(convertCmd)
}
