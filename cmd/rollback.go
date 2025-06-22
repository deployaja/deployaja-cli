package cmd

import (
	"fmt"
	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd())
}

func rollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback",
		Short: "Rollback deployment to previous version",		
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			cfg, err := config.LoadDeploymentConfig()
			if err != nil {
				return err
			}
			name := cfg.Name

			fmt.Printf("%s Rolling back %s...\n", ui.InfoPrint("⏪"), name)

			err = apiClient.Rollback(name)
			if err != nil {
				return err
			}

			fmt.Printf("%s Rollback initiated for %s\n", ui.SuccessPrint("✓"), name)

			return nil
		},
	}
}
