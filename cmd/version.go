package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		version := readVersionFromFile()
		fmt.Printf("DeployAja CLI version: %s\n", version)
	},
}

func readVersionFromFile() string {
	return "beta-0.0.1"
}
