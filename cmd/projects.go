package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/justanoobcoder/tmux-mgr/internal/resurrect"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/justanoobcoder/tmux-mgr/internal/tmux"
	"github.com/justanoobcoder/tmux-mgr/internal/ui"
	"github.com/spf13/cobra"
)

func projectsRun(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	store := resurrect.NewStore(cfg.Resurrect.SaveDir)
	manager := service.NewManager(cfg, store)

	projects, err := manager.GetProjects()
	if err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	m := ui.NewProjectPickerModel(projects)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	if toDelete := m.ProjectToDelete(); toDelete != nil {
		if err := manager.RemoveProject(toDelete.Path); err != nil {
			return fmt.Errorf("remove project: %w", err)
		}
		fmt.Printf("Removed project: %s\n", toDelete.Path)
		return nil
	}

	selected := m.SelectedProject()
	if selected == nil {
		return nil
	}

	if err := manager.TrackProjectSelection(selected.Path); err != nil {
		fmt.Printf("Warning: failed to track project selection: %v\n", err)
	}

	tmuxClient := tmux.NewClient()
	launcher := service.NewLauncher(tmuxClient, &cfg.Tmux, store)

	return launcher.Launch(*selected)
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Open the TUI project picker",
	Long:  `Opens a fuzzy-searchable TUI to select and launch your configured projects.`,
	RunE:  projectsRun,
}

func init() {
	rootCmd.AddCommand(projectsCmd)
}
