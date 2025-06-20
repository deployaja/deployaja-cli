package cmd

import (
	"encoding/base64"
	"encoding/json"
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

			// Get app configuration from API
			response, err := apiClient.InstallApp(appName)
			if err != nil {
				return fmt.Errorf("failed to install app: %v", err)
			}

			if response.Status != "success" {
				return fmt.Errorf("installation failed: %s", response.Message)
			}

			// Decode the base64 config
			configData, err := base64.StdEncoding.DecodeString(response.Config)
			if err != nil {
				return fmt.Errorf("failed to decode configuration: %v", err)
			}

			// Create install data structure
			installData := map[string]interface{}{
				"appName":    response.AppName,
				"config":     string(configData),
				"message":    response.Message,
				"status":     response.Status,
				"installUrl": response.InstallURL,
			}

			// Save to JSON file
			filename := fmt.Sprintf("%s-install.json", appName)
			jsonData, err := json.MarshalIndent(installData, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal JSON: %v", err)
			}

			err = os.WriteFile(filename, jsonData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %v", err)
			}

			absPath, _ := filepath.Abs(filename)
			fmt.Printf("%s Configuration saved to: %s\n", ui.SuccessPrint("✅"), absPath)
			fmt.Printf("%s %s\n", ui.InfoPrint("💡"), response.Message)

			if response.InstallURL != "" {
				fmt.Printf("%s Install URL: %s\n", ui.InfoPrint("🔗"), response.InstallURL)
			}

			return nil
		},
	}

	return cmd
}
