package cmd

import (
	"fmt"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tmux-mgr configuration",
	Long:  "Commands to initialize and view the configuration file located at ~/.config/tmux-mgr/config.json",
}

func configInitRun(cmd *cobra.Command, args []string) error {
	cfg := &config.Config{
		Projects: []string{},
	}

	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println("Config initialized.")
	return nil
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default config",
	Long:  "Generates a default config file at ~/.config/tmux-mgr/config.json if one does not already exist.",
	RunE:  configInitRun,
}

func configShowRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("Projects:\n")
	for _, p := range cfg.Projects {
		fmt.Printf("  - %s\n", p)
	}
	fmt.Printf("Tmux.SessionPrefix: %s\n", cfg.Tmux.SessionPrefix)
	fmt.Printf("Tmux.AttachOnCreate: %v\n", cfg.Tmux.AttachOnCreate)

	return nil
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current config",
	Long:  "Prints the current, fully-parsed configuration to the terminal.",
	RunE:  configShowRun,
}

func init() {
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)

	rootCmd.AddCommand(configCmd)
}
