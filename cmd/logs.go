package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"deployaja-cli/internal/api"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logsCmd())
}

func logsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs NAME",
		Short: "View deployment logs",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			name := args[0]
			tail, _ := cmd.Flags().GetInt("tail")
			follow, _ := cmd.Flags().GetBool("follow")

			fmt.Printf("%s Fetching logs for %s...\n", ui.InfoPrint("📝"), name)

			if follow {
				return streamLogs(name, tail)
			}

			// Regular logs (non-follow mode)
			logs, err := apiClient.GetLogs(name, tail, false)
			if err != nil {
				return err
			}

			// Display logs
			for _, log := range logs {
				levelColor := ui.GetLogLevelColor(log.Level)
				fmt.Printf("[%s] %s %s\n",
					ui.FormatTime(log.Timestamp),
					levelColor(strings.ToUpper(log.Level)),
					log.Message)
			}

			return nil
		},
	}

	cmd.Flags().Int("tail", 100, "Number of lines to show")
	cmd.Flags().BoolP("follow", "f", false, "Follow log output")

	return cmd
}

func streamLogs(name string, tail int) error {
	logChan := make(chan api.LogEntry, 100)
	errorChan := make(chan error, 1)

	// Start streaming in a goroutine
	go apiClient.GetLogsStream(name, tail, logChan, errorChan)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("%s Following logs (press Ctrl+C to stop)...\n", ui.InfoPrint("🔄"))

	for {
		select {
		case log := <-logChan:
			levelColor := ui.GetLogLevelColor(log.Level)
			fmt.Printf("[%s] %s %s\n",
				ui.FormatTime(log.Timestamp),
				levelColor(strings.ToUpper(log.Level)),
				log.Message)

		case err := <-errorChan:
			if err != nil {
				return fmt.Errorf("stream error: %v", err)
			}
			return nil

		case <-sigChan:
			fmt.Printf("\n%s Stopping log stream...\n", ui.InfoPrint("⏹️"))
			return nil
		}
	}
}
