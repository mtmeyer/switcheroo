package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"switcheroo/internal/config"
	"switcheroo/internal/git"
	"switcheroo/internal/ui"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// Load configuration
	var cfg *config.Config
	var err error

	if *configPath != "" {
		cfg, err = config.LoadFromPath(*configPath)
	} else {
		cfg, err = config.Load()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		fmt.Fprintln(os.Stderr, "\nCreate a config file at one of these locations:")
		fmt.Fprintln(os.Stderr, "  ~/.config/switcheroo/config.json")
		fmt.Fprintln(os.Stderr, "  ~/.switcheroo/config.json")
		fmt.Fprintln(os.Stderr, "\nExample config:")
		fmt.Fprintln(os.Stderr, `{
  "directory": "/path/to/repos",
  "output": {
    "type": "path"
  },
  "theme": "dracula",
  "preview": {
    "enabled": true
  }
}`)
		os.Exit(1)
	}

	// Discover repositories from configured directory
	repos, err := git.DiscoverRepositories(cfg.Directory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Fprintf(os.Stderr, "No git repositories found in: %s\n", cfg.Directory)
		os.Exit(1)
	}

	// Convert to UI display format
	repoDisplays := ui.FromGitRepositories(repos)

	// Load theme based on config
	var theme ui.Theme
	if cfg.Theme == "" || cfg.Theme == "default" {
		// Use terminal default theme
		theme = ui.TerminalDefaultTheme()
	} else {
		// Try to load theme from file
		themePath, err := ui.FindThemeFile(cfg.Theme)
		if err != nil {
			// Theme not found, fall back to terminal default
			fmt.Fprintf(os.Stderr, "Warning: Theme '%s' not found, using default\n", cfg.Theme)
			theme = ui.TerminalDefaultTheme()
		} else {
			theme, err = ui.LoadThemeFromFile(themePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Error loading theme '%s': %v, using default\n", cfg.Theme, err)
				theme = ui.TerminalDefaultTheme()
			}
		}
	}

	// Initialize Bubble Tea program with real data and theme
	p := tea.NewProgram(ui.NewModel(repoDisplays, theme), tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	model, ok := finalModel.(ui.Model)
	if !ok {
		os.Exit(0)
	}

	selectedPath := model.SelectedPath()
	if selectedPath == "" {
		os.Exit(0)
	}

	switch cfg.Output.Type {
	case "command":
		if err := runConfiguredCommand(cfg.Output.Value, selectedPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
			os.Exit(1)
		}
	case "path":
		fmt.Println(selectedPath)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported output type: %s\n", cfg.Output.Type)
		os.Exit(1)
	}
}

var execCommandRunner = defaultExecCommandRunner

func defaultExecCommandRunner(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runConfiguredCommand(template, path string) error {
	replaced := strings.ReplaceAll(template, "{{path}}", path)
	parts := strings.Fields(replaced)
	if len(parts) == 0 {
		return fmt.Errorf("invalid command template: %s", template)
	}
	return execCommandRunner(parts[0], parts[1:])
}
