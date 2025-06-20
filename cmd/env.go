package cmd

import (
	"deployaja-cli/internal/ui"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd())
}

func envCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env [edit|set|get]",
		Short: "Manage environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ensureAuthenticated(); err != nil {
				return err
			}

			action := "edit"
			if len(args) > 0 {
				action = args[0]
			}

			switch action {
			case "edit":
				return editEnvVars()
			case "get":
				key := ""
				if len(args) > 1 {
					key = args[1]
				}
				return getEnvVars(key)
			case "set":
				if len(args) < 2 {
					return fmt.Errorf("usage: deployaja env set KEY=VALUE")
				}
				return setEnvVar(args[1])
			default:
				return fmt.Errorf("unknown action: %s", action)
			}
		},
	}

	return cmd
}

func editEnvVars() error {
	// Get current env vars
	envVars, err := apiClient.GetEnvVars()
	if err != nil {
		return err
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "deployaja-env-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Write env vars to temp file
	data, err := json.MarshalIndent(envVars, "", "  ")
	if err != nil {
		return err
	}

	_, err = tmpFile.Write(data)
	if err != nil {
		return err
	}
	tmpFile.Close()

	// Open in vim
	cmd := exec.Command("vim", tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	// Read modified content
	modifiedData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return err
	}

	var modifiedEnvVars map[string]string
	err = json.Unmarshal(modifiedData, &modifiedEnvVars)
	if err != nil {
		return err
	}

	// Update env vars
	return apiClient.UpdateEnvVars(modifiedEnvVars)
}

func getEnvVars(key string) error {
	envVars, err := apiClient.GetEnvVars()
	if err != nil {
		return err
	}

	if key != "" {
		if value, exists := envVars[key]; exists {
			fmt.Printf("%s\n", value)
		} else {
			return fmt.Errorf("environment variable %s not found", key)
		}
	} else {
		for k, v := range envVars {
			fmt.Printf("%s=%s\n", k, v)
		}
	}

	return nil
}

func setEnvVar(keyValue string) error {
	parts := strings.SplitN(keyValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format. Use KEY=VALUE")
	}

	envVars := map[string]string{
		parts[0]: parts[1],
	}

	err := apiClient.UpdateEnvVars(envVars)
	if err != nil {
		return err
	}

	fmt.Printf("%s Set %s\n", ui.SuccessPrint("✓"), keyValue)
	return nil
}
