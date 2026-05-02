package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/justanoobcoder/tmux-mgr/internal/resurrect"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/spf13/cobra"
)

var (
	addFolderFlag   bool
	addExcludeFlags []string
)

func addRun(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	store := resurrect.NewStore(cfg.Resurrect.SaveDir)
	manager := service.NewManager(cfg, store)

	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	if addFolderFlag {
		if err := manager.AddFolder(path, addExcludeFlags); err != nil {
			return fmt.Errorf("add folder: %w", err)
		}
		absPath, _ := filepath.Abs(path)
		fmt.Printf("Added folder: %s (excludes: %v)\n", filepath.Clean(absPath), addExcludeFlags)
	} else {
		if err := manager.AddProject(path); err != nil {
			return fmt.Errorf("add project: %w", err)
		}
		absPath, _ := filepath.Abs(path)
		fmt.Printf("Added project: %s\n", filepath.Clean(absPath))
	}
	return nil
}

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a project or folder to your configuration",
	Long: `Adds a single project or an entire parent folder to your configuration.
If you use the --folder flag, it will dynamically scan all subdirectories as projects.`,
	Example: `  tmux-mgr add  # add current directory
  tmux-mgr add /path/to/project
  tmux-mgr add --folder /path/to/workspace
  tmux-mgr add --folder /path/to/workspace --exclude node_modules`,
	RunE: addRun,
}

func init() {
	addCmd.Flags().BoolVarP(&addFolderFlag, "folder", "f", false, "Add as a dynamic folder instead of a single project")
	addCmd.Flags().StringSliceVarP(&addExcludeFlags, "exclude", "e", []string{}, "Exclude subdirectories (only valid with --folder)")

	rootCmd.AddCommand(addCmd)
}
