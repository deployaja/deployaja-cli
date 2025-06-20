package cmd

import (
	"fmt"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd())
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate deployaja.yaml configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.LoadDeploymentConfig()
			if err != nil {
				return err
			}

			// Basic validation
			if cfg.Name == "" {
				return fmt.Errorf("name is required")
			}

			if cfg.Container.Image == "" {
				return fmt.Errorf("container.image is required")
			}

			if cfg.Container.Port == 0 {
				return fmt.Errorf("container.port is required")
			}

			fmt.Printf("%s Configuration is valid\n", ui.SuccessPrint("✓"))
			return nil
		},
	}
}
