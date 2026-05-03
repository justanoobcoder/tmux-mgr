package resurrect

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/justanoobcoder/tmux-mgr/internal/domain"
)

type Store struct {
	SaveDir string
}

func NewStore(saveDir string) *Store {
	return &Store{SaveDir: saveDir}
}

func (s *Store) ListSessions() ([]domain.SavedSession, string, error) {
	lastFile := filepath.Join(s.SaveDir, "last")

	file, err := os.Open(lastFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", nil
		}
		return nil, "", err
	}
	defer file.Close()

	sessionMap := make(map[string]bool)
	var sessions []domain.SavedSession
	var lastSaved string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			switch parts[0] {
			case "state":
				lastSaved = parts[1]
			case "pane", "window":
				sessName := parts[1]
				if sessName != "" && !sessionMap[sessName] {
					sessionMap[sessName] = true
					sessions = append(sessions, domain.SavedSession{Name: sessName})
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, "", fmt.Errorf("scan last save file: %w", err)
	}

	return sessions, lastSaved, nil
}
