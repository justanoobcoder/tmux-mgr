package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/spf13/cobra"
)

var force bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tmux-mgr configuration",
	Long:  "Commands to initialize and view the configuration file located at ~/.config/tmux-mgr/config.json",
}

func loadConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrConfigNotFound) {
			fmt.Fprintln(os.Stderr, "Error: config file not found; run 'tmux-mgr config init' to create one")
			configInitCmd.Help()
			os.Exit(1)
		}
		return nil, err
	}
	return cfg, nil
}

func configInitRun(cmd *cobra.Command, args []string) error {
	if !force {
		exists, err := config.ConfigExists()
		if err != nil {
			return err
		}
		if exists {
			fmt.Println("Config file already exists (use --force to overwrite)")
			return nil
		}
	}

	resurrectSaveDir, err := config.GetResurrectSaveDir()
	if err != nil {
		return err
	}

	cfg := &config.Config{
		Projects: []string{},
		Folders:  []config.FolderConfig{},
		Tmux: config.TmuxConfig{
			SessionPrefix:  "",
			AttachOnCreate: true,
		},
		Resurrect: config.ResurrectConfig{
			Enabled: true,
			SaveDir: resurrectSaveDir,
		},
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
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	fmt.Printf("Projects:\n")
	for _, p := range cfg.Projects {
		fmt.Printf("  - %s\n", p)
	}
	fmt.Printf("Folders:\n")
	for _, f := range cfg.Folders {
		fmt.Printf("  - %s (excludes: %v)\n", f.Path, f.Excludes)
	}
	fmt.Printf("Tmux.SessionPrefix: %s\n", cfg.Tmux.SessionPrefix)
	fmt.Printf("Tmux.AttachOnCreate: %v\n", cfg.Tmux.AttachOnCreate)
	fmt.Printf("Resurrect.Enabled: %v\n", cfg.Resurrect.Enabled)
	fmt.Printf("Resurrect.SaveDir: %s\n", cfg.Resurrect.SaveDir)

	return nil
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current config",
	Long:  "Prints the current, fully-parsed configuration to the terminal.",
	RunE:  configShowRun,
}

func init() {
	configInitCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing config file")
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)

	rootCmd.AddCommand(configCmd)
}
