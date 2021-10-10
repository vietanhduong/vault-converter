package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vietanhduong/vault-converter/pkg/util/cli"
	"github.com/vietanhduong/vault-converter/pkg/util/env"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"github.com/vietanhduong/vault-converter/pkg/vault"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticates users to Vault",
	Long: `
Authenticates users to Vault using the provided arguments. 
Method using: 'userpass'. The path of 'userpass' should be 'userpass/'  
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.ArgumentsRequired(cmd, []string{"address", "username", "password"}); err != nil {
			return err
		}

		address := cmd.Flag("address").Value.String()
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()

		auth := vault.New(address, username, password)
		if err := auth.Login(); err != nil {
			return err
		}
		output.Printf(fmt.Sprintf("Success! You are now authenticated.\nYou do NOT need to run '%s %s' again.", rootCmd.Use, cmd.Use))
		return nil
	},
}

func init() {
	flags := authCmd.Flags()
	flags.StringP("address", "a", env.GetEnvAsStringOrFallback("VAULT_ADDR", "https://127.0.0.1:8200"), "Address of the Vault server. This can also be specified via the VAULT_ADDR environment variable.")
	flags.StringP("username", "u", env.GetEnvAsStringOrFallback("VAULT_USER", ""), "The username to authenticate with Vault server. This can also be specified via the VAULT_USER environment variables.")
	flags.StringP("password", "p", env.GetEnvAsStringOrFallback("VAULT_PASSWORD", ""), "The user's password. This can also be specified via the VAULT_PASSWORD environment variables.")
	rootCmd.AddCommand(authCmd)
}
