package cmd

import (
	"fmt"
	"os"
	"os/exec"	

	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(genCmd())
}

func genCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen [prompt]",
		Short: "Generate aja configuration based on a prompt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			prompt := args[0]

			fmt.Printf("%s Generating content...\n", ui.InfoPrint("🤖"))

			response, err := apiClient.Gen(prompt)
			if err != nil {
				return fmt.Errorf("failed to generate content: %v", err)
			}

			// Create temporary YAML file with random 2-digit number in filename
			tempFile := "deployaja.yaml"

			err = os.WriteFile(tempFile, []byte(response.Content), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file: %v", err)
			}

			fmt.Printf("%s Content written to %s\n", ui.SuccessPrint("✓"), tempFile)
			fmt.Printf("%s Opening vim...\n", ui.InfoPrint("📝"))

			// Open vim with the temporary file
			vimCmd := exec.Command("vim", tempFile)
			vimCmd.Stdin = os.Stdin
			vimCmd.Stdout = os.Stdout
			vimCmd.Stderr = os.Stderr

			err = vimCmd.Run()
			if err != nil {
				return fmt.Errorf("failed to open vim: %v", err)
			}

			return nil
		},
	}
}
