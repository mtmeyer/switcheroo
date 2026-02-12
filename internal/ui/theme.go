package ui

import "github.com/charmbracelet/lipgloss"

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

// DefaultTheme returns the default theme
func DefaultTheme() Theme {
	t := Theme{
		// Colors
		AccentColor:      lipgloss.Color("12"),  // Cyan
		MutedColor:       lipgloss.Color("8"),   // Gray
		DisabledColor:    lipgloss.Color("240"), // Dark gray
		SuccessColor:     lipgloss.Color("10"),  // Green
		WarningColor:     lipgloss.Color("11"),  // Yellow
		ErrorColor:       lipgloss.Color("9"),   // Red
		BackgroundColor:  lipgloss.Color("0"),   // Black
		ForegroundColor:  lipgloss.Color("15"),  // White
		BorderColor:      lipgloss.Color("12"),  // Cyan for left borders
		CurrentLineColor: lipgloss.Color("14"),  // Bright cyan
		InputBgColor:     lipgloss.Color("234"), // Subtle dark background
		ListBgColor:      lipgloss.Color("235"), // Slightly lighter
		PreviewBgColor:   lipgloss.Color("236"), // Even lighter for depth
		SelectedBgColor:  lipgloss.Color("237"), // Highlighted item bg
	}

	// Nerd Font Icons (using unicode escapes)
	t.Icons.Folder = "\uf07c"       // nf-fa-folder
	t.Icons.Git = "\ue725"          // nf-dev-git_branch
	t.Icons.Branch = "\ue725"       // nf-oct-git_branch
	t.Icons.Worktree = "\uf115"     // nf-fa-folder_open
	t.Icons.Check = "\uf00c"        // nf-fa-check
	t.Icons.Warning = "\uf071"      // nf-fa-exclamation_triangle
	t.Icons.Modified = "\uf040"     // nf-fa-pencil
	t.Icons.Search = "\uf002"       // nf-fa-search
	t.Icons.ChevronRight = "\uf054" // nf-fa-chevron_right

	// Build styles from colors
	t.TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.AccentColor).
		MarginBottom(1)

	t.InputStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor).
		Background(t.InputBgColor).
		Padding(0, 1).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(t.BorderColor)

	t.InputPromptStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.AccentColor)

	t.ItemStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor)

	t.SelectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.CurrentLineColor).
		Background(t.SelectedBgColor)

	t.DisabledStyle = lipgloss.NewStyle().
		Foreground(t.DisabledColor).
		Faint(true)

	t.LabelStyle = lipgloss.NewStyle().
		Foreground(t.MutedColor).
		Bold(true)

	t.ValueStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor)

	t.HelpStyle = lipgloss.NewStyle().
		Foreground(t.MutedColor).
		MarginTop(1)

	t.BorderStyle = lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(t.BorderColor)

	t.PreviewStyle = lipgloss.NewStyle().
		Background(t.PreviewBgColor).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(t.BorderColor)

	return t
}
