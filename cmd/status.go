package cmd

import (
	"fmt"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd())
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [NAME]",
		Short: "Check deployment status and health",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			var name string
			if len(args) > 0 {
				name = args[0]
			}

			response, err := apiClient.GetStatus(name)
			if err != nil {
				return err
			}

			if len(response.Deployments) == 0 {
				fmt.Printf("%s No deployments found\n", ui.WarningPrint("⚠"))
				return nil
			}

			fmt.Printf("%s Deployment Status\n\n", ui.InfoPrint("📊"))

			// Prepare table data
			headers := []string{"NAME", "STATUS", "REPLICAS", "URL", "LAST DEPLOYED"}
			var rows [][]string

			for _, deployment := range response.Deployments {
				statusColor := ui.GetStatusColor(deployment.Status)
				statusText := statusColor(deployment.Status)

				replicas := fmt.Sprintf("%d/%d", deployment.Replicas.Available, deployment.Replicas.Desired)
				url := deployment.URL
				if url == "" {
					url = "-"
				}

				lastDeployed := ui.FormatTime(deployment.LastDeployed)

				rows = append(rows, []string{
					deployment.Name,
					statusText,
					replicas,
					url,
					lastDeployed,
				})
			}

			// Print table
			fmt.Print(ui.FormatTable(headers, rows))

			return nil
		},
	}
}
