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

			// Use the global API client with proper authentication
			// Call API to validate configuration
			validateResp, err := apiClient.Validate(cfg)
			if err != nil {
				// If API validation fails, show the error
				return fmt.Errorf("validation failed: %v", err)
			}

			if !validateResp.Valid {
				return fmt.Errorf("configuration is invalid: %s", validateResp.Message)
			}

			// Configuration is valid
			fmt.Printf("%s Configuration is valid\n", ui.SuccessPrint("✓"))

			// Show any warnings if present
			if len(validateResp.Warnings) > 0 {
				fmt.Println("\nWarnings:")
				for _, warning := range validateResp.Warnings {
					fmt.Printf("%s %s\n", ui.WarningPrint("⚠"), warning)
				}
			}

			return nil
		},
	}
}
