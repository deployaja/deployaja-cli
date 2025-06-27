package cmd

import (
	"fmt"

	"deployaja-cli/internal/version"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		versionStr := version.GetVersion()
		fmt.Printf("DeployAja CLI version: %s\n", versionStr)
	},
}

// ReadVersionFromFile returns the current CLI version
// Made public so it can be used by other packages like the API client
// Deprecated: Use internal/version.GetVersion() instead
func ReadVersionFromFile() string {
	return version.GetVersion()
}
