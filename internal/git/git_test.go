package git

import "testing"

func TestParseShortstat(t *testing.T) {
	added, removed := parseShortstat(" 3 files changed, 14 insertions(+), 9 deletions(-)")
	if added != 14 || removed != 9 {
		t.Fatalf("expected 14/9 got %d/%d", added, removed)
	}
	added, removed = parseShortstat(" 1 file changed, 2 deletions(-)")
	if added != 0 || removed != 2 {
		t.Fatalf("expected 0/2 got %d/%d", added, removed)
	}
}

func TestDetermineBranchStatus(t *testing.T) {
	if s := determineBranchStatus("origin/main", "[ahead 2]"); s != "ahead" {
		t.Fatalf("expected ahead got %s", s)
	}
	if s := determineBranchStatus("origin/main", "[ahead 1, behind 1]"); s != "diverged" {
		t.Fatalf("expected diverged got %s", s)
	}
	if s := determineBranchStatus("origin/main", "[up to date]"); s != "clean" {
		t.Fatalf("expected clean got %s", s)
	}
	if s := determineBranchStatus("", ""); s != "untracked" {
		t.Fatalf("expected untracked got %s", s)
	}
}
