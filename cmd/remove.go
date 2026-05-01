package cmd

import (
	"fmt"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/spf13/cobra"
)

func removeRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	manager := service.NewManager(cfg)

	path := args[0]
	if err := manager.RemoveConfigPath(path); err != nil {
		return err
	}

	fmt.Printf("Successfully removed '%s' from configuration.\n", path)
	return nil
}

var removeCmd = &cobra.Command{
	Use:     "remove <path>",
	Short:   "Remove a configured project",
	Long:    `Remove a configured project from the configuration.`,
	Example: `  tmux-mgr remove /path/to/project`,
	Args:    cobra.ExactArgs(1),
	RunE:    removeRun,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
