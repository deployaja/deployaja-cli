package cmd

import (
	"fmt"
	"time"

	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loginCmd())
}

func loginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Authenticate with DeployAja platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessionCode := uuid.New().String()
			loginURL := fmt.Sprintf("%s/login?ses=%s", apiClient.BaseURL, sessionCode)

			fmt.Printf("%s Opening browser for authentication...\n", ui.InfoPrint("🔐"))
			fmt.Printf("If browser doesn't open, visit: %s\n", loginURL)

			err := browser.OpenURL(loginURL)
			if err != nil {
				fmt.Printf("%s Failed to open browser: %v\n", ui.WarningPrint("⚠"), err)
			}

			fmt.Printf("%s Waiting for authentication...\n", ui.InfoPrint("⏳"))

			// Poll for authentication
			for i := 0; i < 120; i++ { // 2 minutes timeout
				time.Sleep(1 * time.Second)

				token, err := apiClient.CheckAuth(sessionCode)
				if err == nil && token != "" {
					err = config.SaveToken(token)
					if err != nil {
						return fmt.Errorf("failed to save token: %v", err)
					}

					apiClient.Token = token
					fmt.Printf("%s Authentication successful!\n", ui.SuccessPrint("✓"))
					return nil
				}

				if i%5 == 0 {
					fmt.Print(".")
				}
			}

			return fmt.Errorf("authentication timeout")
		},
	}
}
