package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func rootRun(cmd *cobra.Command, args []string) error {
	return projectsRun(cmd, args)
}

func checkTmuxInstalled(cmd *cobra.Command, args []string) {
	_, err := exec.LookPath("tmuxs")
	if err != nil {
		fmt.Fprintln(os.Stderr, "tmux is not installed, please install it before using tmux-mgr.")
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "tmux-mgr",
	Short: "A CLI/TUI to manage tmux sessions and projects",
	Long: `tmux-mgr is a tmux session manager that automates how you create and attach to tmux sessions.
It replaces manual commands with a clean terminal interface for faster navigation.`,
	PersistentPreRun: checkTmuxInstalled,
	RunE:             rootRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
