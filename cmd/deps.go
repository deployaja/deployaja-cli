package cmd

import (
	"fmt"
	"strings"

	"deployaja-cli/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(depsCmd())
}

func depsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps",
		Short: "List available dependencies and versions",
		RunE: func(cmd *cobra.Command, args []string) error {
			depType, _ := cmd.Flags().GetString("type")

			response, err := apiClient.GetDependencies(depType)
			if err != nil {
				return err
			}

			fmt.Printf("%s Available Dependencies\n\n", ui.InfoPrint("🔧"))

			for _, dep := range response.Dependencies {
				fmt.Printf("%s\n", color.New(color.Bold).Sprint(dep.Name))
				fmt.Printf("  Type: %s\n", dep.Type)
				fmt.Printf("  Versions: %s (default: %s)\n",
					strings.Join(dep.Versions, ", "), dep.DefaultVersion)
				fmt.Printf("  Base Cost: $%.2f/month\n", dep.Pricing.Base)

				if dep.Pricing.Storage > 0 {
					fmt.Printf("  Storage: $%.2f/GB/month\n", dep.Pricing.Storage)
				}

				fmt.Printf("  Specs: %s, %s\n", dep.Specs.CPU, dep.Specs.Memory)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().String("type", "", "Filter by dependency type")
	return cmd
}
