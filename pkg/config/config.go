package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	ConfigDirName  = ".devctl/plugins/tempo"
	ConfigFileName = "config"
)

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("Could not determine home directory: " + err.Error())
	}
	return filepath.Join(home, ConfigDirName)
}

func InitConfig() error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	viper.AddConfigPath(configDir)
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType("yaml")
	return nil
}

func SaveConfig() error {
	configDir := getConfigDir()
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	configFile := filepath.Join(configDir, ConfigFileName+".yaml")
	return viper.WriteConfigAs(configFile)
}

func LoadConfig() error {
	return viper.ReadInConfig()
}
