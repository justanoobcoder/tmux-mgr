package cmd

import (
	"fmt"

	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/spf13/cobra"
)

func removeRun(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	manager := service.NewManager(cfg, nil)

	path := args[0]
	if err := manager.RemoveConfigPath(path); err != nil {
		return fmt.Errorf("remove path: %w", err)
	}

	fmt.Printf("Successfully removed '%s' from configuration.\n", path)
	return nil
}

var removeCmd = &cobra.Command{
	Use:   "remove <path>",
	Short: "Remove a configured project",
	Long:  `Remove a configured project from the configuration.`,
	Example: `  tmux-mgr remove /path/to/project
  tmux-mgr remove /path/to/workspace`,
	Args: cobra.ExactArgs(1),
	RunE: removeRun,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
