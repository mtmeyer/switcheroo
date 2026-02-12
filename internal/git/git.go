package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitRepository checks if a directory is a git repository
func IsGitRepository(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetWorktrees returns a list of worktrees for a repository
// Returns nil if the repo doesn't use worktrees (only has main worktree)
func GetWorktrees(repoPath string) ([]Worktree, error) {
	cmd := exec.Command("git", "-C", repoPath, "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	worktrees := parseWorktreeList(string(output))

	// Filter out the main worktree (the repository itself)
	// Only return additional worktrees
	var additionalWorktrees []Worktree
	for _, wt := range worktrees {
		if wt.Path != repoPath {
			additionalWorktrees = append(additionalWorktrees, wt)
		}
	}

	return additionalWorktrees, nil
}

// parseWorktreeList parses the output of `git worktree list --porcelain`
func parseWorktreeList(output string) []Worktree {
	var worktrees []Worktree
	var current *Worktree

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			if current != nil {
				worktrees = append(worktrees, *current)
				current = nil
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current = &Worktree{
				Path: strings.TrimPrefix(line, "worktree "),
			}
			current.Name = filepath.Base(current.Path)
		} else if strings.HasPrefix(line, "branch ") && current != nil {
			branchRef := strings.TrimPrefix(line, "branch ")
			// Remove refs/heads/ prefix
			current.Branch = strings.TrimPrefix(branchRef, "refs/heads/")
		} else if strings.HasPrefix(line, "locked") && current != nil {
			current.IsLocked = true
		} else if strings.HasPrefix(line, "bare") && current != nil {
			current.IsBare = true
		}
	}

	// Don't forget the last one
	if current != nil {
		worktrees = append(worktrees, *current)
	}

	return worktrees
}

// GetBranches returns all local branches for a repository
func GetBranches(repoPath string) ([]Branch, error) {
	// Get branch list with format: name|isCurrent
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)|%(HEAD)")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var branches []Branch
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		branchName := parts[0]
		isCurrent := parts[1] == "*"

		branches = append(branches, Branch{
			Name:      branchName,
			IsCurrent: isCurrent,
		})
	}

	return branches, nil
}

// GetCurrentBranch returns the currently checked out branch
func GetCurrentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
