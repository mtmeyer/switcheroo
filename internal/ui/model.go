package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// AppState represents the current state of the application
type AppState int

const (
	StateRepoSelect AppState = iota
	StateWorktreeSelect
	StateComplete
	StateError
)

// Model is the root Bubble Tea model
type Model struct {
	state               AppState
	repoSelectModel     tea.Model
	worktreeSelectModel tea.Model
	err                 error
}

// NewModel creates a new root model using discovered repositories and theme
func NewModel(repos []RepoDisplay, theme Theme) Model {
	return Model{
		state:           StateRepoSelect,
		repoSelectModel: NewRepoSelectModel(repos, theme),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		// Forward window size to child models
		if m.state == StateRepoSelect {
			updated, cmd := m.repoSelectModel.Update(msg)
			m.repoSelectModel = updated
			return m, cmd
		}
	}

	// Delegate to current state's model
	switch m.state {
	case StateRepoSelect:
		updated, cmd := m.repoSelectModel.Update(msg)
		m.repoSelectModel = updated
		return m, cmd
	case StateWorktreeSelect:
		updated, cmd := m.worktreeSelectModel.Update(msg)
		m.worktreeSelectModel = updated
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m Model) View() string {
	switch m.state {
	case StateRepoSelect:
		return m.repoSelectModel.View()
	case StateWorktreeSelect:
		return m.worktreeSelectModel.View()
	case StateError:
		return "Error: " + m.err.Error()
	default:
		return ""
	}
}
