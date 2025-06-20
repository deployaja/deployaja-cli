package cmd

import (
	"fmt"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(planCmd())
}

func planCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plan",
		Short: "Show deployment plan and cost forecasting",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			cfg, err := config.LoadDeploymentConfig()
			if err != nil {
				return err
			}

			fmt.Printf("%s Calculating deployment costs...\\n", ui.InfoPrint("→"))

			response, err := apiClient.GetCostEstimate(cfg)
			if err != nil {
				return err
			}

			// Display plan
			fmt.Printf("\n%s Deployment Plan\n", ui.InfoPrint("📋"))
			fmt.Printf("Application: %s\n", cfg.Name)
			fmt.Printf("Image: %s\n", cfg.Container.Image)
			fmt.Printf("Replicas: %d\n", cfg.Resources.Replicas)

			if len(cfg.Dependencies) > 0 {
				fmt.Printf("\nDependencies:\n")
				for _, dep := range cfg.Dependencies {
					fmt.Printf("  - %s (%s %s)\n", dep.Name, dep.Type, dep.Version)
				}
			}

			// Display costs
			fmt.Printf("\n%s Cost Estimate\n", ui.InfoPrint("💰"))
			fmt.Printf("Monthly: $%.2f %s\n", response.EstimatedCost.Monthly, response.EstimatedCost.Currency)
			fmt.Printf("Daily: $%.2f %s\n", response.EstimatedCost.Daily, response.EstimatedCost.Currency)

			fmt.Printf("\nBreakdown:\n")
			fmt.Printf("  Compute: $%.2f\n", response.Breakdown.Compute)
			fmt.Printf("  Storage: $%.2f\n", response.Breakdown.Storage)
			fmt.Printf("  Network: $%.2f\n", response.Breakdown.Network)

			if len(response.Breakdown.Dependencies) > 0 {
				for name, cost := range response.Breakdown.Dependencies {
					fmt.Printf("  %s: $%.2f\n", name, cost)
				}
			}

			return nil
		},
	}
}
