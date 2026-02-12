package config

// Config represents the application configuration
type Config struct {
	Directory string        `json:"directory"`
	Command   string        `json:"command"`
	Output    string        `json:"output"`
	Preview   PreviewConfig `json:"preview"`
}

// PreviewConfig represents preview pane configuration
type PreviewConfig struct {
	Enabled        bool     `json:"enabled"`
	RepoFields     []string `json:"repo_fields"`
	WorktreeFields []string `json:"worktree_fields"`
}

// TODO: Implement configuration loading
