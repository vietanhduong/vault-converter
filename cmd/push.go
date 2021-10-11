package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/vietanhduong/vault-converter/pkg/push"
	"github.com/vietanhduong/vault-converter/pkg/util/env"
	"github.com/vietanhduong/vault-converter/pkg/util/os"
	"github.com/vietanhduong/vault-converter/pkg/util/output"
	"github.com/vietanhduong/vault-converter/pkg/util/util"
	"github.com/vietanhduong/vault-converter/pkg/vault"
)

var pushCmd = &cobra.Command{
	Use:   "push SOURCE_FILE SECRET_PATH",
	Short: "Parse source file and push to Vault",
	Long: `Parse source file and push secrets to Vault.
Based on the extension of SOURCE_FILE to determine the file format. 
SECRET_PATH should be a absolute path at Vault and the values should 
be in JSON format.
Supports the following formats: "tfvars"
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("arguments SOURCE_FILE and SECRET_PATH must be specified")
		}
		address := cmd.Flag("address").Value.String()
		sourceFile := args[0]
		secretPath := args[1]

		// Verify that source file is exist and NOT empty
		content, err := os.Cat(sourceFile)
		if err != nil {
			return err
		}

		if util.IsNullOrEmpty(content) {
			return errors.New("Source file is empty")
		}

		// GetExtension return a file extension.
		// Note that if your file is an `env` file.
		// It should be end with `.env`
		pusher, err := push.NewConverter(os.GetExtension(sourceFile))
		if err != nil {
			return err
		}

		values, err := pusher.Convert(content)
		if err != nil {
			return err
		}

		token, err := os.Cat(vault.DefaultTokenPath)
		if err != nil {
			return err
		}

		if util.IsNullOrEmpty(token) {
			return errors.New("Vault: Unauthorized")
		}

		v := vault.New(address, string(token))
		if err = v.Write(secretPath, values); err != nil {
			return err
		}

		output.Printf("Successfully pushed to %s", secretPath)
		return nil
	},
}

func init() {
	flags := pushCmd.Flags()
	flags.StringP("address", "a", env.GetEnvAsStringOrFallback("VAULT_ADDR", "https://127.0.0.1:8200"), "Address of the Auth server. This can also be specified via the VAULT_ADDR environment variable.")
	rootCmd.AddCommand(pushCmd)
}
