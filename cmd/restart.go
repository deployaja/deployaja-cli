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
			if !response.Success {
				return fmt.Errorf("restart failed")
			}

			fmt.Printf("%s %s\n", ui.SuccessPrint("✓"), response.Data.Message)
			fmt.Printf("Status: %s\n", response.Data.Status)
			fmt.Printf("Method: %s\n", response.Data.Method)

			// Display rollout status
			rollout := response.Data.RolloutStatus
			fmt.Printf("Rollout Status:\n")
			fmt.Printf("  Generation: %d\n", rollout.Generation)
			fmt.Printf("  Observed Generation: %d\n", rollout.ObservedGeneration)
			fmt.Printf("  Replicas: %d\n", rollout.Replicas)
			fmt.Printf("  Ready Replicas: %d\n", rollout.ReadyReplicas)
			fmt.Printf("  Updated Replicas: %d\n", rollout.UpdatedReplicas)

			return nil
		},
	}

	return cmd
}
