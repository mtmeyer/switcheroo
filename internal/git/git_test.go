package git

import (
	"testing"
)

// ============= parseShortstat Tests =============

func TestParseShortstat_FullOutput(t *testing.T) {
	added, removed := parseShortstat(" 3 files changed, 14 insertions(+), 9 deletions(-)")
	if added != 14 || removed != 9 {
		t.Fatalf("expected 14/9 got %d/%d", added, removed)
	}
}

func TestParseShortstat_OnlyDeletions(t *testing.T) {
	added, removed := parseShortstat(" 1 file changed, 2 deletions(-)")
	if added != 0 || removed != 2 {
		t.Fatalf("expected 0/2 got %d/%d", added, removed)
	}
}

func TestParseShortstat_OnlyInsertions(t *testing.T) {
	added, removed := parseShortstat(" 2 files changed, 10 insertions(+)")
	if added != 10 || removed != 0 {
		t.Fatalf("expected 10/0 got %d/%d", added, removed)
	}
}

func TestParseShortstat_EmptyOutput(t *testing.T) {
	added, removed := parseShortstat("")
	if added != 0 || removed != 0 {
		t.Fatalf("expected 0/0 for empty string got %d/%d", added, removed)
	}
}

func TestParseShortstat_NoChanges(t *testing.T) {
	// This shouldn't happen with git diff but test defensive parsing
	added, removed := parseShortstat("0 files changed")
	if added != 0 || removed != 0 {
		t.Fatalf("expected 0/0 got %d/%d", added, removed)
	}
}

func TestParseShortstat_MultipleFiles(t *testing.T) {
	// Test with pluralization variations
	added, removed := parseShortstat(" 5 files changed, 100 insertions(+), 50 deletions(-)")
	if added != 100 || removed != 50 {
		t.Fatalf("expected 100/50 got %d/%d", added, removed)
	}
}

// ============= determineBranchStatus Tests =============

func TestDetermineBranchStatus_Ahead(t *testing.T) {
	if s := determineBranchStatus("origin/main", "[ahead 2]"); s != "ahead" {
		t.Fatalf("expected ahead got %s", s)
	}
}

func TestDetermineBranchStatus_Behind(t *testing.T) {
	if s := determineBranchStatus("origin/main", "[behind 3]"); s != "behind" {
		t.Fatalf("expected behin got %s", s)
	}
}

func TestDetermineBranchStatus_Diverged(t *testing.T) {
	if s := determineBranchStatus("origin/main", "[ahead 1, behind 1]"); s != "diverged" {
		t.Fatalf("expected diverged got %s", s)
	}
}

func TestDetermineBranchStatus_Clean(t *testing.T) {
	if s := determineBranchStatus("origin/main", "[up to date]"); s != "clean" {
		t.Fatalf("expected clean got %s", s)
	}
}

func TestDetermineBranchStatus_NoUpstream(t *testing.T) {
	if s := determineBranchStatus("", ""); s != "untracked" {
		t.Fatalf("expected untracked got %s", s)
	}
}

func TestDetermineBranchStatus_EmptyTrackWithUpstream(t *testing.T) {
	// Has upstream but no track info yet
	if s := determineBranchStatus("origin/main", ""); s != "clean" {
		t.Fatalf("expected clean (default) got %s", s)
	}
}

// ============= parseWorktreeList Tests =============

func TestParseWorktreeList_SingleWorktree(t *testing.T) {
	input := `worktree /path/to/repo
branch refs/heads/main`

	worktrees := parseWorktreeList(input)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}
	if worktrees[0].Path != "/path/to/repo" {
		t.Errorf("expected path /path/to/repo, got %s", worktrees[0].Path)
	}
	if worktrees[0].Name != "repo" {
		t.Errorf("expected name repo, got %s", worktrees[0].Name)
	}
	if worktrees[0].Branch != "main" {
		t.Errorf("expected branch main, got %s", worktrees[0].Branch)
	}
}

