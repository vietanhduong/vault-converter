package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/vietanhduong/vault-converter/pkg/util/output"
)

var (
	// GitCommit is updated with the Git tag by the Goreleaser build
	GitCommit = "unknown"
	// BuildDate is updated with the current ISO timestamp by the Goreleaser build
	BuildDate = "unknown"
	// Version is updated with the latest tag by the Goreleaser build
	Version = "unreleased"
)

func showVersion() {
	output.Printf("Version:\t %s\n", Version)
	output.Printf("Git commit:\t %s\n", GitCommit)
	output.Printf("Date:\t\t %s\n", BuildDate)
	output.Printf("License:\t Apache 2.0\n")
}

var rootCmd = &cobra.Command{
	Use:          "vault-converter",
	Short:        "Convert to file from Vault",
	Long:         `Convert to file from Vault. Support multiple file format like '.tfvars', '.env'`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		version := cmd.Flag("version").Value.String()
		if version == "true" {
			showVersion()
		} else {
			_ = cmd.Help()
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version information and exit. This flag is only available at the global level.")
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
