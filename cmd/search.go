package cmd

import (
	"fmt"
	"strings"

	"deployaja-cli/internal/ui"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd())
}

func searchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [QUERY]",
		Short: "Search for apps in the marketplace",
		Long: `Search for applications in the marketplace.
You can search by app name, description, category, or tags.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			query := args[0]

			fmt.Printf("%s Searching for: %s\n\n", ui.InfoPrint("🔍"), query)

			// Search apps via API
			response, err := apiClient.SearchApps(query)
			if err != nil {
				return fmt.Errorf("failed to search apps: %v", err)
			}

			if len(response.Apps) == 0 {
				fmt.Printf("%s No apps found matching '%s'\n", ui.WarningPrint("⚠️"), query)
				return nil
			}

			fmt.Printf("%s Found %d apps\n\n", ui.SuccessPrint("✅"), response.Total)

			// Display results in a table format
			for i, app := range response.Apps {
				fmt.Printf("%s %s\n", color.New(color.Bold, color.FgCyan).Sprint(i+1), color.New(color.Bold).Sprint(app.Name))
				fmt.Printf("   %s\n", app.Description)
				fmt.Printf("   Category: %s\n", app.Category)
				fmt.Printf("   Author: %s\n", app.Author)
				fmt.Printf("   Version: %s\n", app.Version)
				fmt.Printf("   Downloads: %d\n", app.Downloads)
				fmt.Printf("   Rating: %.1f/5.0\n", app.Rating)

				if len(app.Tags) > 0 {
					fmt.Printf("   Tags: %s\n", strings.Join(app.Tags, ", "))
				}

				if app.Repository != "" {
					fmt.Printf("   Repository: %s\n", app.Repository)
				}

				fmt.Println()
			}

			fmt.Printf("%s Use 'aja install <app-name>' to install an app\n", ui.InfoPrint("💡"))

			return nil
		},
	}

	return cmd
}
