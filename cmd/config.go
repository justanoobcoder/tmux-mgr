package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/justanoobcoder/tmux-mgr/internal/ui"
	"github.com/spf13/cobra"
)

var force bool
var editorFlag string

func configRun(cmd *cobra.Command, args []string) error {
	exists, err := config.ConfigExists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Fprintln(os.Stderr, "Error: config file not found; run 'tmux-mgr config init' to create one")
		configInitCmd.Help()
		return nil
	}

	path, err := config.GetConfigFilePath()
	if err != nil {
		return err
	}

	editor := editorFlag
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		return fmt.Errorf("$EDITOR is not set; use --editor <editor> to specify one")
	}

	editorCmd := exec.Command(editor, path)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	return nil
}

func configPruneRun(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	manager := service.NewManager(cfg, nil)
	deadPaths, err := manager.GetDeadPaths()
	if err != nil {
		return err
	}

	if len(deadPaths) == 0 {
		fmt.Println("No dead paths found in configuration.")
		return nil
	}

	var pathsToRemove []string

	if force {
		pathsToRemove = deadPaths
	} else {
		m := ui.NewPrunePickerModel(deadPaths)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running TUI: %w", err)
		}

		if !m.IsConfirmed() {
			fmt.Println("Pruning cancelled.")
			return nil
		}
		pathsToRemove = m.GetSelectedPaths()
	}

	if len(pathsToRemove) == 0 {
		fmt.Println("No paths selected for removal.")
		return nil
	}

	if err := manager.BulkRemoveConfigPaths(pathsToRemove); err != nil {
		return fmt.Errorf("bulk remove paths: %w", err)
	}

	fmt.Printf("Successfully removed %d paths from configuration.\n", len(pathsToRemove))
	for _, p := range pathsToRemove {
		fmt.Printf("  - %s\n", p)
	}

	return nil
}

var configPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove non-existent paths from configuration",
	Long:  "Identifies projects and folders in the configuration that no longer exist on the filesystem and allows you to remove them.",
	RunE:  configPruneRun,
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tmux-mgr configuration",
	Long:  "Commands to initialize and view the configuration file located at ~/.config/tmux-mgr/config.json.\nRunning this command without subcommands will open the config file in your $EDITOR.",
	RunE:  configRun,
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
	configCmd.Flags().StringVarP(&editorFlag, "editor", "e", "", "Editor to open config file with (overrides $EDITOR)")
	configInitCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing config file")
	configPruneCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip interactive TUI and remove all dead paths")
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configPruneCmd)

	rootCmd.AddCommand(configCmd)
}
