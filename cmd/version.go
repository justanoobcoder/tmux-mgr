package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is set at build time via:
//
//	go build -ldflags "-X github.com/justanoobcoder/tmux-mgr/cmd.version=v1.0.0"
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of tmux-mgr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tmux-mgr", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
