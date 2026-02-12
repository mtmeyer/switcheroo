package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	// Max dimensions for the UI
	maxWidth  = 120
	maxHeight = 30

	// Component heights
	headerHeight = 1
	inputHeight  = 3
	helpHeight   = 1
	padding      = 2

	// List sizing
	listPaneWidth    = 40
	previewPaneWidth = 60
	gap              = 2
)

// Mock data structures (will be replaced with real types later)
type mockRepo struct {
	name          string
	path          string
	worktrees     []string // List of worktree names
	branches      []string
	currentBranch string
	status        string
}

// RepoSelectModel handles repository selection
type RepoSelectModel struct {
	repos         []mockRepo
	filteredRepos []mockRepo
	cursor        int
	scrollOffset  int
	searchQuery   string
	searchFocused bool
	selected      bool
	theme         Theme
	width         int
	height        int
}

// NewRepoSelectModel creates a new repository selection model with mock data
func NewRepoSelectModel() RepoSelectModel {
	repos := []mockRepo{
		{
			name:          "my-project",
			path:          "/Users/me/repos/my-project",
			worktrees:     []string{"feature-new-ui", "fix-bug-123", "staging"},
			branches:      []string{"feature/new-ui", "fix/bug-123", "main", "staging"},
			currentBranch: "feature/new-ui",
			status:        "Clean",
		},
		{
			name:          "another-repo",
			path:          "/Users/me/repos/another-repo",
			worktrees:     []string{}, // No worktrees
			branches:      []string{"develop", "main"},
			currentBranch: "main",
			status:        "Modified (2 files)",
		},
		{
			name:          "third-project",
			path:          "/Users/me/repos/third-project",
			worktrees:     []string{"experiment-refactor", "main"},
			branches:      []string{"experiment/refactor", "main"},
			currentBranch: "main",
			status:        "Clean",
		},
		// Add more mock data to test scrolling
		{
			name:          "util-library",
			path:          "/Users/me/repos/util-library",
			worktrees:     []string{},
			branches:      []string{"main"},
			currentBranch: "main",
			status:        "Clean",
		},
		{
			name:          "api-server",
			path:          "/Users/me/repos/api-server",
			worktrees:     []string{"develop", "feature-auth", "feature-db", "hotfix-security", "main"},
			branches:      []string{"develop", "feature/auth", "feature/db", "hotfix/security", "main"},
			currentBranch: "develop",
			status:        "Clean",
		},
	}

	return RepoSelectModel{
		repos:         repos,
		filteredRepos: repos,
		cursor:        0,
		scrollOffset:  0,
		searchQuery:   "",
		searchFocused: true, // Start with search focused
		theme:         DefaultTheme(),
		width:         0,
		height:        0,
	}
}

// Init initializes the repo select model
func (m RepoSelectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for repo selection
func (m RepoSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScroll()
			}
		case "down", "j":
			if m.cursor < len(m.filteredRepos)-1 {
				m.cursor++
				m.adjustScroll()
			}
		case "enter":
			if len(m.filteredRepos) > 0 {
				m.selected = true
				return m, tea.Quit
			}
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filterRepos()
			}
		default:
			// Add character to search query
			if len(msg.String()) == 1 {
				m.searchQuery += msg.String()
				m.filterRepos()
			}
		}
	}
	return m, nil
}

// filterRepos filters repositories based on search query
func (m *RepoSelectModel) filterRepos() {
	if m.searchQuery == "" {
		m.filteredRepos = m.repos
	} else {
		m.filteredRepos = []mockRepo{}
		query := strings.ToLower(m.searchQuery)
		for _, repo := range m.repos {
			if strings.Contains(strings.ToLower(repo.name), query) {
				m.filteredRepos = append(m.filteredRepos, repo)
			}
		}
	}

	// Reset cursor if out of bounds
	if m.cursor >= len(m.filteredRepos) {
		m.cursor = max(0, len(m.filteredRepos)-1)
	}
	m.scrollOffset = 0
}

// adjustScroll adjusts scroll offset to keep cursor visible
func (m *RepoSelectModel) adjustScroll() {
	visibleHeight := m.getVisibleListHeight()

	// Scroll down if cursor is below visible area
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}

	// Scroll up if cursor is above visible area
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}
}

