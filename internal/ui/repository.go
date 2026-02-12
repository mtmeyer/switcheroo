package ui

import (
	"sort"
	"switcheroo/internal/git"
)

// RepoDisplay is an adapter between git.Repository and UI needs
type RepoDisplay struct {
	Name          string
	Path          string
	Worktrees     []string // Worktree names
	Branches      []string // Branch names
	CurrentBranch string
	HasWorktrees  bool
}

// FromGitRepositories converts git.Repository slice to RepoDisplay slice
func FromGitRepositories(repos []git.Repository) []RepoDisplay {
	displays := make([]RepoDisplay, len(repos))

	for i, repo := range repos {
		displays[i] = RepoDisplay{
			Name:         repo.Name,
			Path:         repo.Path,
			HasWorktrees: repo.HasWorktrees,
		}

		if repo.HasWorktrees {
			// Extract worktree names
			worktreeNames := make([]string, len(repo.Worktrees))
			for j, wt := range repo.Worktrees {
				worktreeNames[j] = wt.Name
			}
			displays[i].Worktrees = worktreeNames
		} else {
			// Extract branch names and sort them
			branchNames := make([]string, len(repo.Branches))
			for j, branch := range repo.Branches {
				branchNames[j] = branch.Name
				if branch.IsCurrent {
					displays[i].CurrentBranch = branch.Name
				}
			}
			sort.Strings(branchNames)
			displays[i].Branches = branchNames
		}
	}

	return displays
}
