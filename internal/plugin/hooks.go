package plugin

import "switcheroo/internal/git"

// Hooks represents plugin hook points throughout the application
type Hooks struct {
	// Metadata hooks
	MetadataProviders []MetadataProvider

	// Filter hooks
	PreRepoSelect  []func(repos []git.Repository) []git.Repository
	PostRepoSelect []func(repo git.Repository) error

	PreWorktreeSelect  []func(worktrees []git.Worktree) []git.Worktree
	PostWorktreeSelect []func(worktree git.Worktree) error
}

// GlobalHooks is the global hooks instance (empty for now)
var GlobalHooks = &Hooks{
	MetadataProviders:  []MetadataProvider{},
	PreRepoSelect:      []func([]git.Repository) []git.Repository{},
	PostRepoSelect:     []func(git.Repository) error{},
	PreWorktreeSelect:  []func([]git.Worktree) []git.Worktree{},
	PostWorktreeSelect: []func(git.Worktree) error{},
}

// Plugin system will be implemented in the future
