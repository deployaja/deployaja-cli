package cmd

import (
	"fmt"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd())
}

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all active deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			response, err := apiClient.ListDeployments()
			if err != nil {
				return err
			}

			if len(response.Deployments) == 0 {
				fmt.Printf("%s No deployments found\n", ui.WarningPrint("⚠"))
				return nil
			}

			fmt.Printf("%s Active Deployments\n\n", ui.InfoPrint("📦"))

			for _, deployment := range response.Deployments {
				statusColor := ui.GetStatusColor(deployment.Status)
				fmt.Printf("%-20s %s", deployment.Name, statusColor(deployment.Status))

				if deployment.URL != "" {
					fmt.Printf(" %s", deployment.URL)
				}

				fmt.Printf(" (created %s)\n", ui.FormatTime(deployment.CreatedAt))
			}

			return nil
		},
	}
}
