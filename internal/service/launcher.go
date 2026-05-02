package service

import (
	"fmt"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
	"github.com/justanoobcoder/tmux-mgr/internal/resurrect"
	"github.com/justanoobcoder/tmux-mgr/internal/tmux"
)

type Launcher struct {
	tmuxClient *tmux.Client
	config     *config.TmuxConfig
	store      *resurrect.Store
}

func NewLauncher(client *tmux.Client, cfg *config.TmuxConfig, store *resurrect.Store) *Launcher {
	return &Launcher{tmuxClient: client, config: cfg, store: store}
}

func (l *Launcher) Launch(p domain.Project) error {
	sessionName := p.SessionName(l.config.SessionPrefix)

	if !l.tmuxClient.HasSession(sessionName) {
		isSaved := false
		if l.store != nil {
			savedSessions, _, err := l.store.ListSessions()
			if err == nil {
				for _, s := range savedSessions {
					if s.Name == sessionName {
						isSaved = true
						break
					}
				}
			}
		}

		if isSaved {
			if err := l.tmuxClient.RestoreResurrect(); err != nil {
				if err := l.tmuxClient.NewSession(sessionName, p.Path); err != nil {
					return fmt.Errorf("create session: %w", err)
				}
			}
		} else {
			if err := l.tmuxClient.NewSession(sessionName, p.Path); err != nil {
				return fmt.Errorf("create session: %w", err)
			}
		}
	}

	if l.config.AttachOnCreate {
		return l.tmuxClient.Attach(sessionName)
	}

	return nil
}
