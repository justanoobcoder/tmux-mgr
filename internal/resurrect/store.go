package resurrect

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *Store) DeleteSession(name string) error {
	files, err := filepath.Glob(filepath.Join(s.SaveDir, "tmux_resurrect_*.txt"))
	if err != nil {
		return fmt.Errorf("glob resurrect files: %w", err)
	}

	for _, file := range files {
		err := s.deleteSessionFromFile(file, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) deleteSessionFromFile(path string, sessionName string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var newLines []string
	hasSession := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")

		if len(parts) >= 2 && (parts[0] == "pane" || parts[0] == "window" || parts[0] == "state") && parts[1] == sessionName {
			hasSession = true
			continue
		}

		newLines = append(newLines, line)
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		return err
	}

	if !hasSession {
		return nil
	}

	hasAnyPaneOrWindow := false
	for _, l := range newLines {
		parts := strings.Split(l, "\t")
		if len(parts) >= 1 && (parts[0] == "pane" || parts[0] == "window") {
			hasAnyPaneOrWindow = true
			break
		}
	}

	if !hasAnyPaneOrWindow {
		return os.Remove(path)
	}

	content := strings.Join(newLines, "\n") + "\n"
	return os.WriteFile(path, []byte(content), 0644)
}
