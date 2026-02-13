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
	theme               Theme
	selectedPath        string
	err                 error
}

// NewModel creates a new root model using discovered repositories and theme
func NewModel(repos []RepoDisplay, theme Theme) Model {
	return Model{
		state:           StateRepoSelect,
		repoSelectModel: NewRepoSelectModel(repos, theme),
		theme:           theme,
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
		switch m.state {
		case StateRepoSelect:
			updated, cmd := m.repoSelectModel.Update(msg)
			m.repoSelectModel = updated
			return m, cmd
		case StateWorktreeSelect:
			if m.worktreeSelectModel != nil {
				updated, cmd := m.worktreeSelectModel.Update(msg)
				m.worktreeSelectModel = updated
				return m, cmd
			}
		}
	case RepoSelectedMsg:
		return m.handleRepoSelected(msg.Repo)
	case WorktreeSelectedMsg:
		return m, pathSelectedCmd(msg.Worktree.Path)
	case BackToReposMsg:
		m.state = StateRepoSelect
		m.worktreeSelectModel = nil
		return m, nil
	case PathSelectedMsg:
		m.selectedPath = msg.Path
		m.state = StateComplete
		return m, tea.Quit
	}

	switch m.state {
	case StateRepoSelect:
		updated, cmd := m.repoSelectModel.Update(msg)
		m.repoSelectModel = updated
		return m, cmd
	case StateWorktreeSelect:
		if m.worktreeSelectModel != nil {
			updated, cmd := m.worktreeSelectModel.Update(msg)
			m.worktreeSelectModel = updated
			return m, cmd
		}
	}

	return m, nil
}

func (m Model) handleRepoSelected(repo RepoDisplay) (tea.Model, tea.Cmd) {
	if repo.HasWorktrees && len(repo.Worktrees) > 0 {
		wtModel := NewWorktreeSelectModel(repo, m.theme)
		m.worktreeSelectModel = wtModel
		m.state = StateWorktreeSelect
		return m, nil
	}
	return m, pathSelectedCmd(repo.Path)
}

// View renders the model
func (m Model) View() string {
	switch m.state {
	case StateRepoSelect:
		return m.repoSelectModel.View()
	case StateWorktreeSelect:
		if m.worktreeSelectModel != nil {
			return m.worktreeSelectModel.View()
		}
		return ""
	case StateError:
		return "Error: " + m.err.Error()
	default:
		return ""
	}
}

// SelectedPath returns the final path chosen by the user
func (m Model) SelectedPath() string {
	return m.selectedPath
}
