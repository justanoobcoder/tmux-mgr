package tmux

import (
	"fmt"
	"strconv"
	"strings"
)

func (c *Client) HasSession(name string) bool {
	_, err := c.Run("has-session", "-t", name)
	return err == nil
}

func (c *Client) NewSession(name, startDir string) error {
	_, err := c.Run("new-session", "-d", "-s", name, "-c", startDir)
	return err
}

func (c *Client) ListSessions() ([]string, string, error) {
	out, err := c.Run("list-sessions", "-F", "#{session_name} #{session_last_attached}")
	if err != nil {
		return nil, "", fmt.Errorf("list tmux sessions: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var sessions []string
	var lastActive string
	var maxTime int64

	for _, l := range lines {
		if l == "" {
			continue
		}
		parts := strings.SplitN(l, " ", 2)
		name := parts[0]
		sessions = append(sessions, name)

		if len(parts) == 2 {
			ts, err := strconv.ParseInt(parts[1], 10, 64)
			if err == nil && ts > maxTime {
				maxTime = ts
				lastActive = name
			}
		}
	}
	return sessions, lastActive, nil
}

func (c *Client) KillSession(name string) error {
	if c.HasSession(name) {
		_, err := c.Run("kill-session", "-t", name)
		return err
	}
	return nil
}

func (c *Client) RestoreResurrect() error {
	if _, err := c.Run("start-server"); err != nil {
		return fmt.Errorf("start tmux server: %w", err)
	}

	out, err := c.Run("list-keys")
	if err != nil {
		return err
	}

	var restoreScript string
	for line := range strings.SplitSeq(out, "\n") {
		if strings.Contains(line, "resurrect") && strings.Contains(line, "restore.sh") {
			parts := strings.Split(line, "run-shell ")
			if len(parts) >= 2 {
				restoreScript = strings.TrimSpace(parts[1])
				restoreScript = strings.Trim(restoreScript, "'\"")
				break
			}
		}
	}

	if restoreScript == "" {
		return fmt.Errorf("could not find tmux-resurrect restore script in tmux bindings")
	}

	_, err = c.Run("run-shell", restoreScript)
	return err
}
