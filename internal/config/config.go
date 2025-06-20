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
	if _, err := os.Stat(DeployFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("deployaja.yaml not found. Run 'deployaja init' to create one")
	}

	data, err := os.ReadFile(DeployFile)
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
