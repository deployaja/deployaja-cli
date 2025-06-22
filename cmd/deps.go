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
		Use:   "deps [instance]",
		Short: "List available dependencies and versions",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var instance string
			if len(args) > 0 {
				instance = args[0]
			}

			depType, _ := cmd.Flags().GetString("type")
			
			if instance != "" {
				response, err := apiClient.GetDependencyInstance()
				if err != nil {
					return err
				}
				fmt.Printf("%s Dependency Instance Details\n\n", ui.InfoPrint("🔍"))
				for _, inst := range *&response.Instances {
					fmt.Printf("ID:        %v\n", inst.ID)
					fmt.Printf("User ID:   %v\n", inst.UserID)
					fmt.Printf("Type:      %v\n", inst.Type)
					fmt.Printf("Created:   %v\n", inst.CreatedAt)
					fmt.Printf("Updated:   %v\n", inst.UpdatedAt)
					fmt.Printf("Config:\n")
					// Pretty print config (assume it's a map or struct)
					switch cfg := inst.Config.(type) {
					case map[string]interface{}:
						for k, v := range cfg {
							fmt.Printf("  %s: %v\n", k, v)
						}
					default:
						fmt.Printf("  %v\n", inst.Config)
					}
					fmt.Println()
				}
				return nil
			} else {
				fmt.Printf("%s Available Dependencies\n\n", ui.InfoPrint("🔧"))
			}
			response, err := apiClient.GetDependencies(depType)
			if err != nil {
				return err
			}

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
