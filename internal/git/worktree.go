package git

// Worktree represents a Git worktree
type Worktree struct {
	Name       string
	Path       string
	Branch     string
	IsLocked   bool
	IsBare     bool
	CommitHash string
}

// TODO: Implement worktree detection
