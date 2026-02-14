package git

// Worktree represents a Git worktree
type Worktree struct {
	Name        string
	Path        string
	Branch      string
	IsLocked    bool
	IsBare      bool
	CommitHash  string
	Status      string
	DiffAdded   int
	DiffRemoved int
}

// TODO: Implement worktree detection
