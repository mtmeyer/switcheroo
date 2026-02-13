package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme contains all the styling for the UI
// This structure allows for easy theming in the future
type Theme struct {
	// Colors
	AccentColor      lipgloss.Color
	MutedColor       lipgloss.Color
	DisabledColor    lipgloss.Color
	SuccessColor     lipgloss.Color
	WarningColor     lipgloss.Color
	ErrorColor       lipgloss.Color
	BackgroundColor  lipgloss.Color
	ForegroundColor  lipgloss.Color
	BorderColor      lipgloss.Color
	CurrentLineColor lipgloss.Color
	InputBgColor     lipgloss.Color
	ListBgColor      lipgloss.Color
	PreviewBgColor   lipgloss.Color
	SelectedBgColor  lipgloss.Color

	// Styles
	TitleStyle       lipgloss.Style
	InputStyle       lipgloss.Style
	InputPromptStyle lipgloss.Style
	ItemStyle        lipgloss.Style
	SelectedStyle    lipgloss.Style
	DisabledStyle    lipgloss.Style
	LabelStyle       lipgloss.Style
	ValueStyle       lipgloss.Style
	HelpStyle        lipgloss.Style
	BorderStyle      lipgloss.Style
	PreviewStyle     lipgloss.Style

	// Nerd Font Icons
	Icons struct {
		Folder       string
		Git          string
		Branch       string
		Worktree     string
		Check        string
		Warning      string
		Modified     string
		Search       string
		ChevronRight string
	}
}

// ThemeDefinition describes the color palette names used to construct a Theme.
// Fields are optional; empty values defer to the user's terminal colors.
type ThemeDefinition struct {
	Accent      string `json:"accent"`
	Muted       string `json:"muted"`
	Disabled    string `json:"disabled"`
	Success     string `json:"success"`
	Warning     string `json:"warning"`
	Error       string `json:"error"`
	Background  string `json:"background"`
	Foreground  string `json:"foreground"`
	Border      string `json:"border"`
	CurrentLine string `json:"current_line"`
	InputBg     string `json:"input_bg"`
	ListBg      string `json:"list_bg"`
	PreviewBg   string `json:"preview_bg"`
	SelectedBg  string `json:"selected_bg"`
}

// TerminalDefaultTheme returns a theme that defers to terminal colors.
func TerminalDefaultTheme() Theme {
	return buildTheme(ThemeDefinition{})
}

// buildTheme converts a ThemeDefinition into a full Theme with styles.
func buildTheme(def ThemeDefinition) Theme {
	t := Theme{
		AccentColor:      makeColor(def.Accent),
		MutedColor:       makeColor(def.Muted),
		DisabledColor:    makeColor(def.Disabled),
		SuccessColor:     makeColor(def.Success),
		WarningColor:     makeColor(def.Warning),
		ErrorColor:       makeColor(def.Error),
		BackgroundColor:  makeColor(def.Background),
		ForegroundColor:  makeColor(def.Foreground),
		BorderColor:      makeColor(def.Border),
		CurrentLineColor: makeColor(def.CurrentLine),
		InputBgColor:     makeColor(def.InputBg),
		ListBgColor:      makeColor(def.ListBg),
		PreviewBgColor:   makeColor(def.PreviewBg),
		SelectedBgColor:  makeColor(def.SelectedBg),
	}

	setDefaultIcons(&t)
	buildStyles(&t)
	return t
}

func buildStyles(t *Theme) {
	t.TitleStyle = lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)
	if t.AccentColor != "" {
		t.TitleStyle = t.TitleStyle.Foreground(t.AccentColor)
	}

	t.InputStyle = lipgloss.NewStyle().
		Padding(0, 1).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder())
	if t.ForegroundColor != "" {
		t.InputStyle = t.InputStyle.Foreground(t.ForegroundColor)
	}
	if t.InputBgColor != "" {
		t.InputStyle = t.InputStyle.Background(t.InputBgColor)
	}
	if t.BorderColor != "" {
		t.InputStyle = t.InputStyle.BorderForeground(t.BorderColor)
	}

	t.InputPromptStyle = lipgloss.NewStyle().
		Bold(true)
	if t.AccentColor != "" {
		t.InputPromptStyle = t.InputPromptStyle.Foreground(t.AccentColor)
	}

	t.ItemStyle = lipgloss.NewStyle()
	if t.ForegroundColor != "" {
		t.ItemStyle = t.ItemStyle.Foreground(t.ForegroundColor)
	}

	selectedFg := firstColor(t.CurrentLineColor, t.AccentColor, t.ForegroundColor)
	t.SelectedStyle = lipgloss.NewStyle().Bold(true)
	if selectedFg != "" {
		t.SelectedStyle = t.SelectedStyle.Foreground(selectedFg)
	}
	if t.SelectedBgColor != "" {
		t.SelectedStyle = t.SelectedStyle.Background(t.SelectedBgColor)
	}

	t.DisabledStyle = lipgloss.NewStyle().Faint(true)
	if t.DisabledColor != "" {
		t.DisabledStyle = t.DisabledStyle.Foreground(t.DisabledColor)
	}

	t.LabelStyle = lipgloss.NewStyle().Bold(true)
	if t.MutedColor != "" {
		t.LabelStyle = t.LabelStyle.Foreground(t.MutedColor)
	}

	t.ValueStyle = lipgloss.NewStyle()
	if t.ForegroundColor != "" {
		t.ValueStyle = t.ValueStyle.Foreground(t.ForegroundColor)
	}

	t.HelpStyle = lipgloss.NewStyle().MarginTop(1)
	if t.MutedColor != "" {
		t.HelpStyle = t.HelpStyle.Foreground(t.MutedColor)
	}

	t.BorderStyle = lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder())
	if t.BorderColor != "" {
		t.BorderStyle = t.BorderStyle.BorderForeground(t.BorderColor)
	}
	if t.BackgroundColor != "" {
		t.BorderStyle = t.BorderStyle.Background(t.BackgroundColor)
	}

	t.PreviewStyle = lipgloss.NewStyle().
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder())
	if t.PreviewBgColor != "" {
		t.PreviewStyle = t.PreviewStyle.Background(t.PreviewBgColor)
	}
	if t.BorderColor != "" {
		t.PreviewStyle = t.PreviewStyle.BorderForeground(t.BorderColor)
	}
}

func setDefaultIcons(t *Theme) {
	t.Icons.Folder = "\uf07c"
	t.Icons.Git = "\ue725"
	t.Icons.Branch = "\ue725"
	t.Icons.Worktree = "\uf115"
	t.Icons.Check = "\uf00c"
	t.Icons.Warning = "\uf071"
	t.Icons.Modified = "\uf040"
	t.Icons.Search = "\uf002"
	t.Icons.ChevronRight = "\uf054"
}

func makeColor(value string) lipgloss.Color {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return lipgloss.Color("")
	}
	return lipgloss.Color(trimmed)
}

func firstColor(colors ...lipgloss.Color) lipgloss.Color {
	for _, c := range colors {
		if c != "" {
			return c
		}
	}
	return lipgloss.Color("")
}
