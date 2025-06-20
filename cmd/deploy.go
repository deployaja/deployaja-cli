package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd())
}

func deployCmd() *cobra.Command {
	var buildFlag bool
	var fileFlag string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy application to cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			// Validate Dockerfile exists if --build flag is used
			if buildFlag {
				if err := validateDockerfileExists(); err != nil {
					return err
				}
			}

			// Load config from specified file or default
			var cfg *config.DeploymentConfig
			var err error
			if fileFlag != "" {
				cfg, err = config.LoadDeploymentConfigFromFile(fileFlag)
			} else {
				cfg, err = config.LoadDeploymentConfig()
			}
			if err != nil {
				return err
			}

			fmt.Printf("%s Deploying %s...\n", ui.InfoPrint("🚀"), cfg.Name)

			response, err := apiClient.Deploy(cfg, buildFlag)
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

	cmd.Flags().BoolVar(&buildFlag, "build", false, "Build the application before deployment (requires Dockerfile)")
	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to deployment configuration file (default: deployaja.yaml)")

	return cmd
}

// validateDockerfileExists checks if a Dockerfile exists in the current directory
func validateDockerfileExists() error {
	dockerfilePath := filepath.Join(".", "Dockerfile")

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("%s Dockerfile not found in current directory. Use --build flag only when Dockerfile is present", ui.ErrorPrint("✗"))
	}

	return nil
}
