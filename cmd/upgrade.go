package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade DeployAja CLI to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		var upgradeCmd *exec.Cmd
		var manualCmd, osName string

		if runtime.GOOS != "windows" {
			sudoCheck := exec.Command("sudo", "-v")
			sudoCheck.Stdout = cmd.OutOrStdout()
			sudoCheck.Stderr = cmd.OutOrStderr()
			if err := sudoCheck.Run(); err != nil {
				fmt.Println("Sudo access is required to upgrade.")
				return
			}
		}

		switch runtime.GOOS {
		case "windows":
			fmt.Println("Running upgrade for Windows...")
			upgradeCmd = exec.Command("powershell", "-Command", "iwr -useb https://deployaja.id/setup.bat | iex")
			manualCmd = "iwr -useb https://deployaja.id/setup.bat | iex"
			osName = "Windows"
		default:
			fmt.Println("Running upgrade for macOS/Linux...")
			upgradeCmd = exec.Command("bash", "-c", "set -e -o pipefail; curl -sSL https://deployaja.id/setup.sh | bash")
			manualCmd = "curl -sSL https://deployaja.id/setup.sh | bash"
			osName = "macOS/Linux"
		}

		upgradeCmd.Stdout = cmd.OutOrStdout()
		upgradeCmd.Stderr = cmd.OutOrStderr()

		done := make(chan error)
		go func() {
			done <- upgradeCmd.Run()
		}()

		spinChars := []rune{'|', '/', '-', '\\'}
		i := 0
	loop:
		for {
			select {
			case err := <-done:
				if err != nil {
					fmt.Println("\nUpgrade failed.")
					fmt.Printf("Please upgrade manually using the following command for %s:\n", osName)
					fmt.Printf("\n%s\n", manualCmd)
					fmt.Println("\nOr visit https://deployaja.id/ for the guide.")
				} else {
					fmt.Println("\nUpgrade successful!")
				}
				break loop
			default:
				fmt.Printf("\rUpgrading... %c", spinChars[i%len(spinChars)])
				time.Sleep(200 * time.Millisecond)
				i++
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
