package ui

import (
	"sort"
	"strings"

	"switcheroo/internal/git"
)

// RepoDisplay is an adapter between git.Repository and UI needs
type RepoDisplay struct {
	Name            string
	Path            string
	Worktrees       []WorktreeDisplay
	Branches        []BranchDisplay
	CurrentBranch   string
	HasWorktrees    bool
	LineDiffAdded   int
	LineDiffRemoved int
	DefaultBranch   string
}

type BranchDisplay struct {
	Name        string
	IsCurrent   bool
	Status      string
	DiffAdded   int
	DiffRemoved int
}

// WorktreeDisplay represents worktree data needed by the UI
type WorktreeDisplay struct {
	Name        string
	Path        string
	Branch      string
	Locked      bool
	Status      string
	DiffAdded   int
	DiffRemoved int
}

// FromGitRepositories converts git.Repository slice to RepoDisplay slice
func FromGitRepositories(repos []git.Repository) []RepoDisplay {
	displays := make([]RepoDisplay, len(repos))

	for i, repo := range repos {
		displays[i] = RepoDisplay{
			Name:          repo.Name,
			Path:          repo.Path,
			HasWorktrees:  repo.HasWorktrees,
			DefaultBranch: repo.DefaultBranch,
		}

		if repo.HasWorktrees {
			worktrees := make([]WorktreeDisplay, len(repo.Worktrees))
			for j, wt := range repo.Worktrees {
				worktrees[j] = WorktreeDisplay{
					Name:        wt.Name,
					Path:        wt.Path,
					Branch:      wt.Branch,
					Locked:      wt.IsLocked,
					Status:      wt.Status,
					DiffAdded:   wt.DiffAdded,
					DiffRemoved: wt.DiffRemoved,
				}
			}
			displays[i].Worktrees = worktrees
		} else {
			branches := make([]BranchDisplay, len(repo.Branches))
			for j, branch := range repo.Branches {
				branches[j] = BranchDisplay{
					Name:        branch.Name,
					IsCurrent:   branch.IsCurrent,
					Status:      branch.Status,
					DiffAdded:   branch.DiffAdded,
					DiffRemoved: branch.DiffRemoved,
				}
				if branch.IsCurrent {
					displays[i].CurrentBranch = branch.Name
					if branch.Status == "untracked" && repo.DefaultBranch != "" && repo.DefaultBranch != branch.Name {
						if added, removed, err := git.DiffBranches(repo.Path, branch.Name, repo.DefaultBranch); err == nil {
							displays[i].LineDiffAdded = added
							displays[i].LineDiffRemoved = removed
						}
					} else {
						displays[i].LineDiffAdded = branch.DiffAdded
						displays[i].LineDiffRemoved = branch.DiffRemoved
					}
				}
			}
			sort.Slice(branches, func(a, b int) bool {
				return strings.ToLower(branches[a].Name) < strings.ToLower(branches[b].Name)
			})
			displays[i].Branches = branches
		}
	}

	return displays
}
