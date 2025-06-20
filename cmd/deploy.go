package cmd

import (
	"fmt"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd())
}

func deployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy application to cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			cfg, err := config.LoadDeploymentConfig()
			if err != nil {
				return err
			}

			fmt.Printf("%s Deploying %s...\n", ui.InfoPrint("🚀"), cfg.Name)

			response, err := apiClient.Deploy(cfg, false)
			if err != nil {
				return err
			}

			fmt.Printf("%s %s\n", ui.SuccessPrint("✓"), response.Message)
			fmt.Printf("Deployment ID: %s\n", response.DeploymentID)

			if response.EstimatedTime != "" {
				fmt.Printf("Estimated time: %s\n", response.EstimatedTime)
			}

			if response.URL != "" {
				fmt.Printf("URL: %s\n", response.URL)
			}

			return nil
		},
	}

	return cmd
}
