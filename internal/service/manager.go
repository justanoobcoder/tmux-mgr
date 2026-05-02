package service

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/justanoobcoder/tmux-mgr/internal/config"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
	"github.com/justanoobcoder/tmux-mgr/internal/resurrect"
)

type Manager struct {
	cfg   *config.Config
	store *resurrect.Store
}

func NewManager(cfg *config.Config, store *resurrect.Store) *Manager {
	return &Manager{cfg: cfg, store: store}
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

	for i, f := range m.cfg.Folders {
		dir := filepath.Dir(cleanPath)
		if dir == f.Path {
			base := filepath.Base(cleanPath)
			var newExcludes []string
			unexcluded := false
			for _, ex := range f.Excludes {
				if ex == base {
					unexcluded = true
				} else {
					newExcludes = append(newExcludes, ex)
				}
			}
			if unexcluded {
				m.cfg.Folders[i].Excludes = newExcludes
				if err := config.Save(m.cfg); err != nil {
					return fmt.Errorf("save config: %w", err)
				}
				return nil
			}
		}
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

func (m *Manager) AddFolder(path string, excludes []string) error {
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

	for _, f := range m.cfg.Folders {
		if f.Path == cleanPath {
			return fmt.Errorf("folder already exists in config")
		}
	}

	m.cfg.Folders = append(m.cfg.Folders, config.FolderConfig{
		Path:     cleanPath,
		Excludes: excludes,
	})

	if err := config.Save(m.cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	return nil
}

func (m *Manager) GetProjects() []domain.Project {
	projectMap := make(map[string]bool)
	var projects []domain.Project

	for _, p := range m.cfg.Projects {
		if !projectMap[p] {
			projectMap[p] = true
			projects = append(projects, domain.NewProject(p))
		}
	}

	for _, f := range m.cfg.Folders {
		entries, err := os.ReadDir(f.Path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			excluded := slices.Contains(f.Excludes, entry.Name())

			if !excluded {
				projPath := filepath.Join(f.Path, entry.Name())
				if !projectMap[projPath] {
					projectMap[projPath] = true
					projects = append(projects, domain.NewProject(projPath))
				}
			}
		}
	}

	sort.Slice(projects, func(i, j int) bool {
		scoreI := m.cfg.Scores[projects[i].Path]
		scoreJ := m.cfg.Scores[projects[j].Path]
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return projects[i].Path < projects[j].Path
	})

	return projects
}

func (m *Manager) DeleteSession(name string) error {
	if m.store != nil {
		if err := m.store.DeleteSession(name); err != nil {
			return fmt.Errorf("delete session from resurrect: %w", err)
		}
	}
	return nil
}

func (m *Manager) RemoveProject(path string) error {
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

	for i, f := range m.cfg.Folders {
		dir := filepath.Dir(cleanPath)
		if dir == f.Path {
			base := filepath.Base(cleanPath)
			alreadyExcluded := slices.Contains(f.Excludes, base)
			if !alreadyExcluded {
				m.cfg.Folders[i].Excludes = append(m.cfg.Folders[i].Excludes, base)
				removed = true
			}
		}
	}

	if !removed {
		return fmt.Errorf("project not found")
	}

	if m.cfg.Scores != nil {
		delete(m.cfg.Scores, cleanPath)
	}

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

	var newFolders []config.FolderConfig
	for _, f := range m.cfg.Folders {
		if f.Path == cleanPath {
			removed = true
		} else {
			newFolders = append(newFolders, f)
		}
	}
	m.cfg.Folders = newFolders

	if !removed {
		return fmt.Errorf("path not found in config: %s", cleanPath)
	}

	if m.cfg.Scores != nil {
		delete(m.cfg.Scores, cleanPath)
		prefix := cleanPath + string(filepath.Separator)
		for k := range m.cfg.Scores {
			if strings.HasPrefix(k, prefix) {
				delete(m.cfg.Scores, k)
			}
		}
	}

	if err := config.Save(m.cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	return nil
}

func (m *Manager) TrackProjectSelection(path string) error {
	if m.cfg.Scores == nil {
		m.cfg.Scores = make(map[string]int)
	}
	m.cfg.Scores[path]++
	return config.Save(m.cfg)
}
