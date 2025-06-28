package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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

			// Generate random name using Wayang myth characters
			wayangNames := []string{"bima", "arjuna", "werkudara", "nakula", "sadewa", "gatotkaca", "krishna", "yudhistira", "srikandi", "abimanyu"}
			rand.Seed(time.Now().UnixNano())
			randomName := wayangNames[rand.Intn(len(wayangNames))]
			randomNumber := rand.Intn(90) + 10 // 2-digit number
			appName := fmt.Sprintf("%s-%d-app", randomName, randomNumber)

			cfg := config.DeploymentConfig{
				Name:        appName,				
				Description: "Simple web application with nginx and postgres",
			}

			// Container configuration
			cfg.Container.Image = "nginx:latest"
			cfg.Container.Port = 80

			// Resource allocation
			cfg.Resources.CPU = "500m"
			cfg.Resources.Memory = "1Gi"
			cfg.Resources.Replicas = 2

			// Dependencies
			cfg.Dependencies = []config.Dependency{
				{
					Name:    "postgresql",
					Type:    "postgresql",
					Version: "15",
					Storage: "1Gi",
				},
			}

			// Environment variables
			cfg.Env = []config.EnvVar{
				{Name: "NODE_ENV", Value: "production"},
				{Name: "LOG_LEVEL", Value: "info"},
				{Name: "CACHE_TTL", Value: "3600"},
				{Name: "MAX_CONNECTIONS", Value: "1000"},
			}

			// Health check configuration
			cfg.HealthCheck.Path = "/api/health"
			cfg.HealthCheck.Port = 8080
			cfg.HealthCheck.InitialDelaySeconds = 60
			cfg.HealthCheck.PeriodSeconds = 30

			// Domain (optional)
			// Set a random domain using the same Wayang character name
			cfg.Domain = fmt.Sprintf("%s%d.deployaja.id", randomName, randomNumber)

			// Volumes (optional)
			cfg.Volumes = []config.Volume{
				{
					Name:      "app-storage",
					Size:      "1Gi",
					MountPath: "/usr/share/nginx/html",
				},
			}

			data, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}

			err = os.WriteFile(config.DeployFile, data, 0644)
			if err != nil {
				return err
			}

			fmt.Printf("%s Created %s with comprehensive configuration\n", ui.SuccessPrint("✓"), config.DeployFile)
			fmt.Printf("%s Edit the file to customize your application configuration\n", ui.InfoPrint("→"))
			fmt.Printf("%s Remove unused sections (dependencies, volumes, domain) if not needed\n", ui.InfoPrint("💡"))

			return nil
		},
	}
}