// getVisibleListHeight returns how many list items can be visible
func (m *RepoSelectModel) getVisibleListHeight() int {
	contentHeight := m.getContentHeight()
	availableHeight := contentHeight - inputHeight - padding
	return max(1, availableHeight)
}

// getContentHeight returns the available height for content
func (m *RepoSelectModel) getContentHeight() int {
	// Use max height or terminal height, whichever is smaller
	effectiveHeight := min(maxHeight, m.height)
	return effectiveHeight - headerHeight - helpHeight - (padding * 2)
}

// getContentWidth returns the available width for content
func (m *RepoSelectModel) getContentWidth() int {
	// Use max width or terminal width, whichever is smaller
	return min(maxWidth, m.width)
}

// View renders the repo selection view
func (m RepoSelectModel) View() string {
	contentWidth := m.getContentWidth()
	contentHeight := m.getContentHeight()

	// Build the main content
	var sections []string

	// Header
	header := m.renderHeader()
	sections = append(sections, header)

	// Search input
	searchInput := m.renderSearchInput(contentWidth)
	sections = append(sections, searchInput)

	// Main content (list + preview)
	mainContent := m.renderMainContent(contentWidth, contentHeight)
	sections = append(sections, mainContent)

	// Help text
	help := m.renderHelp()
	sections = append(sections, help)

	// Join all sections
	innerContent := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Wrap everything in a container box with border
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Padding(1, 2)

	content := containerStyle.Render(innerContent)

	// Center the entire UI
	return m.centerContent(content)
}

// renderHeader renders the header
func (m *RepoSelectModel) renderHeader() string {
	return m.theme.TitleStyle.Render("Select Repository")
}

// renderSearchInput renders the search input field
func (m *RepoSelectModel) renderSearchInput(width int) string {
	icon := m.theme.Icons.Search
	prompt := m.theme.InputPromptStyle.Render(icon + " ")

	// Input content
	queryText := m.searchQuery
	if queryText == "" {
		queryText = "Search repositories..."
	}

	cursor := "█" // Block cursor
	inputContent := queryText + cursor

	// Calculate total width to match panes below
	totalPaneWidth := listPaneWidth + gap + previewPaneWidth

	// Create input with background, left border, and more vertical padding
	inputStyle := lipgloss.NewStyle().
		Width(totalPaneWidth).
		Background(m.theme.InputBgColor).
		Foreground(m.theme.ForegroundColor).
		Padding(1, 2). // Increased vertical padding from 0 to 1
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor).
		MarginBottom(1)

	return inputStyle.Render(prompt + inputContent)
}

// renderMainContent renders the list and preview panes
func (m *RepoSelectModel) renderMainContent(width, height int) string {
	availableHeight := height - inputHeight - padding

	// Render left pane (list) with background
	leftPane := m.renderRepoList(availableHeight)
	leftStyle := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(availableHeight).
		Background(m.theme.ListBgColor).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	// Render right pane (preview) with different background for depth
	rightPane := m.renderPreview(availableHeight)
	rightStyle := lipgloss.NewStyle().
		Width(previewPaneWidth).
		Height(availableHeight).
		Background(m.theme.PreviewBgColor).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	// Join horizontally
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(leftPane),
		strings.Repeat(" ", gap),
		rightStyle.Render(rightPane),
	)
}

// renderRepoList renders the scrollable list of repositories
func (m *RepoSelectModel) renderRepoList(height int) string {
	if len(m.filteredRepos) == 0 {
		return m.theme.LabelStyle.Render("No repositories found")
	}

	var lines []string
	visibleHeight := height - 2 // Account for padding

	// Calculate visible range
	start := m.scrollOffset
	end := min(start+visibleHeight, len(m.filteredRepos))

	// Render visible items
	for i := start; i < end; i++ {
		repo := m.filteredRepos[i]

		// Icon based on repo type
		var icon string
		if len(repo.worktrees) > 0 {
			icon = m.theme.Icons.Worktree
		} else {
			icon = m.theme.Icons.Git
		}

		var line string
		if i == m.cursor {
			// Selected item with background
			cursor := m.theme.Icons.ChevronRight
			content := cursor + " " + icon + " " + repo.name
			line = m.theme.SelectedStyle.Render(content)
		} else {
			// Normal item
			content := "  " + icon + " " + repo.name
			line = m.theme.ItemStyle.Render(content)
		}

		lines = append(lines, line)
	}

	// Add scroll indicator if needed
	if len(m.filteredRepos) > visibleHeight {
		total := len(m.filteredRepos)
		showing := fmt.Sprintf("  (%d/%d)", end, total)
		lines = append(lines, "")
		lines = append(lines, m.theme.HelpStyle.Render(showing))
	}

	return strings.Join(lines, "\n")
}

