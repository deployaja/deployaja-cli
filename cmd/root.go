package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"deployaja-cli/internal/api"
	"deployaja-cli/internal/config"
	"deployaja-cli/internal/ui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	apiClient *api.APIClient
)

var rootCmd = &cobra.Command{
	Use:   "aja",
	Short: "Deploy applications with managed dependencies",
	Long: `DeployAja is a CLI tool that simplifies container deployment 
with managed dependencies like PostgreSQL, Redis, and more.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", ui.ErrorPrint("Error:"), err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.deployaja/config.yaml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		configPath := filepath.Join(home, config.ConfigDir)
		os.MkdirAll(configPath, 0755)

		viper.AddConfigPath(configPath)
		viper.SetConfigType("yaml")
		viper.SetConfigName(config.ConfigFile)
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()

	token := config.LoadToken()
	apiClient = api.NewApiClient(token)
}

func ensureAuthenticated() error {
	if apiClient.Token == "" {
		return fmt.Errorf("not authenticated. Run 'aja login' first")
	}
	return nil
}
