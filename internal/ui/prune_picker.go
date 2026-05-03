package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type PrunePickerModel struct {
	paths     []string
	selected  map[int]bool
	cursor    int
	confirmed bool
	quitting  bool
}

func NewPrunePickerModel(paths []string) *PrunePickerModel {
	return &PrunePickerModel{
		paths:    paths,
		selected: make(map[int]bool),
	}
}

func (m *PrunePickerModel) Init() tea.Cmd {
	return nil
}

func (m *PrunePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.paths)-1 {
				m.cursor++
			}
		case "space", "x":
			m.selected[m.cursor] = !m.selected[m.cursor]
		case "a":
			allSelected := true
			for i := range m.paths {
				if !m.selected[i] {
					allSelected = false
					break
				}
			}
			for i := range m.paths {
				m.selected[i] = !allSelected
			}
		case "enter":
			if len(m.GetSelectedPaths()) > 0 {
				m.confirmed = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *PrunePickerModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var b strings.Builder
	b.WriteString("The following paths no longer exist. Select paths to remove from config:\n\n")

	for i, path := range m.paths {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if m.selected[i] {
			checked = "x"
		}

		fmt.Fprintf(&b, "%s [%s] %s\n", cursor, checked, path)
	}

	b.WriteString("\n(j/k: navigate, space/x: toggle, a: toggle all, enter: confirm, q: cancel)\n")
	return tea.NewView(b.String())
}

func (m *PrunePickerModel) GetSelectedPaths() []string {
	var selected []string
	for i, path := range m.paths {
		if m.selected[i] {
			selected = append(selected, path)
		}
	}
	return selected
}

func (m *PrunePickerModel) IsConfirmed() bool {
	return m.confirmed
}
