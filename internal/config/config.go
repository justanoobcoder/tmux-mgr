package config

type FolderConfig struct {
	Path     string   `mapstructure:"path" json:"path"`
	Excludes []string `mapstructure:"excludes" json:"excludes"`
}

type Config struct {
	Projects  []string        `mapstructure:"projects" json:"projects"`
	Folders   []FolderConfig  `mapstructure:"folders" json:"folders"`
	Scores    map[string]int  `mapstructure:"scores" json:"scores"`
	Tmux      TmuxConfig      `mapstructure:"tmux" json:"tmux"`
	Resurrect ResurrectConfig `mapstructure:"resurrect" json:"resurrect"`
}

type TmuxConfig struct {
	SessionPrefix  string `mapstructure:"session_prefix" json:"session_prefix"`
	AttachOnCreate bool   `mapstructure:"attach_on_create" json:"attach_on_create"`
}

type ResurrectConfig struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	SaveDir string `mapstructure:"save_dir" json:"save_dir"`
}
