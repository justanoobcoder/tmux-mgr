package config

type Config struct {
	Projects  []string        `mapstructure:"projects" json:"projects"`
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
