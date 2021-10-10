package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "vault-converter",
	Short: "Convert to file from Vault",
	Long: `
Convert to file from Vault. Support multiple file format 
like '.tfvars', '.env'
`}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

