package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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
	for i := range worktrees {
		wt := &worktrees[i]
		if wt.Path == repoPath {
			continue
		}
		populateWorktreeMetadata(wt)
		additionalWorktrees = append(additionalWorktrees, *wt)
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
	format := "%(refname:short)|%(HEAD)|%(upstream:short)|%(upstream:track)"
	cmd := exec.Command("git", "-C", repoPath, "branch", "--format="+format)
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
		if len(parts) < 4 {
			continue
		}
		branch := Branch{
			Name:      parts[0],
			IsCurrent: parts[1] == "*",
			Upstream:  parts[2],
		}
		track := parts[3]
		branch.Status = determineBranchStatus(branch.Upstream, track)
		if branch.Upstream != "" {
			added, removed, err := diffAgainstUpstream(repoPath, branch.Name, branch.Upstream)
			if err == nil {
				branch.DiffAdded = added
				branch.DiffRemoved = removed
			}
		}
		branches = append(branches, branch)
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

func populateWorktreeMetadata(wt *Worktree) {
	status, err := workingTreeStatus(wt.Path)
	if err == nil {
		wt.Status = status
	}
	added, removed, err := diffInWorktree(wt.Path)
	if err == nil {
		wt.DiffAdded = added
		wt.DiffRemoved = removed
	}
}

func workingTreeStatus(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return "clean", nil
	}
	if strings.Contains(trimmed, "??") {
		return "untracked", nil
	}
	return "modified", nil
}

func diffInWorktree(path string) (int, int, error) {
	cmd := exec.Command("git", "-C", path, "diff", "--shortstat")
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	added, removed := parseShortstat(string(output))
	return added, removed, nil
}

func diffAgainstUpstream(repoPath, branch, upstream string) (int, int, error) {
	if branch == "" || upstream == "" {
		return 0, 0, nil
	}
	arg := branch + "..." + upstream
	cmd := exec.Command("git", "-C", repoPath, "diff", "--shortstat", arg)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	added, removed := parseShortstat(string(output))
	return added, removed, nil
}

func DiffBranches(repoPath, left, right string) (int, int, error) {
	if left == "" || right == "" {
		return 0, 0, nil
	}
	arg := left + "..." + right
	cmd := exec.Command("git", "-C", repoPath, "diff", "--shortstat", arg)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	added, removed := parseShortstat(string(output))
	return added, removed, nil
}

func determineBranchStatus(upstream, track string) string {
	if upstream == "" {
		return "untracked"
	}
	track = strings.ToLower(track)
	switch {
	case strings.Contains(track, "ahead") && strings.Contains(track, "behind"):
		return "diverged"
	case strings.Contains(track, "ahead"):
		return "ahead"
	case strings.Contains(track, "behind"):
		return "behind"
	case strings.Contains(track, "up to date"):
		return "clean"
	default:
		return "clean"
	}
}

var (
	insertionsRe = regexp.MustCompile(`(?P<count>\d+)\s+insertions?\(\+\)`)
	deletionsRe  = regexp.MustCompile(`(?P<count>\d+)\s+deletions?\(-\)`)
)

func parseShortstat(output string) (int, int) {
	added := extractCount(insertionsRe, output)
	removed := extractCount(deletionsRe, output)
	return added, removed
}

func extractCount(re *regexp.Regexp, output string) int {
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return 0
	}
	return atoiSafe(match[1])
}

func atoiSafe(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return i
}

func GetDefaultBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(output))
		if strings.HasPrefix(ref, "refs/remotes/origin/") {
			candidate := strings.TrimPrefix(ref, "refs/remotes/origin/")
			if branchExists(repoPath, candidate) {
				return candidate, nil
			}
		}
	}

	candidates := []string{"main", "master"}
	for _, candidate := range candidates {
		if branchExists(repoPath, candidate) {
			return candidate, nil
		}
	}

	if current, err := GetCurrentBranch(repoPath); err == nil && current != "" {
		return current, nil
	}

	return "", fmt.Errorf("could not determine default branch")
}

func branchExists(repoPath, branch string) bool {
	if branch == "" {
		return false
	}
	cmd := exec.Command("git", "-C", repoPath, "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	return cmd.Run() == nil
}
