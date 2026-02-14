package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	maxContainerWidth = 100
	baseListWidth     = 35
	basePreviewWidth  = 50
	gapWidth          = 2
	paneBorderWidth   = 1
	containerPadding  = 2
	containerBorder   = 2
	paneHeight        = 20
)

type layoutMetrics struct {
	containerWidth int
	contentWidth   int
	listWidth      int
	previewWidth   int
}

func calculateLayout(termWidth int) layoutMetrics {
	if termWidth <= 0 {
		naturalContent := baseListWidth + basePreviewWidth + gapWidth + paneBorderWidth*2
		return layoutMetrics{
			containerWidth: naturalContent + containerPadding*2 + containerBorder,
			contentWidth:   naturalContent,
			listWidth:      baseListWidth,
			previewWidth:   basePreviewWidth,
		}
	}

	containerWidth := termWidth
	if containerWidth > maxContainerWidth {
		containerWidth = maxContainerWidth
	}

	innerWidth := containerWidth - (containerPadding*2 + containerBorder)
	if innerWidth < 4 {
		innerWidth = 4
	}

	usable := innerWidth - gapWidth - paneBorderWidth*2
	if usable < 2 {
		usable = 2
	}

	totalBase := baseListWidth + basePreviewWidth
	listWidth := usable * baseListWidth / totalBase
	previewWidth := usable - listWidth

	if listWidth < 1 {
		listWidth = 1
	}
	if previewWidth < 1 {
		previewWidth = 1
	}

	contentWidth := innerWidth

	return layoutMetrics{
		containerWidth: containerWidth,
		contentWidth:   contentWidth,
		listWidth:      listWidth,
		previewWidth:   previewWidth,
	}
}

// renderSideBySide renders the list and preview panes side by side
func renderSideBySide(listContent, previewContent string, layout layoutMetrics, theme Theme) string {
	listPaneWidth := layout.listWidth
	previewPaneWidth := layout.previewWidth

	if listPaneWidth < 1 {
		listPaneWidth = 1
	}
	if previewPaneWidth < 1 {
		previewPaneWidth = 1
	}

	// Render left pane (list)
	leftStyle := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(theme.BorderColor)

	if theme.ListBgColor != lipgloss.Color("") {
		leftStyle = leftStyle.Background(theme.ListBgColor)
	}
	leftPane := leftStyle.Render(listContent)

	// Render right pane (preview)
	rightStyle := lipgloss.NewStyle().
		Width(previewPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(theme.BorderColor)

	if theme.PreviewBgColor != lipgloss.Color("") {
		rightStyle = rightStyle.Background(theme.PreviewBgColor)
	}
	rightPane := rightStyle.Render(previewContent)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, strings.Repeat(" ", gapWidth), rightPane)
}

// wrapInContainer wraps the inner content in the main container with border
func wrapInContainer(inner string, layout layoutMetrics, hasSize bool, theme Theme, width, height int) string {
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(theme.BorderColor).
		Padding(1, containerPadding)

	if hasSize {
		containerStyle = containerStyle.Width(layout.containerWidth)
	}

	if theme.BackgroundColor != lipgloss.Color("") {
		containerStyle = containerStyle.Background(theme.BackgroundColor)
	}

	contentBox := containerStyle.Render(inner)

	if !hasSize {
		return contentBox
	}

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		contentBox,
	)
}
