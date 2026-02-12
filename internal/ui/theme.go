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
		BorderColor:      lipgloss.Color("8"),   // Gray
		CurrentLineColor: lipgloss.Color("14"),  // Bright cyan
	}

	// Build styles from colors
	t.TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.AccentColor)

	t.InputStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor)

	t.InputPromptStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.AccentColor)

	t.ItemStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor)

	t.SelectedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.CurrentLineColor)

	t.DisabledStyle = lipgloss.NewStyle().
		Foreground(t.DisabledColor).
		Faint(true)

	t.LabelStyle = lipgloss.NewStyle().
		Foreground(t.MutedColor)

	t.ValueStyle = lipgloss.NewStyle().
		Foreground(t.ForegroundColor)

	t.HelpStyle = lipgloss.NewStyle().
		Foreground(t.MutedColor)

	t.BorderStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderColor)

	t.PreviewStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderColor)

	return t
}
