package ui

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/justanoobcoder/tmux-mgr/internal/domain"
)

type SessionPickerModel struct {
	sessions   []domain.DisplaySession
	cursor     int
	ToDelete   *domain.DisplaySession
	ToAttach   *domain.DisplaySession
	confirming bool
}

func NewSessionPickerModel(sessions []domain.DisplaySession) *SessionPickerModel {
	return &SessionPickerModel{
		sessions: sessions,
	}
}

func (m *SessionPickerModel) Init() tea.Cmd {
	return nil
}

func (m *SessionPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.confirming {
			switch msg.String() {
			case "y", "Y", "enter":
				m.ToDelete = &m.sessions[m.cursor]
				return m, tea.Quit
			case "n", "N", "esc", "q":
				m.confirming = false
				return m, nil
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "a", "enter":
			if len(m.sessions) > 0 {
				m.ToAttach = &m.sessions[m.cursor]
				return m, tea.Quit
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}
		case "d", "delete":
			if len(m.sessions) > 0 {
				m.confirming = true
			}
		}
	}
	return m, nil
}

func (m *SessionPickerModel) View() tea.View {
	if len(m.sessions) == 0 {
		return tea.NewView("No sessions found.\n\n(q to quit)\n")
	}

	var b strings.Builder

	if m.confirming {
		fmt.Fprintf(&b, "Are you sure you want to delete session '%s'? (y/N)\n", m.sessions[m.cursor].Name)
		return tea.NewView(b.String())
	}

	b.WriteString("Sessions:\n\n")

	maxNameLen := 0
	for _, sess := range m.sessions {
		if len(sess.Name) > maxNameLen {
			maxNameLen = len(sess.Name)
		}
	}
	maxNameLen += 2

	for i, sess := range m.sessions {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		namePadded := fmt.Sprintf("%-*s", maxNameLen, sess.Name)

		var tags []string
		if sess.IsActive {
			if sess.IsLastActive {
				tags = append(tags, "Active*")
			} else {
				tags = append(tags, "Active")
			}
		}
		if sess.IsSaved {
			if sess.IsLastSaved {
				tags = append(tags, "Saved*")
			} else {
				tags = append(tags, "Saved")
			}
		}

		tagStr := ""
		if len(tags) > 0 {
			tagStr = fmt.Sprintf("[%s]", strings.Join(tags, " | "))
		}

		fmt.Fprintf(&b, "%s %s %s\n", cursor, namePadded, tagStr)
	}

	b.WriteString("\n(j/k/up/down: navigate, a/enter: attach, d: delete, q: quit)\n")
	return tea.NewView(b.String())
}
