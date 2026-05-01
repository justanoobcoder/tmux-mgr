package service

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
)

type Manager struct {
	cfg *config.Config
}

func NewManager(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) AddProject(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	cleanPath := filepath.Clean(absPath)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("stat path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", cleanPath)
	}

	if slices.Contains(m.cfg.Projects, cleanPath) {
		return fmt.Errorf("project already exists in config")
	}

	m.cfg.Projects = append(m.cfg.Projects, cleanPath)

	if err := config.Save(m.cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

func (m *Manager) RemoveConfigPath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	cleanPath := filepath.Clean(absPath)

	removed := false

	var newProjects []string
	for _, p := range m.cfg.Projects {
		if p == cleanPath {
			removed = true
		} else {
			newProjects = append(newProjects, p)
		}
	}
	m.cfg.Projects = newProjects

	if !removed {
		return fmt.Errorf("path not found in config: %s", cleanPath)
	}

	if err := config.Save(m.cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}
