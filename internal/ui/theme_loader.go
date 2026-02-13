package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadThemeFromFile parses a JSON theme definition and returns a Theme.
func LoadThemeFromFile(path string) (Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Theme{}, err
	}

	var def ThemeDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return Theme{}, err
	}

	return buildTheme(def), nil
}

// FindThemeFile attempts to resolve a theme name to an on-disk JSON file.
// It searches explicit paths first, followed by common theme directories.
func FindThemeFile(name string) (string, error) {
	if strings.TrimSpace(name) == "" {
		return "", errors.New("theme name is empty")
	}

	if looksLikePath(name) {
		candidate := expandPath(name)
		if filepath.Ext(candidate) == "" {
			candidate += ".json"
		}
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		return "", fmt.Errorf("theme path %q not found", candidate)
	}

	filename := ensureJSONExtension(name)
	candidates := []string{}

	// Current working directory /themes
	cwdTheme := filepath.Join("themes", filename)
	candidates = append(candidates, cwdTheme)

	// Executable directory /themes
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		candidates = append(candidates, filepath.Join(execDir, "themes", filename))
	}

	// XDG config directory or ~/.config fallback
	if configHome := os.Getenv("XDG_CONFIG_HOME"); configHome != "" {
		candidates = append(candidates, filepath.Join(configHome, "switcheroo", "themes", filename))
	} else if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".config", "switcheroo", "themes", filename))
	}

	// ~/.switcheroo fallback
	if home, err := os.UserHomeDir(); err == nil {
		candidates = append(candidates, filepath.Join(home, ".switcheroo", "themes", filename))
	}

	for _, candidate := range dedupeStrings(candidates) {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("theme %q not found", name)
}

func looksLikePath(input string) bool {
	if filepath.IsAbs(input) {
		return true
	}
	if strings.HasPrefix(input, "~") || strings.HasPrefix(input, ".") {
		return true
	}
	return strings.ContainsRune(input, '/') || strings.ContainsRune(input, '\\')
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			rest := strings.TrimPrefix(path, "~")
			rest = strings.TrimLeft(rest, string(filepath.Separator)+"/\\")
			if rest == "" {
				return home
			}
			return filepath.Join(home, rest)
		}
	}
	return path
}

func ensureJSONExtension(name string) string {
	if strings.HasSuffix(name, ".json") {
		return name
	}
	return name + ".json"
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, v := range values {
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
