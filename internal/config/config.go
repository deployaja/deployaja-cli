package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = ".deployaja"
	ConfigFile = "config.yaml"
	TokenFile  = "token"
	DeployFile = "deployaja.yaml"
)

func LoadDeploymentConfig() (*DeploymentConfig, error) {
	return LoadDeploymentConfigFromFile(DeployFile)
}

func LoadDeploymentConfigFromFile(filePath string) (*DeploymentConfig, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("deployment config file '%s' not found. Run 'aja init' to create one", filePath)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config DeploymentConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadToken() string {
	home, _ := os.UserHomeDir()
	tokenPath := filepath.Join(home, ConfigDir, TokenFile)

	if data, err := os.ReadFile(tokenPath); err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

func SaveToken(token string) error {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ConfigDir)
	os.MkdirAll(configPath, 0755)

	tokenPath := filepath.Join(configPath, TokenFile)
	return os.WriteFile(tokenPath, []byte(token), 0600)
}
