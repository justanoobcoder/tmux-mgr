package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
	"github.com/sahilm/fuzzy"
)

type ProjectPickerModel struct {
	projects         []domain.Project
	filteredProjects []domain.Project
	cursor           int
	selected         *domain.Project
	searchInput      string
	filtering        bool
	exactSearch      bool
	confirming       bool
	ToDelete         *domain.Project
}

func NewProjectPickerModel(projects []domain.Project) *ProjectPickerModel {
	return &ProjectPickerModel{
		projects:         projects,
		filteredProjects: projects,
		filtering:        false,
		exactSearch:      false,
	}
}

func (m *ProjectPickerModel) Init() tea.Cmd {
	return nil
}

func (m *ProjectPickerModel) filterProjects() {
	if m.searchInput == "" {
		m.filteredProjects = m.projects
	} else {
		var paths []string
		for _, p := range m.projects {
			paths = append(paths, p.Path)
		}

		var filtered []domain.Project
		if !m.exactSearch {
			matches := fuzzy.Find(m.searchInput, paths)
			for _, match := range matches {
				filtered = append(filtered, m.projects[match.Index])
			}
		} else {
			for _, p := range m.projects {
				if strings.Contains(p.Path, m.searchInput) {
					filtered = append(filtered, p)
				}
			}
		}
		m.filteredProjects = filtered
	}
	m.cursor = 0
}

func (m *ProjectPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if m.confirming {
			switch msg.String() {
			case "y", "Y", "enter":
				m.ToDelete = &m.filteredProjects[m.cursor]
				return m, tea.Quit
			case "n", "N", "esc", "q":
				m.confirming = false
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.exactSearch = !m.exactSearch
			m.filterProjects()
			return m, nil
		case "enter":
			if len(m.filteredProjects) > 0 {
				m.selected = &m.filteredProjects[m.cursor]
				return m, tea.Quit
			}
		case "ctrl+d", "delete":
			if len(m.filteredProjects) > 0 {
				m.confirming = true
			}
			return m, nil
		}

		if m.filtering {
			switch msg.String() {
			case "backspace":
				if len(m.searchInput) > 0 {
					runes := []rune(m.searchInput)
					m.searchInput = string(runes[:len(runes)-1])
					m.filterProjects()
				} else {
					m.filtering = false
					m.filterProjects()
				}
			case "space":
				m.searchInput += " "
				m.filterProjects()
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.filteredProjects)-1 {
					m.cursor++
				}
			default:
				if len(msg.String()) == 1 {
					m.searchInput += msg.String()
					m.filterProjects()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "/":
			m.filtering = true
			return m, nil
		case "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filteredProjects)-1 {
				m.cursor++
			}
		case "d":
			if len(m.filteredProjects) > 0 {
				m.confirming = true
			}
		}
	}
	return m, nil
}

func (m *ProjectPickerModel) View() tea.View {
	if len(m.projects) == 0 {
		return tea.NewView("No projects found. Add one with 'tmux-mgr add <path>'\n\n(q to quit)\n")
	}

	var b strings.Builder

	if m.confirming {
		fmt.Fprintf(&b, "Are you sure you want to remove project '%s'? (y/N)\n", m.filteredProjects[m.cursor].Path)
		return tea.NewView(b.String())
	}

	b.WriteString("Select a project:\n\n")

	if m.filtering {
		mode := "Fuzzy"
		if m.exactSearch {
			mode = "Exact"
		}
		fmt.Fprintf(&b, "%s Search: %s_\n\n", mode, m.searchInput)
	}

	for i, proj := range m.filteredProjects {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		fmt.Fprintf(&b, "%s %s\n", cursor, proj.Path)
	}

	if len(m.filteredProjects) == 0 {
		b.WriteString("  No projects match search.\n")
	}

	b.WriteString("\n(j/k/up/down: navigate, enter: select, d/ctrl+d: remove, /: search, esc: toggle fuzzy search, q: quit)\n")
	return tea.NewView(b.String())
}

func (m *ProjectPickerModel) SelectedProject() *domain.Project {
	return m.selected
}

func (m *ProjectPickerModel) ProjectToDelete() *domain.Project {
	return m.ToDelete
}
