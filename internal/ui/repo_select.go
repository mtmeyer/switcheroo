package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// RepoSelectModel handles repository selection using the generic ListSelectModel
type RepoSelectModel struct {
	list ListSelectModel
}

// NewRepoSelectModel creates a new repository selection model
func NewRepoSelectModel(repos []RepoDisplay, theme Theme, settings PreviewSettings) RepoSelectModel {
	items := make([]Selectable, len(repos))
	for i := range repos {
		items[i] = repos[i]
	}

	onSelect := func(item Selectable) tea.Cmd {
		if item.HasChildren() {
			// Repo has worktrees, emit message to transition
			return repoSelectedCmd(item.(RepoDisplay))
		}
		// No worktrees, emit path selected
		return pathSelectedCmd(item.OnSelect())
	}

	return RepoSelectModel{
		list: NewListSelectModel("Select Repository", items, theme, settings, onSelect),
	}
}

// Init initializes the repo select model
func (m RepoSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for repo selection
func (m RepoSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetDimensions(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.list.MoveCursor(-1)
		case "down", "j":
			m.list.MoveCursor(1)
		case "enter":
			return m, m.list.SelectCurrent()
		case "backspace":
			if len(m.list.GetSearchQuery()) > 0 {
				newQuery := m.list.GetSearchQuery()[:len(m.list.GetSearchQuery())-1]
				m.list.UpdateSearchQuery(newQuery)
			}
		default:
			if len(msg.String()) == 1 {
				newQuery := m.list.GetSearchQuery() + msg.String()
				m.list.UpdateSearchQuery(newQuery)
			}
		}
	}
	return m, nil
}

// View renders the repo selection view
func (m RepoSelectModel) View() string {
	layout := calculateLayout(m.list.Width)
	hasSize := m.list.Width > 0 && m.list.Height > 0

	// Build components
	title := m.list.RenderTitle()
	search := m.list.RenderSearch(layout, hasSize)
	list := m.list.RenderList(layout)
	preview := m.list.RenderPreview()

	// Combine list and preview side by side
	content := renderSideBySide(list, preview, layout, m.list.Theme)

	// Combine all sections
	sections := []string{title, search, content, m.renderHelp()}
	inner := strings.Join(sections, "\n")

	// Wrap in container
	return wrapInContainer(inner, layout, hasSize, m.list.Theme, m.list.Width, m.list.Height)
}

func (m RepoSelectModel) renderHelp() string {
	helpText := "↑/k up • ↓/j down • enter select • q/ctrl+c quit"
	return m.list.Theme.HelpStyle.Render(helpText)
}
