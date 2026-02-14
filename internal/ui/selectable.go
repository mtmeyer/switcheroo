package ui

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

const maxPreviewListItems = 8

// Selectable represents an item that can be displayed and selected in a list
type Selectable interface {
	GetID() string
	GetName() string
	GetPath() string
	GetBranch() string
	HasChildren() bool

	// Preview rendering
	RenderHeader(theme Theme) string
	RenderChildren(theme Theme, settings PreviewSettings) []string
	RenderMetadata(theme Theme, settings PreviewSettings) []string

	// Selection handling - returns the path to output, or empty string if has children
	OnSelect() string
}

// ListSelectModel handles generic list selection for any Selectable items
type ListSelectModel struct {
	Title         string
	items         []Selectable
	filteredItems []Selectable
	cursor        int
	scrollOffset  int
	searchQuery   string
	Theme         Theme
	settings      PreviewSettings
	Width         int
	Height        int

	onSelect func(Selectable) tea.Cmd
}

// NewListSelectModel creates a new list selection model
func NewListSelectModel(title string, items []Selectable, theme Theme, settings PreviewSettings, onSelect func(Selectable) tea.Cmd) ListSelectModel {
	return ListSelectModel{
		Title:         title,
		items:         items,
		filteredItems: items,
		Theme:         theme,
		settings:      settings,
		onSelect:      onSelect,
	}
}

// GetCursor returns the current cursor position
func (m ListSelectModel) GetCursor() int {
	return m.cursor
}

// GetSelected returns the currently selected item
func (m ListSelectModel) GetSelected() Selectable {
	if len(m.filteredItems) == 0 || m.cursor >= len(m.filteredItems) {
		return nil
	}
	return m.filteredItems[m.cursor]
}

// GetSearchQuery returns the current search query
func (m ListSelectModel) GetSearchQuery() string {
	return m.searchQuery
}

// UpdateSearchQuery updates the search query and filters items
func (m *ListSelectModel) UpdateSearchQuery(query string) {
	m.searchQuery = query
	m.filterItems()
}

// MoveCursor moves the cursor by delta and adjusts scroll offset
func (m *ListSelectModel) MoveCursor(delta int) {
	newCursor := m.cursor + delta
	if newCursor < 0 {
		newCursor = 0
	}
	if newCursor >= len(m.filteredItems) {
		newCursor = len(m.filteredItems) - 1
	}
	if newCursor < 0 {
		newCursor = 0
	}
	m.cursor = newCursor
	m.adjustScrollOffset()
}

// ClearSearch clears the search query and resets filtering
func (m *ListSelectModel) ClearSearch() {
	m.searchQuery = ""
	m.filterItems()
}

// SetDimensions updates the width and height
func (m *ListSelectModel) SetDimensions(width, height int) {
	m.Width = width
	m.Height = height
}

// SelectCurrent triggers selection of the current item
func (m ListSelectModel) SelectCurrent() tea.Cmd {
	if selected := m.GetSelected(); selected != nil {
		return m.onSelect(selected)
	}
	return nil
}

func (m *ListSelectModel) filterItems() {
	if m.searchQuery == "" {
		m.filteredItems = m.items
	} else {
		// Build slice of names for fuzzy search
		names := make([]string, len(m.items))
		for i, item := range m.items {
			names[i] = item.GetName()
		}

		// Perform fuzzy search
		matches := fuzzy.Find(m.searchQuery, names)

		// Sort by score (highest first)
		sort.Slice(matches, func(i, j int) bool {
			return matches[i].Score > matches[j].Score
		})

		// Rebuild filtered items in score order
		m.filteredItems = make([]Selectable, 0, len(matches))
		for _, match := range matches {
			m.filteredItems = append(m.filteredItems, m.items[match.Index])
		}
	}

	if m.cursor >= len(m.filteredItems) {
		m.cursor = 0
		m.scrollOffset = 0
	}
}

func (m *ListSelectModel) adjustScrollOffset() {
	visibleHeight := paneHeight - 2

	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}

	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}
}

// RenderList renders the scrollable list of items
func (m ListSelectModel) RenderList(layout layoutMetrics) string {
	if len(m.filteredItems) == 0 {
		return m.Theme.LabelStyle.Render("No items found")
	}

	visibleHeight := paneHeight - 2
	var lines []string

	startIdx := m.scrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.filteredItems) {
		endIdx = len(m.filteredItems)
	}

	for i := startIdx; i < endIdx; i++ {
		item := m.filteredItems[i]
		lines = append(lines, m.renderItemLine(item, i == m.cursor, layout))
	}

	return strings.Join(lines, "\n")
}

func (m ListSelectModel) renderItemLine(item Selectable, isSelected bool, layout layoutMetrics) string {
	var icon string
	if item.HasChildren() {
		icon = m.Theme.Icons.Worktree
	} else {
		icon = m.Theme.Icons.Git
	}

	// Calculate max width dynamically based on container
	// Account for: 4 padding + 1 cursor + 1 icon + 3 spaces = 9 characters
	maxNameWidth := layout.listWidth - 9
	if maxNameWidth < 10 {
		maxNameWidth = 10
	}

	name := truncateWithEllipsis(item.GetName(), maxNameWidth)

	if isSelected {
		cursor := m.Theme.Icons.ChevronRight
		content := cursor + " " + icon + " " + name
		return m.Theme.SelectedStyle.Render(content)
	} else {
		content := "  " + icon + " " + name
		return m.Theme.ItemStyle.Render(content)
	}
}

// RenderPreview renders the preview pane for the selected item
func (m ListSelectModel) RenderPreview() string {
	if len(m.filteredItems) == 0 || m.cursor >= len(m.filteredItems) {
		return ""
	}

	item := m.filteredItems[m.cursor]
	var lines []string

	lines = append(lines, item.RenderHeader(m.Theme))
	lines = append(lines, "")

	children := item.RenderChildren(m.Theme, m.settings)
	if len(children) > 0 {
		lines = append(lines, children...)
		lines = append(lines, "")
	}

	metadata := item.RenderMetadata(m.Theme, m.settings)
	lines = append(lines, metadata...)

	return strings.Join(lines, "\n")
}

// RenderSearch renders the search input
func (m ListSelectModel) RenderSearch(layout layoutMetrics, hasSize bool) string {
	icon := m.Theme.Icons.Search
	prompt := m.Theme.InputPromptStyle.Render(icon + " ")

	queryText := m.searchQuery
	if queryText == "" {
		queryText = "Search..."
	}

	cursor := "█"
	inputContent := prompt + queryText + cursor

	inputStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.Theme.BorderColor)

	if hasSize && layout.contentWidth > 0 {
		inputStyle = inputStyle.Width(layout.contentWidth)
	}

	if m.Theme.InputBgColor != lipgloss.Color("") {
		inputStyle = inputStyle.Background(m.Theme.InputBgColor)
	}

	return inputStyle.Render(inputContent)
}

// RenderTitle renders the title
func (m ListSelectModel) RenderTitle() string {
	return m.Theme.TitleStyle.Render(m.Title)
}
