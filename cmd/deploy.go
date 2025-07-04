package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"deployaja-cli/internal/api"
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

			fmt.Printf("%s Starting deployment process...\n", ui.InfoPrint("🔧"))

			// Load config from specified file or default
			var cfg *config.DeploymentConfig
			var err error
			if fileFlag != "" {
				fmt.Printf("%s Loading configuration from file: %s\n", ui.InfoPrint("📄"), fileFlag)
				cfg, err = config.LoadDeploymentConfigFromFile(fileFlag)
			} else {
				fmt.Printf("%s Checking for default configuration file...\n", ui.InfoPrint("📄"))
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
				fmt.Printf("%s Overriding deployment name: %s\n", ui.InfoPrint("✏️"), nameFlag)
				cfg.Name = nameFlag
			}

			// Apply --set overrides
			if len(setFlags) > 0 {
				fmt.Printf("%s Applying configuration overrides...\n", ui.InfoPrint("🛠️"))
			}
			if err := applySetOverrides(cfg, setFlags); err != nil {
				return err
			}

			fmt.Printf("%s Deploying %s...\n", ui.InfoPrint("🚀"), cfg.Name)

			fmt.Printf("%s Sending deployment request to API...\n", ui.InfoPrint("📡"))
			response, err := apiClient.Deploy(cfg, dryRun)
			if err != nil {
				fmt.Printf("%s Deployment request failed: %v\n", ui.ErrorPrint("❌"), err)
				return nil // Prevent Cobra from printing usage for runtime errors
			}

			fmt.Printf("%s %s\n", ui.SuccessPrint("✓"), response.Message)

			if response.URL != "" {
				fmt.Printf("URL: %s\n", response.URL)
			}

			// Don't poll status if it's a dry run
			if dryRun {
				fmt.Printf("%s Dry run complete. No resources were created.\n", ui.InfoPrint("💡"))
				return nil
			}

			// Poll for deployment status until completion
			fmt.Printf("%s Monitoring deployment progress...\n", ui.InfoPrint("🔍"))

			var lastProgress int
			var lastStep string
			progressCallback := func(progress *api.DeploymentProgress) {
				if progress.Progress.Percentage != lastProgress || progress.Progress.CurrentStep != lastStep {
					// Create a simple progress bar
					progressBar := createProgressBar(progress.Progress.Percentage)

					fmt.Printf("\r%s Progress: [%s] %d%% - %s",
						ui.InfoPrint("📊"),
						progressBar,
						progress.Progress.Percentage,
						progress.Progress.CurrentStep)

					// Add queue details if available
					if progress.Queue.Status == "active" && progress.Queue.Message != "" {
						fmt.Printf(" (%s)", progress.Queue.Message)
					}

					lastProgress = progress.Progress.Percentage
					lastStep = progress.Progress.CurrentStep

					// Add newline if complete
					if progress.Progress.IsComplete {
						fmt.Printf("\n")
					}
				}
			}

			finalDeployment, err := apiClient.PollDeploymentProgress(cfg.Name, progressCallback)
			if err != nil {
				fmt.Printf("\n%s Warning: Failed to monitor deployment progress: %v\n", ui.WarningPrint("⚠️"), err)
				fmt.Printf("%s You can check the status manually using: deployaja status\n", ui.InfoPrint("💡"))
				return nil // Prevent Cobra from printing usage for runtime errors
			}

			// Show final status
			if finalDeployment != nil && (finalDeployment.Status == "running" || finalDeployment.Status == "success") {
				fmt.Printf("\n%s Deployment completed successfully!\n", ui.SuccessPrint("🎉"))
				if finalDeployment.URL != "" {
					fmt.Printf("%s Access your application at: %s\n", ui.InfoPrint("🌐"), finalDeployment.URL)
				}
			} else if finalDeployment != nil {
				fmt.Printf("\n%s Deployment failed with status: %s\n", ui.ErrorPrint("❌"), finalDeployment.Status)
				fmt.Printf("%s Use 'deployaja describe %s' for more details\n", ui.InfoPrint("💡"), cfg.Name)

				// Try to fetch and print more info from the server
				fmt.Printf("%s Fetching deployment details...\n", ui.InfoPrint("🔍"))
				if describe, err := apiClient.Describe(cfg.Name); err == nil {
					fmt.Printf("\n%s Deployment Error Details:\n", ui.ErrorPrint("---"))
					if len(describe.Pod) > 0 {
						fmt.Printf("Pod Info:\n")
						for k, v := range describe.Pod {
							fmt.Printf("  %s: %v\n", k, v)
						}
					}
					if len(describe.Events) > 0 {
						fmt.Printf("\nRecent Events:\n")
						for _, event := range describe.Events {
							if reason, ok := event["reason"]; ok {
								fmt.Printf("- Reason: %v", reason)
							} else {
								fmt.Printf("- Event:")
							}
							if msg, ok := event["message"]; ok {
								fmt.Printf(", Message: %v", msg)
							}
							if t, ok := event["type"]; ok {
								fmt.Printf(", Type: %v", t)
							}
							fmt.Printf("\n")
						}
					}
					if len(describe.Pod) == 0 && len(describe.Events) == 0 {
						fmt.Printf("No detailed error information available from server.\n")
					}
				} else {
					fmt.Printf("%s Failed to fetch deployment details: %v\n", ui.WarningPrint("⚠️"), err)
				}
				return nil // Prevent Cobra from printing usage for runtime errors
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

// createProgressBar creates a visual progress bar
func createProgressBar(percentage int) string {
	const width = 30
	filled := (percentage * width) / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}
