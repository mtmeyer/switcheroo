package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// WorktreeSelectModel handles worktree selection using the generic ListSelectModel
type WorktreeSelectModel struct {
	list     ListSelectModel
	repoName string
}

// NewWorktreeSelectModel creates a new worktree selection model
func NewWorktreeSelectModel(repo RepoDisplay, theme Theme, settings PreviewSettings) WorktreeSelectModel {
	items := make([]Selectable, len(repo.Worktrees))
	for i := range repo.Worktrees {
		items[i] = repo.Worktrees[i]
	}

	onSelect := func(item Selectable) tea.Cmd {
		return worktreeSelectedCmd(item.(WorktreeDisplay))
	}

	return WorktreeSelectModel{
		list:     NewListSelectModel("Select Worktree: "+repo.Name, items, theme, settings, onSelect),
		repoName: repo.Name,
	}
}

// Init initializes the worktree select model
func (m WorktreeSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for worktree selection
func (m WorktreeSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "esc":
			return m, backToReposCmd()
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

// View renders the worktree selection view
func (m WorktreeSelectModel) View() string {
	helpText := "↑/k up • ↓/j down • enter select • esc back • q/ctrl+c quit"
	return RenderSelectionView(m.list, helpText)
}

// SetDimensions sets the terminal dimensions on the underlying list model
func (m *WorktreeSelectModel) SetDimensions(width, height int) {
	m.list.SetDimensions(width, height)
}
