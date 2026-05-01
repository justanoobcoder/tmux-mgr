package config

type Config struct {
	Projects []string   `mapstructure:"projects" json:"projects"`
	Tmux     TmuxConfig `mapstructure:"tmux" json:"tmux"`
}

type TmuxConfig struct {
	SessionPrefix  string `mapstructure:"session_prefix" json:"session_prefix"`
	AttachOnCreate bool   `mapstructure:"attach_on_create" json:"attach_on_create"`
}
