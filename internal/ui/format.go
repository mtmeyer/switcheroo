package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func formatStatusInline(theme Theme, status string) string {
	return "(" + renderStatus(theme, status) + ")"
}

func formatStatusValue(theme Theme, status string) string {
	return renderStatus(theme, status)
}

func renderStatus(theme Theme, status string) string {
	if status == "" {
		status = "unknown"
	}
	style := lipgloss.NewStyle()
	switch status {
	case "clean", "ahead":
		if theme.SuccessColor != "" {
			style = style.Foreground(theme.SuccessColor)
		}
	case "behind", "modified":
		if theme.WarningColor != "" {
			style = style.Foreground(theme.WarningColor)
		}
	case "diverged", "untracked":
		if theme.ErrorColor != "" {
			style = style.Foreground(theme.ErrorColor)
		}
	default:
		if theme.MutedColor != "" {
			style = style.Foreground(theme.MutedColor)
		}
	}
	return style.Render(status)
}

func formatDiffText(theme Theme, added, removed int) string {
	addedText := fmt.Sprintf("+%d", added)
	removedText := fmt.Sprintf("-%d", removed)
	if theme.SuccessColor != "" {
		addedText = lipgloss.NewStyle().Foreground(theme.SuccessColor).Render(addedText)
	}
	if theme.ErrorColor != "" {
		removedText = lipgloss.NewStyle().Foreground(theme.ErrorColor).Render(removedText)
	}
	return addedText + " " + removedText
}

// truncateWithEllipsis truncates a string to maxWidth and adds "..." if truncated
func truncateWithEllipsis(s string, maxWidth int) string {
	if maxWidth <= 3 {
		return "..."
	}
	if len(s) <= maxWidth {
		return s
	}
	return s[:maxWidth-3] + "..."
}
