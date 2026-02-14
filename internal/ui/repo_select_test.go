package ui

import (
	"fmt"
	"testing"
)

func TestAdjustScrollOffset_CursorWithinVisibleArea(t *testing.T) {
	m := &RepoSelectModel{
		cursor:        3,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 20),
	}
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0, got %d", m.scrollOffset)
	}
}

func TestAdjustScrollOffset_CursorMovesDown(t *testing.T) {
	// 20 repos, visible height is 18 (paneHeight 20 - 2 padding)
	// If cursor moves to position 18, we need to scroll
	m := &RepoSelectModel{
		cursor:        18,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 20),
	}
	m.adjustScrollOffset()
	expected := 1 // cursor 18 - visibleHeight 18 + 1 = 1
	if m.scrollOffset != expected {
		t.Errorf("Expected scrollOffset to be %d, got %d", expected, m.scrollOffset)
	}
}

func TestAdjustScrollOffset_CursorMovesUp(t *testing.T) {
	// If we're scrolled down and cursor moves up above visible area
	m := &RepoSelectModel{
		cursor:        5,
		scrollOffset:  10,
		filteredRepos: make([]RepoDisplay, 20),
	}
	m.adjustScrollOffset()
	if m.scrollOffset != 5 {
		t.Errorf("Expected scrollOffset to be 5, got %d", m.scrollOffset)
	}
}

func TestAdjustScrollOffset_CursorAtEnd(t *testing.T) {
	// Cursor at last item
	m := &RepoSelectModel{
		cursor:        99,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 100),
	}
	m.adjustScrollOffset()
	expected := 82 // 99 - 18 + 1 = 82
	if m.scrollOffset != expected {
		t.Errorf("Expected scrollOffset to be %d, got %d", expected, m.scrollOffset)
	}
}

func TestAdjustScrollOffset_EmptyList(t *testing.T) {
	m := &RepoSelectModel{
		cursor:        0,
		scrollOffset:  0,
		filteredRepos: []RepoDisplay{},
	}
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 for empty list, got %d", m.scrollOffset)
	}
}

func TestAdjustScrollOffset_SingleRepo(t *testing.T) {
	m := &RepoSelectModel{
		cursor:        0,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 1),
	}
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 for single repo, got %d", m.scrollOffset)
	}
}

func TestAdjustScrollOffset_ExactFit(t *testing.T) {
	// Exactly 18 repos (visible height)
	m := &RepoSelectModel{
		cursor:        17,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 18),
	}
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 when list fits exactly, got %d", m.scrollOffset)
	}
}

func TestFilterRepos_ResetsCursorWhenOutOfBounds(t *testing.T) {
	// Create repos with names that will be filtered
	repos := make([]RepoDisplay, 20)
	for i := 0; i < 20; i++ {
		repos[i] = RepoDisplay{Name: fmt.Sprintf("repo-%d", i)}
	}

	m := &RepoSelectModel{
		cursor:        15,
		scrollOffset:  5,
		searchQuery:   "", // Will match all repos initially
		filteredRepos: repos,
		repos:         repos,
	}

	// Now set a search query that will only match first 10 repos
	m.searchQuery = "repo-0" // Will match repo-0, repo-01, repo-02, ..., repo-09
	m.filterRepos()

	// Cursor should reset since it was at 15 but now only 10 repos match
	if m.cursor != 0 {
		t.Errorf("Expected cursor to reset to 0 when out of bounds, got %d", m.cursor)
	}
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to reset to 0, got %d", m.scrollOffset)
	}
}

func TestFilterRepos_CursorWithinBounds(t *testing.T) {
	m := &RepoSelectModel{
		cursor:        5,
		scrollOffset:  0,
		filteredRepos: make([]RepoDisplay, 20),
		repos:         make([]RepoDisplay, 20),
	}

	// Filtering but cursor stays within bounds
	m.filteredRepos = make([]RepoDisplay, 15)
	m.filterRepos()

	// Cursor should stay at 5 since it's within the new bounds
	if m.cursor != 5 {
		t.Errorf("Expected cursor to stay at 5, got %d", m.cursor)
	}
}

func TestFilterRepos_ClearSearch(t *testing.T) {
	m := &RepoSelectModel{
		cursor:        3,
		scrollOffset:  1,
		searchQuery:   "test",
		filteredRepos: make([]RepoDisplay, 5),
		repos:         make([]RepoDisplay, 10),
	}

	// Clear search
	m.searchQuery = ""
	m.filterRepos()

	// Should restore all repos
	if len(m.filteredRepos) != 10 {
		t.Errorf("Expected 10 repos after clearing search, got %d", len(m.filteredRepos))
	}
}