// renderPreview renders the preview pane for the selected repo
func (m *RepoSelectModel) renderPreview(height int) string {
	if len(m.filteredRepos) == 0 {
		return ""
	}

	if m.cursor >= len(m.filteredRepos) {
		return ""
	}

	repo := m.filteredRepos[m.cursor]
	var lines []string

	// Name with folder icon
	nameLabel := m.theme.LabelStyle.Render(m.theme.Icons.Folder + " ")
	lines = append(lines, nameLabel+m.theme.ValueStyle.Render(repo.name))

	// Path
	pathStyle := lipgloss.NewStyle().Foreground(m.theme.MutedColor)
	lines = append(lines, pathStyle.Render("  "+repo.path))
	lines = append(lines, "")

	// Show EITHER worktrees OR branches, not both
	if len(repo.worktrees) > 0 {
		// Worktrees section (only if repo has worktrees)
		worktreeIcon := m.theme.Icons.Worktree
		lines = append(lines, m.theme.LabelStyle.Render(worktreeIcon+" Worktrees"))
		lines = append(lines, "")
		for _, worktree := range repo.worktrees {
			line := "  " + m.theme.Icons.ChevronRight + " " + worktree
			lines = append(lines, m.theme.ValueStyle.Render(line))
		}
		lines = append(lines, "")
	} else {
		// Branches section (only if repo has no worktrees)
		branchIcon := m.theme.Icons.Branch
		lines = append(lines, m.theme.LabelStyle.Render(branchIcon+" Branches"))
		lines = append(lines, "")

		// Sort branches alphabetically
		sortedBranches := make([]string, len(repo.branches))
		copy(sortedBranches, repo.branches)
		sort.Strings(sortedBranches)

		for _, branch := range sortedBranches {
			if branch == repo.currentBranch {
				// Highlight current branch with check icon
				line := "  " + m.theme.Icons.Check + " " + branch
				lines = append(lines, m.theme.SelectedStyle.Render(line))
			} else {
				line := "  " + m.theme.Icons.ChevronRight + " " + branch
				lines = append(lines, m.theme.ValueStyle.Render(line))
			}
		}
		lines = append(lines, "")
	}

	// Status with icon
	var statusIcon string
	var statusStyle lipgloss.Style
	if strings.Contains(repo.status, "Modified") {
		statusIcon = m.theme.Icons.Modified
		statusStyle = lipgloss.NewStyle().Foreground(m.theme.WarningColor)
	} else {
		statusIcon = m.theme.Icons.Check
		statusStyle = lipgloss.NewStyle().Foreground(m.theme.SuccessColor)
	}
	statusLabel := m.theme.LabelStyle.Render(statusIcon + " ")
	lines = append(lines, statusLabel+statusStyle.Render(repo.status))

	return strings.Join(lines, "\n")
}

// renderHelp renders the help text
func (m *RepoSelectModel) renderHelp() string {
	helpText := "↑/k up • ↓/j down • enter select • q/ctrl+c quit"
	return m.theme.HelpStyle.Render(helpText)
}

// centerContent centers the content in the terminal
func (m *RepoSelectModel) centerContent(content string) string {
	if m.width == 0 || m.height == 0 {
		return content
	}

	contentWidth := m.getContentWidth()
	contentHeight := m.getContentHeight()

	// Horizontal centering
	horizontalPadding := max(0, (m.width-contentWidth)/2)

	// Vertical centering
	verticalPadding := max(0, (m.height-contentHeight-4)/2)

	// Apply centering
	centered := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(verticalPadding, horizontalPadding)

	return centered.Render(content)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
