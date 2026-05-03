package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var tmuxNotRequired = map[string]bool{
	"version":    true,
	"help":       true,
	"completion": true,
}

func rootRun(cmd *cobra.Command, args []string) error {
	return projectsRun(cmd, args)
}

func checkTmuxInstallRun(cmd *cobra.Command, args []string) {
	if !tmuxNotRequired[cmd.Name()] {
		if _, err := exec.LookPath("tmux"); err != nil {
			fmt.Fprintln(os.Stderr, "tmux is not installed or not found in $PATH")
			os.Exit(1)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "tmux-mgr",
	Short: "A CLI/TUI to manage tmux sessions and projects",
	Long: `tmux-mgr is a tmux session manager that automates how you create and attach to tmux sessions.
It replaces manual commands with a clean terminal interface for faster navigation.`,
	PersistentPreRun: checkTmuxInstallRun,
	RunE:             rootRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
