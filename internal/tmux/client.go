package tmux

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func cleanEnv() []string {
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "DIRENV_") {
			env = append(env, e)
		}
	}
	return env
}

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Run(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	cmd.Env = cleanEnv()
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("tmux error: %v, stderr: %s", err, stderr.String())
	}
	return out.String(), nil
}

func (c *Client) Attach(sessionName string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf("find tmux: %w", err)
	}

	args := []string{"tmux", "attach-session", "-t", sessionName}

	return syscall.Exec(tmuxPath, args, cleanEnv())
}
