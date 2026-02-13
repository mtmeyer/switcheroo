package ui

import tea "github.com/charmbracelet/bubbletea"

// RepoSelectedMsg indicates a repository was chosen from the repo list
type RepoSelectedMsg struct {
	Repo RepoDisplay
}

// WorktreeSelectedMsg indicates a worktree was chosen from the worktree list
type WorktreeSelectedMsg struct {
	Worktree WorktreeDisplay
}

// BackToReposMsg indicates the user wants to return to the repo list
type BackToReposMsg struct{}

// PathSelectedMsg tells the root model that a final path was chosen
type PathSelectedMsg struct {
	Path string
}

func repoSelectedCmd(repo RepoDisplay) tea.Cmd {
	return func() tea.Msg {
		return RepoSelectedMsg{Repo: repo}
	}
}

func worktreeSelectedCmd(worktree WorktreeDisplay) tea.Cmd {
	return func() tea.Msg {
		return WorktreeSelectedMsg{Worktree: worktree}
	}
}

func backToReposCmd() tea.Cmd {
	return func() tea.Msg {
		return BackToReposMsg{}
	}
}

func pathSelectedCmd(path string) tea.Cmd {
	return func() tea.Msg {
		return PathSelectedMsg{Path: path}
	}
}
