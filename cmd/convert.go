package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert SECRET_PATH",
	Short: "Convert secrets to file",
	Long: `
Convert secrets to files.
Supports the following formats: .tfvars (Terraform variables), .env 
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("argument SECRET_PATH must be specified")
		}
		return nil
	},
}

func init() {
	flags := convertCmd.Flags()
	flags.StringP("format", "f", "tfvars", "Output format")
	rootCmd.AddCommand(convertCmd)
}
