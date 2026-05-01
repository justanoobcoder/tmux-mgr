package service

import (
	"fmt"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
	"github.com/justanoobcoder/tmux-mgr/internal/tmux"
)

type Launcher struct {
	tmuxClient *tmux.Client
	config     *config.TmuxConfig
}

func NewLauncher(client *tmux.Client, cfg *config.TmuxConfig) *Launcher {
	return &Launcher{tmuxClient: client, config: cfg}
}

func (l *Launcher) Launch(p domain.Project) error {
	sessionName := p.SessionName(l.config.SessionPrefix)

	if !l.tmuxClient.HasSession(sessionName) {
		isSaved := false

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
