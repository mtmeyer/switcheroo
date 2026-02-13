package ui

import (
	"sort"
	"switcheroo/internal/git"
)

// RepoDisplay is an adapter between git.Repository and UI needs
type RepoDisplay struct {
	Name          string
	Path          string
	Worktrees     []WorktreeDisplay
	Branches      []string // Branch names
	CurrentBranch string
	HasWorktrees  bool
}

// WorktreeDisplay represents worktree data needed by the UI
type WorktreeDisplay struct {
	Name   string
	Path   string
	Branch string
	Locked bool
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
			worktrees := make([]WorktreeDisplay, len(repo.Worktrees))
			for j, wt := range repo.Worktrees {
				worktrees[j] = WorktreeDisplay{
					Name:   wt.Name,
					Path:   wt.Path,
					Branch: wt.Branch,
					Locked: wt.IsLocked,
				}
			}
			displays[i].Worktrees = worktrees
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
