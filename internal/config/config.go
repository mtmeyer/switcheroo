package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the application configuration
type Config struct {
	Directory string        `json:"directory"`
	Preview   PreviewConfig `json:"preview"`
	Output    OutputConfig  `json:"output"`
	Theme     string        `json:"theme,omitempty"`
}

// OutputConfig controls what happens after selecting a repo/worktree
type OutputConfig struct {
	Type  string `json:"type"`
	Value string `json:"value,omitempty"`
}

// PreviewConfig represents preview pane configuration
type PreviewConfig struct {
	Enabled        bool     `json:"enabled"`
	RepoFields     []string `json:"repo_fields,omitempty"`
	WorktreeFields []string `json:"worktree_fields,omitempty"`
}

// Load loads the configuration from the standard locations
// Searches in order: ~/.config/switcheroo/config.json, ~/.switcheroo/config.json
func Load() (*Config, error) {
	configPath, err := findConfigFile()
	if err != nil {
		// No config file found, return defaults
		return defaultConfig(), nil
	}

	return loadFromFile(configPath)
}

// LoadFromPath loads configuration from a specific path
func LoadFromPath(path string) (*Config, error) {
	return loadFromFile(path)
}

// findConfigFile searches for config.json in standard locations
func findConfigFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Check standard locations in order
	locations := []string{
		filepath.Join(homeDir, ".config", "switcheroo", "config.json"),
		filepath.Join(homeDir, ".switcheroo", "config.json"),
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			return location, nil
		}
	}

	return "", errors.New("no config file found")
}

// loadFromFile loads and parses a config file
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Apply defaults for missing values
	applyDefaults(&config)
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	// Validate required fields
	if config.Directory == "" {
		return nil, errors.New("directory is required in config")
	}

	return &config, nil
}

// defaultConfig returns a config with all default values
func defaultConfig() *Config {
	config := &Config{
		Directory: "",
		Preview: PreviewConfig{
			Enabled:        true,
			RepoFields:     []string{"name", "branches", "status", "last_commit"},
			WorktreeFields: []string{"branch", "status", "last_commit", "ahead_behind"},
		},
		Output: OutputConfig{
			Type: "path",
		},
	}
	return config
}

// applyDefaults fills in default values for missing config fields
func applyDefaults(config *Config) {
	if config.Output.Type == "" {
		config.Output.Type = "path"
	}

	// Preview defaults
	if len(config.Preview.RepoFields) == 0 {
		config.Preview.RepoFields = []string{"name", "branches", "status", "last_commit"}
	}

	if len(config.Preview.WorktreeFields) == 0 {
		config.Preview.WorktreeFields = []string{"branch", "status", "last_commit", "ahead_behind"}
	}

	// no validation here to avoid panics; handled separately
}

func validateConfig(config *Config) error {
	switch config.Output.Type {
	case "path":
		return nil
	case "command":
		if strings.TrimSpace(config.Output.Value) == "" {
			return errors.New("output.value is required when output.type is 'command'")
		}
		return nil
	default:
		return errors.New("unsupported output.type: " + config.Output.Type)
	}
}
