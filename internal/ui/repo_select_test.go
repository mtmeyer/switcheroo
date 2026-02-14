package ui

import (
	"strings"
	"testing"
)

// Test ListSelectModel scroll behavior

func TestListSelectModel_AdjustScrollOffset_CursorWithinVisibleArea(t *testing.T) {
	items := []Selectable{
		RepoDisplay{Name: "repo1"},
		RepoDisplay{Name: "repo2"},
		RepoDisplay{Name: "repo3"},
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 3
	m.scrollOffset = 0
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_CursorMovesDown(t *testing.T) {
	items := make([]Selectable, 20)
	for i := 0; i < 20; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 18
	m.scrollOffset = 0
	m.adjustScrollOffset()
	expected := 1 // cursor 18 - visibleHeight 18 + 1 = 1
	if m.scrollOffset != expected {
		t.Errorf("Expected scrollOffset to be %d, got %d", expected, m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_CursorMovesUp(t *testing.T) {
	items := make([]Selectable, 20)
	for i := 0; i < 20; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 5
	m.scrollOffset = 10
	m.adjustScrollOffset()
	if m.scrollOffset != 5 {
		t.Errorf("Expected scrollOffset to be 5, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_CursorAtEnd(t *testing.T) {
	items := make([]Selectable, 100)
	for i := 0; i < 100; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 99
	m.scrollOffset = 0
	m.adjustScrollOffset()
	expected := 82 // 99 - 18 + 1 = 82
	if m.scrollOffset != expected {
		t.Errorf("Expected scrollOffset to be %d, got %d", expected, m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_EmptyList(t *testing.T) {
	m := NewListSelectModel("Test", []Selectable{}, Theme{}, PreviewSettings{}, nil)
	m.cursor = 0
	m.scrollOffset = 0
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 for empty list, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_SingleRepo(t *testing.T) {
	items := []Selectable{RepoDisplay{Name: "repo1"}}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 0
	m.scrollOffset = 0
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 for single repo, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_AdjustScrollOffset_ExactFit(t *testing.T) {
	items := make([]Selectable, 18)
	for i := 0; i < 18; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 17
	m.scrollOffset = 0
	m.adjustScrollOffset()
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to stay at 0 when list fits exactly, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_FilterItems_ResetsCursorWhenOutOfBounds(t *testing.T) {
	// Create repos with specific names that we can filter
	items := []Selectable{
		RepoDisplay{Name: "alpha"},
		RepoDisplay{Name: "beta"},
		RepoDisplay{Name: "gamma"},
		RepoDisplay{Name: "delta"},
		RepoDisplay{Name: "epsilon"},
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 4 // Last item (epsilon)
	m.scrollOffset = 0

	// Now filter to only items containing "a" - this should match alpha, beta, gamma, delta
	// Epsilon will be filtered out, so cursor at 4 will be out of bounds
	m.UpdateSearchQuery("a")

	if m.cursor != 0 {
		t.Errorf("Expected cursor to reset to 0 when out of bounds, got %d", m.cursor)
	}
	if m.scrollOffset != 0 {
		t.Errorf("Expected scrollOffset to reset to 0, got %d", m.scrollOffset)
	}
}

func TestListSelectModel_FilterItems_CursorWithinBounds(t *testing.T) {
	items := make([]Selectable, 20)
	for i := 0; i < 20; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 5
	m.scrollOffset = 0

	// Filter but cursor stays within bounds
	m.filteredItems = items[:15]
	m.filterItems()

	// Cursor should stay at 5 since it's within the new bounds
	if m.cursor != 5 {
		t.Errorf("Expected cursor to stay at 5, got %d", m.cursor)
	}
}

func TestListSelectModel_MoveCursor(t *testing.T) {
	items := make([]Selectable, 10)
	for i := 0; i < 10; i++ {
		items[i] = RepoDisplay{Name: "repo"}
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)
	m.cursor = 0

	m.MoveCursor(1)
	if m.cursor != 1 {
		t.Errorf("Expected cursor to be 1, got %d", m.cursor)
	}

	m.MoveCursor(-1)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to be 0, got %d", m.cursor)
	}

	// Try to go below 0
	m.MoveCursor(-1)
	if m.cursor != 0 {
		t.Errorf("Expected cursor to stay at 0, got %d", m.cursor)
	}

	// Try to go past end
	m.cursor = 9
	m.MoveCursor(1)
	if m.cursor != 9 {
		t.Errorf("Expected cursor to stay at 9, got %d", m.cursor)
	}
}

func TestListSelectModel_UpdateSearchQuery(t *testing.T) {
	items := []Selectable{
		RepoDisplay{Name: "alpha"},
		RepoDisplay{Name: "beta"},
		RepoDisplay{Name: "gamma"},
	}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)

	// Search for "al"
	m.UpdateSearchQuery("al")
	if m.searchQuery != "al" {
		t.Errorf("Expected searchQuery to be 'al', got %s", m.searchQuery)
	}
	// Should only match "alpha"
	if len(m.filteredItems) != 1 {
		t.Errorf("Expected 1 filtered item, got %d", len(m.filteredItems))
	}

	// Clear search
	m.ClearSearch()
	if m.searchQuery != "" {
		t.Errorf("Expected searchQuery to be empty, got %s", m.searchQuery)
	}
	if len(m.filteredItems) != 3 {
		t.Errorf("Expected 3 items after clearing, got %d", len(m.filteredItems))
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		maxWidth int
		expected string
	}{
		{"short", 10, "short"},
		{"exactlyten", 10, "exactlyten"},
		{"thisisaverylongreponame", 10, "thisisa..."},
		{"test", 3, "..."},
		{"test", 2, "..."},
		{"", 10, ""},
		{"ab", 5, "ab"},
	}

	for _, tt := range tests {
		result := truncateWithEllipsis(tt.input, tt.maxWidth)
		if result != tt.expected {
			t.Errorf("truncateWithEllipsis(%q, %d) = %q, want %q",
				tt.input, tt.maxWidth, result, tt.expected)
		}
	}
}

func TestRenderItemLine_Truncation(t *testing.T) {
	// Test that long names are truncated in renderItemLine
	longName := "this-is-a-very-long-repository-name-that-needs-truncation"
	items := []Selectable{RepoDisplay{Name: longName}}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)

	// Create a layout with limited width
	layout := layoutMetrics{listWidth: 40} // 40 - 9 = 31 chars for name

	// Render the line
	rendered := m.renderItemLine(items[0], false, layout)

	// The rendered output should contain "..." if the name is longer than available space
	if !strings.Contains(rendered, "...") {
		t.Errorf("Expected truncated name to contain '...', got: %s", rendered)
	}

	// Test with a short name - should not be truncated
	shortName := "short-repo"
	items2 := []Selectable{RepoDisplay{Name: shortName}}
	m2 := NewListSelectModel("Test", items2, Theme{}, PreviewSettings{}, nil)
	rendered2 := m2.renderItemLine(items2[0], false, layout)

	// Should contain the full short name
	if !strings.Contains(rendered2, shortName) {
		t.Errorf("Expected short name to be rendered fully, got: %s", rendered2)
	}
}

func TestRenderItemLine_DynamicWidth(t *testing.T) {
	longName := "this-is-a-very-long-repository-name-that-needs-truncation"
	items := []Selectable{RepoDisplay{Name: longName}}
	m := NewListSelectModel("Test", items, Theme{}, PreviewSettings{}, nil)

	// Test with narrow layout - should truncate
	narrowLayout := layoutMetrics{listWidth: 30} // 30 - 9 = 21 chars for name
	renderedNarrow := m.renderItemLine(items[0], false, narrowLayout)
	if !strings.Contains(renderedNarrow, "...") {
		t.Errorf("Expected truncation with narrow layout, got: %s", renderedNarrow)
	}

	// Test with wide layout - should not truncate
	wideLayout := layoutMetrics{listWidth: 100} // 100 - 9 = 91 chars for name
	renderedWide := m.renderItemLine(items[0], false, wideLayout)
	// With 91 chars available, the long name (56 chars) should fit without truncation
	if strings.Contains(renderedWide, "...") {
		t.Errorf("Did not expect truncation with wide layout, got: %s", renderedWide)
	}
}
