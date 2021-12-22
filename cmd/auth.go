package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vietanhduong/vault-converter/pkg/util/cli"
	"github.com/vietanhduong/vault-converter/pkg/util/env"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"github.com/vietanhduong/vault-converter/pkg/vault"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticates users to vault",
	Long: `Authenticates users to vault using the provided arguments. 
Method using: 'userpass'. The path of 'userpass' should be 'userpass/'  
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cli.FlagsRequired(cmd, []string{"username", "password"}); err != nil {
			return err
		}

		address := cmd.Flag("address").Value.String()
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()

		auth := vault.NewAuth(address, username, password)
		clientToken, err := auth.Login()
		if err != nil {
			return err
		}

		// To reduce the number of variables passed during program execution.
		// After successful login, user's token will be saved at a fixed path (like the way vault is using).
		if err = os.Write([]byte(clientToken), vault.DefaultTokenPath); err != nil {
			return err
		}

		output.Printf(fmt.Sprintf("Success! You are now authenticated.\nYou do NOT need to run '%s %s' again.", rootCmd.Use, cmd.Use))
		return nil
	},
}

func init() {
	flags := authCmd.Flags()
	flags.StringP("address", "a", env.GetEnvAsStringOrFallback("VAULT_ADDR", "https://127.0.0.1:8200"), "addr of the Auth server. This can also be specified via the VAULT_ADDR environment variable.")
	flags.StringP("username", "u", env.GetEnvAsStringOrFallback("VAULT_USER", ""), "The username to authenticate with Auth server. This can also be specified via the VAULT_USER environment variables.")
	flags.StringP("password", "p", env.GetEnvAsStringOrFallback("VAULT_PASSWORD", ""), "The user's password. This can also be specified via the VAULT_PASSWORD environment variables.")
	rootCmd.AddCommand(authCmd)
}
