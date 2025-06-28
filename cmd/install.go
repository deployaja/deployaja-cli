package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd())
}

func installCmd() *cobra.Command {
	var domain string
	var dryRun bool
	var name string

	cmd := &cobra.Command{
		Use:   "install [APPNAME]",
		Short: "Install an app from the marketplace",
		Long: `Install an app from the marketplace by downloading its configuration.
The configuration will be saved as APPNAME-install.json in the current directory.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			appName := args[0]

			fmt.Printf("%s Installing %s from marketplace...\n", ui.InfoPrint("📦"), appName)
			if domain != "" {
				fmt.Printf("%s Using custom domain: %s\n", ui.InfoPrint("🌐"), domain)
			}
			if dryRun {
				fmt.Printf("%s Dry run mode enabled\n", ui.InfoPrint("🔍"))
			}

			// Get app configuration from API
			response, err := apiClient.InstallApp(appName, domain, name, dryRun)
			if err != nil {
				return fmt.Errorf("failed to install app: %v", err)
			}
			// Check for successful installation (accept "success", "initiated", "pending" as successful)
			successStatuses := []string{"success", "initiated", "pending", "deploying", "running", "validated"}
			isSuccess := false
			for _, status := range successStatuses {
				if response.Status == status {
					isSuccess = true
					break
				}
			}

			if !isSuccess {
				return fmt.Errorf("installation failed: %s (status: %s)", response.Message, response.Status)
			}

			// Decode the base64 config
			configData, err := base64.StdEncoding.DecodeString(response.Config)
			if err != nil {
				return fmt.Errorf("failed to decode configuration: %v", err)
			}

			// Save config to YAML file
			filename := fmt.Sprintf("%s.yaml", response.DeploymentName)
			err = os.WriteFile(filename, configData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write config file: %v", err)
			}

			absPath, _ := filepath.Abs(filename)
			fmt.Printf("%s Configuration saved to: %s\n", ui.SuccessPrint("✅"), absPath)
			fmt.Printf("%s %s\n", ui.InfoPrint("💡"), response.Message)			
			

			// Poll for deployment status until completion
			fmt.Printf("%s Waiting for deployment to complete...\n", ui.InfoPrint("🔍"))

			var lastStatus string
			statusCallback := func(status string) {
				if status != lastStatus {
					fmt.Printf("%s Status: %s\n", ui.InfoPrint("📊"), status)
					lastStatus = status
				}
			}

			finalDeployment, err := apiClient.PollDeploymentStatus(response.DeploymentName, statusCallback)
			if err != nil {
				fmt.Printf("%s Warning: Failed to monitor deployment status: %v\n", ui.WarningPrint("⚠️"), err)
				fmt.Printf("%s You can check the status manually using: deployaja status\n", ui.InfoPrint("💡"))
				return nil
			}

			// Show final status
			if finalDeployment.Status == "running" || finalDeployment.Status == "success" {
				fmt.Printf("%s Installation completed successfully!\n", ui.SuccessPrint("🎉"))
				if finalDeployment.URL != "" {
					fmt.Printf("%s Access your application at: %s\n", ui.InfoPrint("🌐"), finalDeployment.URL)
				}
			} else {
				fmt.Printf("%s Installation failed with status: %s\n", ui.ErrorPrint("❌"), finalDeployment.Status)
				fmt.Printf("%s Use 'deployaja describe %s' for more details\n", ui.InfoPrint("💡"), response.DeploymentName)
				return fmt.Errorf("installation failed")
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&domain, "domain", "d", "", "Custom domain for the ingress URL")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Perform a dry run without actually installing")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Custom name for the deployment")

	return cmd
}
