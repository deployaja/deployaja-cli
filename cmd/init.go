package cmd

import (
	"fmt"
	"os"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(initCmd())
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create deployaja.yaml configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(config.DeployFile); err == nil {
				return fmt.Errorf("deployaja.yaml already exists")
			}

			cfg := config.DeploymentConfig{
				Name:        "my-application",
				Version:     "1.0.0",
				Description: "My awesome application",
			}

			cfg.Container.Image = "nginx:latest"
			cfg.Container.Port = 80

			cfg.Resources.CPU = "100m"
			cfg.Resources.Memory = "128Mi"
			cfg.Resources.Replicas = 1

			cfg.HealthCheck.Path = "/health"
			cfg.HealthCheck.Port = 80
			cfg.HealthCheck.InitialDelaySeconds = 30
			cfg.HealthCheck.PeriodSeconds = 10

			cfg.Env = []config.EnvVar{
				{Name: "NODE_ENV", Value: "production"},
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}

			err = os.WriteFile(config.DeployFile, data, 0644)
			if err != nil {
				return err
			}

			fmt.Printf("%s Created %s\n", ui.SuccessPrint("✓"), config.DeployFile)
			fmt.Printf("%s Edit the file to configure your application\n", ui.InfoPrint("→"))

			return nil
		},
	}
}
