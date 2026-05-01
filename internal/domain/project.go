package domain

import (
	"path/filepath"
	"strings"
)

type Project struct {
	Path string
}

func NewProject(path string) Project {
	return Project{Path: filepath.Clean(path)}
}

func (p Project) SessionName(prefix string) string {
	base := filepath.Base(p.Path)

	// tmux session names cannot contain . or :
	name := strings.ReplaceAll(base, ".", "_")
	name = strings.ReplaceAll(name, ":", "_")

	if prefix != "" {
		return prefix + name
	}
	return name
}
