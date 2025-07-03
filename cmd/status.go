package cmd

import (
	"fmt"

	"deployaja-cli/internal/api"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd())
}

func statusCmd() *cobra.Command {
	var detailed bool

	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"ls"},
		Short:   "Check deployment status and health",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			response, err := apiClient.GetStatus()
			if err != nil {
				return err
			}

			if len(response.Deployments) == 0 {
				fmt.Printf("%s No deployments found\n", ui.WarningPrint("⚠"))
				return nil
			}

			fmt.Printf("%s Deployment Status\n\n", ui.InfoPrint("📊"))

			// Prepare table data
			headers := []string{"NAME", "STATUS", "REPLICAS", "READY", "URL", "LAST DEPLOYED"}
			var rows [][]string

			for _, deployment := range response.Deployments {
				statusColor := ui.GetStatusColor(deployment.Status)
				statusText := statusColor(deployment.Status)

				// Use new replica fields if available, fallback to old format
				var replicas, ready string
				if deployment.DesiredReplicas > 0 || deployment.AvailableReplicas > 0 {
					replicas = fmt.Sprintf("%d/%d", deployment.AvailableReplicas, deployment.DesiredReplicas)
					ready = fmt.Sprintf("%d/%d", deployment.ReadyReplicas, deployment.DesiredReplicas)
				} else {
					replicas = fmt.Sprintf("%d/%d", deployment.Replicas.Available, deployment.Replicas.Desired)
					ready = fmt.Sprintf("%d/%d", deployment.Replicas.Available, deployment.Replicas.Desired)
				}

				url := deployment.URL
				if url == "" {
					url = "-"
				}

				lastDeployed := ui.FormatTime(deployment.LastDeployed)

				rows = append(rows, []string{
					deployment.Name,
					statusText,
					replicas,
					ready,
					url,
					lastDeployed,
				})
			}

			// Print table
			fmt.Print(ui.FormatTable(headers, rows))

			// Show detailed pod information if requested or if there are issues
			if detailed || hasIssues(response.Deployments) {
				fmt.Printf("\n%s Pod Details\n\n", ui.InfoPrint("🔍"))

				for _, deployment := range response.Deployments {
					if len(deployment.Pods) > 0 {
						fmt.Printf("%s %s\n", ui.InfoPrint("📦"), deployment.Name)

						// Pod table
						podHeaders := []string{"POD NAME", "STATUS", "READY", "RESTARTS", "AGE"}
						var podRows [][]string

						for _, pod := range deployment.Pods {
							statusColor := ui.GetStatusColor(pod.Status)
							statusText := statusColor(pod.Status)

							readyText := "No"
							if pod.Ready {
								readyText = ui.SuccessPrint("Yes")
							} else {
								readyText = ui.ErrorPrint("No")
							}

							restarts := fmt.Sprintf("%d", pod.RestartCount)
							if pod.RestartCount > 0 {
								restarts = ui.WarningPrint(restarts)
							}

							podRows = append(podRows, []string{
								pod.Name,
								statusText,
								readyText,
								restarts,
								pod.Age,
							})
						}

						fmt.Print(ui.FormatTable(podHeaders, podRows))

						// Show container details for problematic pods
						for _, pod := range deployment.Pods {
							if !pod.Ready || pod.RestartCount > 0 || pod.Status == "CRASH_LOOP" || pod.Status == "FAILED" || pod.Status == "ERROR" {
								fmt.Printf("\n%s Pod: %s\n", ui.WarningPrint("⚠"), pod.Name)

								if pod.Reason != "" {
									fmt.Printf("  Reason: %s\n", ui.ErrorPrint(pod.Reason))
								}

								if pod.Message != "" {
									fmt.Printf("  Message: %s\n", ui.ErrorPrint(pod.Message))
								}

								if len(pod.ContainerStatuses) > 0 {
									fmt.Printf("  Containers:\n")
									for _, container := range pod.ContainerStatuses {
										statusIcon := "✓"
										statusColor := ui.SuccessPrint
										if !container.Ready {
											statusIcon = "✗"
											statusColor = ui.ErrorPrint
										}

										fmt.Printf("    %s %s: %s", statusColor(statusIcon), container.Name, container.State)

										if container.RestartCount > 0 {
											fmt.Printf(" (restarts: %s)", ui.WarningPrint(fmt.Sprintf("%d", container.RestartCount)))
										}

										if container.Reason != "" {
											fmt.Printf(" - %s", ui.ErrorPrint(container.Reason))
										}

										fmt.Printf("\n")

										if container.Message != "" {
											fmt.Printf("      %s\n", ui.ErrorPrint(container.Message))
										}
									}
								}
								fmt.Printf("\n")
							}
						}
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&detailed, "detailed", "d", false, "Show detailed pod information")

	return cmd
}

// hasIssues checks if any deployment has problematic pods
func hasIssues(deployments []api.DeploymentStatus) bool {
	for _, deployment := range deployments {
		for _, pod := range deployment.Pods {
			if !pod.Ready || pod.RestartCount > 0 || pod.Status == "CRASH_LOOP" || pod.Status == "FAILED" || pod.Status == "ERROR" {
				return true
			}
		}
	}
	return false
}
