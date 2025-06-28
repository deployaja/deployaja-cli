package cmd

import (
	"fmt"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd())
}

func restartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart <DeploymentName>",
		Short: "Restart a deployment",
		Long:  "Restart a deployment by deleting and recreating its pods",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			deploymentName := args[0]

			fmt.Printf("%s Restarting deployment %s...\n", ui.InfoPrint("🔄"), deploymentName)

			response, err := apiClient.Restart(deploymentName)
			if err != nil {
				return err
			}

			// Display the response
			fmt.Printf("%s %s\n", ui.SuccessPrint("✓"), response.Message)
			fmt.Printf("Status: %s\n", response.Status)
			fmt.Printf("Total pods: %d\n", response.TotalPods)
			fmt.Printf("Successful deletions: %d\n", response.SuccessfulDeletions)

			if response.FailedDeletions > 0 {
				fmt.Printf("%s Failed deletions: %d\n", ui.ErrorPrint("⚠"), response.FailedDeletions)
			}

			return nil
		},
	}

	return cmd
}
