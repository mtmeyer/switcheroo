package git

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Branch struct {
	Name        string
	IsCurrent   bool
	Upstream    string
	Status      string
	DiffAdded   int
	DiffRemoved int
}

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

type Repository struct {
	Name          string
	Path          string
	HasWorktrees  bool
	Worktrees     []Worktree
	Branches      []Branch
	CurrentBranch string
	RemoteURL     string
	DefaultBranch string
}

func DiscoverRepositories(directory string) ([]Repository, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var repos []Repository

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		repoPath := filepath.Join(directory, entry.Name())

		// Check if it's a git repository
		if !IsGitRepository(repoPath) {
			continue
		}

		repo, err := LoadRepository(repoPath)
		if err != nil {
			// Skip repos that fail to load
			continue
		}

		repos = append(repos, *repo)
	}

	// Sort repositories alphabetically by name
	sort.Slice(repos, func(i, j int) bool {
		return strings.ToLower(repos[i].Name) < strings.ToLower(repos[j].Name)
	})

	return repos, nil
}

// LoadRepository loads a repository and determines if it uses worktrees
func LoadRepository(path string) (*Repository, error) {
	repo := &Repository{
		Name: filepath.Base(path),
		Path: path,
	}

	if defaultBranch, err := GetDefaultBranch(path); err == nil {
		repo.DefaultBranch = defaultBranch
	}

	// Try to get worktrees
	worktrees, err := GetWorktrees(path)
	if err == nil && len(worktrees) > 0 {
		// This repo uses worktrees
		repo.HasWorktrees = true
		repo.Worktrees = worktrees
	} else {
		// Regular repo without worktrees, get branches
		repo.HasWorktrees = false

		branches, err := GetBranches(path)
		if err == nil {
			repo.Branches = branches
		}

		currentBranch, err := GetCurrentBranch(path)
		if err == nil {
			repo.CurrentBranch = currentBranch
		}
	}

	return repo, nil
}
