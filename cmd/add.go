package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/spf13/cobra"
)

func addRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	manager := service.NewManager(cfg)

	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	if err := manager.AddProject(path); err != nil {
		return fmt.Errorf("add project: %w", err)
	}
	absPath, _ := filepath.Abs(path)
	fmt.Printf("Added project: %s\n", filepath.Clean(absPath))

	return nil
}

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a project to configuration",
	Long:  `Adds a single project to your configuration.`,
	Example: `  tmux-mgr add   # add current directory
  tmux-mgr add /path/to/project`,
	RunE: addRun,
}

func init() {
	rootCmd.AddCommand(addCmd)
}
