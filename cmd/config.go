package cmd

import (
	"fmt"
	"path/filepath"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd())
}

func configCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			configPath := filepath.Join(home, config.ConfigDir)

			fmt.Printf("Configuration file: %s\n", configPath)
			fmt.Printf("Token file: %s\n", filepath.Join(home, config.ConfigDir, config.TokenFile))

			if apiClient.Token != "" {
				fmt.Printf("Authentication: %s\n", ui.SuccessPrint("✓ Authenticated"))
			} else {
				fmt.Printf("Authentication: %s\n", ui.ErrorPrint("✗ Not authenticated"))
			}

			return nil
		},
	}
}
