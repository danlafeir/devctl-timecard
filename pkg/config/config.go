package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	ConfigFileName = "config"
)

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine home directory: " + err.Error())
	}

	// Detect binary name to determine config location
	// If running as "timecard", use .timecard/config.yaml
	// Otherwise (devctl-timecard), use .devctl/config.yaml
	var execName string
	if len(os.Args) > 0 {
		execPath := os.Args[0]
		execName = filepath.Base(execPath)
		// Remove any extensions and check if it's the timecard binary
		execName = strings.TrimSuffix(execName, filepath.Ext(execName))
	}

	if execName == "timecard" {
		return filepath.Join(homeDir, ".timecard", "config.yaml")
	}

	// Default to devctl config location
	return filepath.Join(homeDir, ".devctl", "config.yaml")
}

func InitConfig() error {
	configFilePath := getConfigPath()
	configDir := filepath.Dir(configFilePath)
	configName := strings.TrimSuffix(filepath.Base(configFilePath), filepath.Ext(configFilePath))

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	// Create config file if it doesn't exist (with empty tempo structure)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Create empty config file with tempo key structure (tempo.* keys are used for Tempo API settings)
		emptyConfig := []byte("tempo:\n")
		if err := os.WriteFile(configFilePath, emptyConfig, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)
	return nil
}

func SaveConfig() error {
	configFilePath := getConfigPath()
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	return viper.WriteConfigAs(configFilePath)
}

func LoadConfig() error {
	return viper.ReadInConfig()
}
