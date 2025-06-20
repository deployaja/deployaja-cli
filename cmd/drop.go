package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dropCmd())
}

func dropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drop NAME",
		Short: "Delete/destroy deployment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			name := args[0]
			force, _ := cmd.Flags().GetBool("force")

			if !force {
				fmt.Printf("%s Are you sure you want to delete %s? (y/N): ", ui.WarningPrint("⚠"), name)
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response != "y" && response != "yes" {
					fmt.Printf("Cancelled\n")
					return nil
				}
			}

			fmt.Printf("%s Deleting %s...\n", ui.InfoPrint("🗑"), name)

			err := apiClient.Drop(name)
			if err != nil {
				return err
			}

			fmt.Printf("%s Deletion initiated for %s\n", ui.SuccessPrint("✓"), name)

			return nil
		},
	}

	cmd.Flags().Bool("force", false, "Force deletion without confirmation")
	return cmd
}
