package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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
	var setFlags []string
	var dryRun bool
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

			if err != nil {
				return err
			}

			// Override name if provided via flag
			if nameFlag != "" {
				cfg.Name = nameFlag
			}

			// Apply --set overrides
			if err := applySetOverrides(cfg, setFlags); err != nil {
				return err
			}

			fmt.Printf("%s Deploying %s...\n", ui.InfoPrint("🚀"), cfg.Name)

			response, err := apiClient.Deploy(cfg, dryRun)
			if err != nil {
				return err
			}

			fmt.Printf("%s %s\n", ui.SuccessPrint("✓"), response.Message)			

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
	cmd.Flags().StringSliceVar(&setFlags, "set", []string{}, "Set configuration values using dot notation (e.g., --set container.image=nginx:latest --set resources.replicas=3)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Dry run the deployment")
	return cmd
}

// applySetOverrides applies configuration overrides from --set flags
func applySetOverrides(cfg *config.DeploymentConfig, setFlags []string) error {
	for _, override := range setFlags {
		parts := strings.SplitN(override, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid --set format: %s (expected key=value)", override)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if err := setConfigValue(cfg, key, value); err != nil {
			return fmt.Errorf("failed to set %s: %v", key, err)
		}
	}
	return nil
}

// setConfigValue sets a configuration value using dot notation
func setConfigValue(cfg *config.DeploymentConfig, path, value string) error {
	parts := strings.Split(path, ".")

	switch parts[0] {
	case "name":
		cfg.Name = value
	case "description":
		cfg.Description = value
	case "domain":
		cfg.Domain = value
	case "container":
		if len(parts) < 2 {
			return fmt.Errorf("container requires a subfield (e.g., container.image)")
		}
		switch parts[1] {
		case "image":
			cfg.Container.Image = value
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("container.port must be a number: %v", err)
			}
			cfg.Container.Port = port
		default:
			return fmt.Errorf("unknown container field: %s", parts[1])
		}
	case "resources":
		if len(parts) < 2 {
			return fmt.Errorf("resources requires a subfield (e.g., resources.cpu)")
		}
		switch parts[1] {
		case "cpu":
			cfg.Resources.CPU = value
		case "memory":
			cfg.Resources.Memory = value
		case "replicas":
			replicas, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("resources.replicas must be a number: %v", err)
			}
			cfg.Resources.Replicas = replicas
		default:
			return fmt.Errorf("unknown resources field: %s", parts[1])
		}
	case "healthCheck":
		if len(parts) < 2 {
			return fmt.Errorf("healthCheck requires a subfield (e.g., healthCheck.path)")
		}
		switch parts[1] {
		case "path":
			cfg.HealthCheck.Path = value
		case "port":
			port, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("healthCheck.port must be a number: %v", err)
			}
			cfg.HealthCheck.Port = port
		case "initialDelaySeconds":
			delay, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("healthCheck.initialDelaySeconds must be a number: %v", err)
			}
			cfg.HealthCheck.InitialDelaySeconds = delay
		case "periodSeconds":
			period, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("healthCheck.periodSeconds must be a number: %v", err)
			}
			cfg.HealthCheck.PeriodSeconds = period
		default:
			return fmt.Errorf("unknown healthCheck field: %s", parts[1])
		}
	default:
		return fmt.Errorf("unknown configuration field: %s", parts[0])
	}

	return nil
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
