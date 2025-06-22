package cmd

import (
	"fmt"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd())
}

func rollbackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rollback NAME",
		Short: "Rollback deployment to previous version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			name := args[0]

			fmt.Printf("%s Rolling back %s...\n", ui.InfoPrint("⏪"), name)

			err := apiClient.Rollback(name, "previous")
			if err != nil {
				return err
			}

			fmt.Printf("%s Rollback initiated for %s\n", ui.SuccessPrint("✓"), name)

			return nil
		},
	}
}
