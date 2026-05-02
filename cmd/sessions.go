package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
	"github.com/justanoobcoder/tmux-mgr/internal/resurrect"
	"github.com/justanoobcoder/tmux-mgr/internal/service"
	"github.com/justanoobcoder/tmux-mgr/internal/tmux"
	"github.com/justanoobcoder/tmux-mgr/internal/ui"
	"github.com/spf13/cobra"
)

func sessionsRun(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.Resurrect.Enabled {
		fmt.Println("Resurrect management is disabled in config.")
		return nil
	}

	store := resurrect.NewStore(cfg.Resurrect.SaveDir)
	manager := service.NewManager(cfg, store)
	tmuxClient := tmux.NewClient()

	sessionMap := make(map[string]*domain.DisplaySession)
	var sessionNames []string

	activeSessions, lastActive, _ := tmuxClient.ListSessions()
	for _, name := range activeSessions {
		sessionMap[name] = &domain.DisplaySession{
			Name:         name,
			IsActive:     true,
			IsLastActive: name == lastActive,
		}
		sessionNames = append(sessionNames, name)
	}

	savedSessions, lastSaved, _ := store.ListSessions()
	for _, sess := range savedSessions {
		if ds, exists := sessionMap[sess.Name]; exists {
			ds.IsSaved = true
			if sess.Name == lastSaved {
				ds.IsLastSaved = true
			}
		} else {
			sessionMap[sess.Name] = &domain.DisplaySession{
				Name:        sess.Name,
				IsSaved:     true,
				IsLastSaved: sess.Name == lastSaved,
			}
			sessionNames = append(sessionNames, sess.Name)
		}
	}

	var displaySessions []domain.DisplaySession
	for _, name := range sessionNames {
		displaySessions = append(displaySessions, *sessionMap[name])
	}

	m := ui.NewSessionPickerModel(displaySessions)
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	if m.ToAttach != nil {
		if m.ToAttach.IsActive {
			if err := tmuxClient.Attach(m.ToAttach.Name); err != nil {
				return fmt.Errorf("attach active session: %w", err)
			}
		} else {
			if err := tmuxClient.RestoreResurrect(); err != nil {
				if err := tmuxClient.NewSession(m.ToAttach.Name, "."); err != nil {
					return fmt.Errorf("create new session: %w", err)
				}
			}
			if err := tmuxClient.Attach(m.ToAttach.Name); err != nil {
				return fmt.Errorf("attach new session: %w", err)
			}
		}
	} else if m.ToDelete != nil {
		if m.ToDelete.IsSaved {
			if err := manager.DeleteSession(m.ToDelete.Name); err != nil {
				return fmt.Errorf("delete session: %w", err)
			}
		}

		if m.ToDelete.IsActive {
			if err := tmuxClient.KillSession(m.ToDelete.Name); err != nil {
				fmt.Printf("Failed to kill active tmux session '%s': %v\n", m.ToDelete.Name, err)
			}
		}

		fmt.Printf("Deleted session '%s'\n", m.ToDelete.Name)
	}

	return nil
}

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage running and saved tmux sessions",
	Long: `Opens a unified TUI picker merging both your actively running tmux sessions
and your dead-but-saved sessions stored on disk by tmux-resurrect.
Includes smart tags to easily identify your [Last Active] and [Last Saved] contexts.`,
	RunE: sessionsRun,
}

func init() {
	rootCmd.AddCommand(sessionsCmd)
}
