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
	var fileFlag string
	var nameFlag string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy application to cloud",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			// Load config from specified file or default
			var cfg *config.DeploymentConfig
			var err error
			if fileFlag != "" {
				cfg, err = config.LoadDeploymentConfigFromFile(fileFlag)
			} else {
				if err := validateDefaultConfigExists(); err != nil {
					return err
				}
				cfg, err = config.LoadDeploymentConfig()
			}

			// Override name if provided via flag
			if nameFlag != "" {
				cfg.Name = nameFlag
			}

			if err := validateDockerfileExists(); err != nil {
				return err
			}

			fmt.Printf("%s Deploying %s...\n", ui.InfoPrint("🚀"), cfg.Name)

			response, err := apiClient.Deploy(cfg, true)
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

	cmd.Flags().StringVarP(&fileFlag, "file", "f", "", "Path to deployment configuration file (required if deployaja.yaml doesn't exist)")
	cmd.Flags().StringVarP(&nameFlag, "name", "n", "", "Override the API name for deployment")

	return cmd
}

// validateDockerfileExists checks if a Dockerfile exists in the current directory
func validateDockerfileExists() error {
	dockerfilePath := filepath.Join(".", "Dockerfile")

	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		return fmt.Errorf("%s Dockerfile not found in current directory. Dockerfile is required for deployment", ui.ErrorPrint("✗"))
	}

	return nil
}

// validateDefaultConfigExists checks if the default config file exists
func validateDefaultConfigExists() error {
	defaultConfigPath := filepath.Join(".", "deployaja.yaml")

	if _, err := os.Stat(defaultConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("%s Default deployment configuration file not found in current directory. Use 'deployaja init' to create one or specify a config file with -f flag", ui.ErrorPrint("✗"))
	}

	return nil
}
