package plugin

import "switcheroo/internal/git"

// MetadataProvider allows plugins to add custom metadata fields
type MetadataProvider interface {
	Name() string
	GetMetadata(item interface{}) (key string, value string, error error)
}

// FilterProvider allows plugins to filter/transform lists
type FilterProvider interface {
	Name() string
	FilterRepositories(repos []git.Repository) []git.Repository
	FilterWorktrees(worktrees []git.Worktree) []git.Worktree
}

// ActionProvider allows plugins to run on selection
type ActionProvider interface {
	Name() string
	OnRepoSelect(repo git.Repository) error
	OnWorktreeSelect(worktree git.Worktree) error
}

// Plugin system will be implemented in the future