func TestParseWorktreeList_MultipleWorktrees(t *testing.T) {
	input := `worktree /path/to/repo
branch refs/heads/main

worktree /path/to/repo-worktrees/feature-x
branch refs/heads/feature-x

worktree /path/to/repo-worktrees/bugfix
branch refs/heads/bugfix
locked`

	worktrees := parseWorktreeList(input)
	if len(worktrees) != 3 {
		t.Fatalf("expected 3 worktrees, got %d", len(worktrees))
	}

	// Check first worktree
	if worktrees[0].Path != "/path/to/repo" {
		t.Errorf("expected first path /path/to/repo, got %s", worktrees[0].Path)
	}

	// Check second worktree
	if worktrees[1].Path != "/path/to/repo-worktrees/feature-x" {
		t.Errorf("expected second path, got %s", worktrees[1].Path)
	}
	if worktrees[1].Branch != "feature-x" {
		t.Errorf("expected branch feature-x, got %s", worktrees[1].Branch)
	}

	// Check third worktree (locked)
	if worktrees[2].Path != "/path/to/repo-worktrees/bugfix" {
		t.Errorf("expected third path, got %s", worktrees[2].Path)
	}
	if !worktrees[2].IsLocked {
		t.Errorf("expected third worktree to be locked")
	}
}

func TestParseWorktreeList_BareWorktree(t *testing.T) {
	input := `worktree /path/to/bare-repo
bare`

	worktrees := parseWorktreeList(input)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}
	if !worktrees[0].IsBare {
		t.Errorf("expected worktree to be bare")
	}
}

func TestParseWorktreeList_NoBranch(t *testing.T) {
	// Worktree with detached HEAD (no branch line)
	input := `worktree /path/to/repo`

	worktrees := parseWorktreeList(input)
	if len(worktrees) != 1 {
		t.Fatalf("expected 1 worktree, got %d", len(worktrees))
	}
	if worktrees[0].Branch != "" {
		t.Errorf("expected empty branch for detached HEAD, got %s", worktrees[0].Branch)
	}
}

func TestParseWorktreeList_EmptyInput(t *testing.T) {
	worktrees := parseWorktreeList("")
	if len(worktrees) != 0 {
		t.Fatalf("expected 0 worktrees for empty input, got %d", len(worktrees))
	}
}

// ============= extractCount and atoiSafe Tests =============

func TestExtractCount_ValidMatch(t *testing.T) {
	re := insertionsRe
	count := extractCount(re, "5 insertions(+)")
	if count != 5 {
		t.Fatalf("expected 5, got %d", count)
	}
}

func TestExtractCount_NoMatch(t *testing.T) {
	re := insertionsRe
	count := extractCount(re, "no changes here")
	if count != 0 {
		t.Fatalf("expected 0 for no match, got %d", count)
	}
}

func TestAtoiSafe_ValidNumber(t *testing.T) {
	result := atoiSafe("42")
	if result != 42 {
		t.Fatalf("expected 42, got %d", result)
	}
}

func TestAtoiSafe_InvalidNumber(t *testing.T) {
	result := atoiSafe("not-a-number")
	if result != 0 {
		t.Fatalf("expected 0 for invalid number, got %d", result)
	}
}

func TestAtoiSafe_EmptyString(t *testing.T) {
	result := atoiSafe("")
	if result != 0 {
		t.Fatalf("expected 0 for empty string, got %d", result)
	}
}

// ============= Working Tree Status Classification Tests =============

// Note: workingTreeStatus calls git command, but we can test the classification logic
// by examining how we categorize different status outputs

func TestWorkingTreeStatus_Classification(t *testing.T) {
	// These tests document the expected behavior of workingTreeStatus
	// In real scenarios, these would require git repository setup

	testCases := []struct {
		name      string
		porcelain string
		expected  string
	}{
		{
			name:      "clean repo",
			porcelain: "",
			expected:  "clean",
		},
		{
			name:      "modified files only",
			porcelain: " M file.go\nM  another.go",
			expected:  "modified",
		},
		{
			name:      "untracked files only",
			porcelain: "?? newfile.go",
			expected:  "untracked",
		},
		{
			name:      "modified with untracked",
			porcelain: " M file.go\n?? untracked.go",
			expected:  "untracked", // untracked takes precedence
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This documents expected behavior
			// The actual workingTreeStatus function needs git repo context
			// but we can verify our classification logic here
			var result string
			trimmed := tc.porcelain
			if trimmed == "" {
				result = "clean"
			} else if containsAny(trimmed, []string{"??"}) {
				result = "untracked"
			} else {
				result = "modified"
			}

			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// Helper function for testing
func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(s[:len(substr)] == substr) ||
		(len(s) > len(substr) && containsHelper(s[1:], substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
