package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var ErrConfigNotFound = errors.New("config file not found")

func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "tmux-mgr", "config.json"), nil
}

func ConfigExists() (bool, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("check config file: %w", err)
	}
	return true, nil
}

func GetResurrectSaveDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(home, ".tmux", "resurrect"), nil
}

func Load() (*Config, error) {
	exists, err := ConfigExists()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrConfigNotFound
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	saveDir, err := GetResurrectSaveDir()
	if err != nil {
		return nil, err
	}

	viper.SetConfigFile(configPath)

	viper.SetDefault("tmux.attach_on_create", true)
	viper.SetDefault("resurrect.enabled", true)
	viper.SetDefault("resurrect.save_dir", saveDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	configDir := filepath.Join(home, ".config", "tmux-mgr")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}
