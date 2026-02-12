package git

// Repository represents a Git repository
type Repository struct {
	Name          string
	Path          string
	HasWorktrees  bool
	Worktrees     []Worktree
	Branches      []Branch
	CurrentBranch string
	RemoteURL     string
}

// TODO: Implement repository discovery
