package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func rootRun(cmd *cobra.Command, args []string) error {
	return projectsRun(cmd, args)
}

var rootCmd = &cobra.Command{
	Use:   "tmux-mgr",
	Short: "A CLI/TUI to manage tmux sessions and projects",
	Long: `tmux-mgr is a tmux session manager that automates how you create and attach to tmux sessions.
It replaces manual commands with a clean terminal interface for faster navigation.`,
	RunE: rootRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
