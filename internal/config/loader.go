package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get home directory: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(home, ".config", "tmux-mgr", "config.json")
}

func ConfigExists() bool {
	config := GetConfigFilePath()
	_, err := os.Stat(config)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Fprintf(os.Stderr, "Error checking for config file: %v\n", err)
		os.Exit(1)
	}
	return true
}

func GetResurrectSaveDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get home directory: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(home, ".tmux", "resurrect")
}

func Load() (*Config, error) {
	if !ConfigExists() {
		fmt.Fprintf(os.Stderr, "Config file does not exist.\n")
		os.Exit(1)
	}

	configDir := GetConfigFilePath()
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.SetDefault("tmux.attach_on_create", true)
	viper.SetDefault("resurrect.enabled", true)
	viper.SetDefault("resurrect.save_dir", GetResurrectSaveDir())

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
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
