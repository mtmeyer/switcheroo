package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxPreviewListItems = 8

// RepoSelectModel handles repository selection
type RepoSelectModel struct {
	repos         []RepoDisplay
	filteredRepos []RepoDisplay
	cursor        int
	scrollOffset  int
	searchQuery   string
	theme         Theme
	settings      PreviewSettings
	width         int
	height        int
}

// NewRepoSelectModel creates a new repository selection model
func NewRepoSelectModel(repos []RepoDisplay, theme Theme, settings PreviewSettings) RepoSelectModel {
	return RepoSelectModel{
		repos:         repos,
		filteredRepos: repos,
		cursor:        0,
		scrollOffset:  0,
		searchQuery:   "",
		theme:         theme,
		settings:      settings,
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

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.adjustScrollOffset()
			}
		case "down", "j":
			if m.cursor < len(m.filteredRepos)-1 {
				m.cursor++
				m.adjustScrollOffset()
			}
		case "enter":
			if len(m.filteredRepos) > 0 {
				selected := m.filteredRepos[m.cursor]
				return m, repoSelectedCmd(selected)
			}
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filterRepos()
			}
		default:
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
		m.filteredRepos = []RepoDisplay{}
		query := strings.ToLower(m.searchQuery)
		for _, repo := range m.repos {
			if strings.Contains(strings.ToLower(repo.Name), query) {
				m.filteredRepos = append(m.filteredRepos, repo)
			}
		}
	}
	if m.cursor >= len(m.filteredRepos) {
		m.cursor = 0
		m.scrollOffset = 0
	}
}

// adjustScrollOffset ensures the cursor stays visible by adjusting scroll offset
func (m *RepoSelectModel) adjustScrollOffset() {
	visibleHeight := paneHeight - 2

	// If cursor is above the visible area, scroll up
	if m.cursor < m.scrollOffset {
		m.scrollOffset = m.cursor
	}

	// If cursor is below the visible area, scroll down
	if m.cursor >= m.scrollOffset+visibleHeight {
		m.scrollOffset = m.cursor - visibleHeight + 1
	}
}

// View renders the repo selection view
func (m RepoSelectModel) View() string {
	layout := calculateLayout(m.width)
	hasSize := m.width > 0 && m.height > 0

	// Build components
	title := m.renderTitle()
	search := m.renderSearch(layout, hasSize)
	content := m.renderContent(layout)
	help := m.renderHelp()

	// Combine all sections
	sections := []string{title, search, content, help}
	inner := strings.Join(sections, "\n")

	// Wrap in container with border
	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Padding(1, containerPadding)

	if hasSize {
		containerStyle = containerStyle.Width(layout.containerWidth)
	}

	// Only set background if theme defines one
	if m.theme.BackgroundColor != lipgloss.Color("") {
		containerStyle = containerStyle.Background(m.theme.BackgroundColor)
	}

	contentBox := containerStyle.Render(inner)

	if !hasSize {
		return contentBox
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		contentBox,
	)
}

func (m *RepoSelectModel) renderTitle() string {
	return m.theme.TitleStyle.Render("Select Repository")
}

func (m *RepoSelectModel) renderSearch(layout layoutMetrics, hasSize bool) string {
	icon := m.theme.Icons.Search
	prompt := m.theme.InputPromptStyle.Render(icon + " ")

	queryText := m.searchQuery
	if queryText == "" {
		queryText = "Search repositories..."
	}

	cursor := "█"
	inputContent := prompt + queryText + cursor

	inputStyle := lipgloss.NewStyle().
		Padding(1, 2).
		Margin(0, 0, 1, 0).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	if hasSize && layout.contentWidth > 0 {
		inputStyle = inputStyle.Width(layout.contentWidth - 1)
	}

	// Only set background if theme defines one
	if m.theme.InputBgColor != lipgloss.Color("") {
		inputStyle = inputStyle.Background(m.theme.InputBgColor)
	}

	return inputStyle.Render(inputContent)
}

// renderContent renders the list and preview side by side
func (m *RepoSelectModel) renderContent(layout layoutMetrics) string {
	listPaneWidth := layout.listWidth
	previewPaneWidth := layout.previewWidth

	if listPaneWidth < 1 {
		listPaneWidth = 1
	}
	if previewPaneWidth < 1 {
		previewPaneWidth = 1
	}
	// Render left pane (list)
	leftContent := m.renderList(layout)
	leftStyle := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	// Only set background if theme defines one
	if m.theme.ListBgColor != lipgloss.Color("") {
		leftStyle = leftStyle.Background(m.theme.ListBgColor)
	}
	leftPane := leftStyle.Render(leftContent)

	// Render right pane (preview)
	rightContent := m.renderPreview(layout)
	rightStyle := lipgloss.NewStyle().
		Width(previewPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	// Only set background if theme defines one
	if m.theme.PreviewBgColor != lipgloss.Color("") {
		rightStyle = rightStyle.Background(m.theme.PreviewBgColor)
	}
	rightPane := rightStyle.Render(rightContent)

	// Join horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, strings.Repeat(" ", gapWidth), rightPane)
}

// renderList renders the repository list
func (m *RepoSelectModel) renderList(layout layoutMetrics) string {
	if len(m.filteredRepos) == 0 {
		return m.theme.LabelStyle.Render("No repositories found")
	}

	visibleHeight := paneHeight - 2 // Account for padding
	var lines []string

	// Calculate max width for repo names (accounting for padding, cursor, icon)
	maxNameWidth := layout.listWidth - 9 // 4 padding + 1 cursor + 1 icon + 3 spaces
	if maxNameWidth < 10 {
		maxNameWidth = 10
	}

	// Show repos that fit in the visible area, accounting for scroll offset
	startIdx := m.scrollOffset
	endIdx := startIdx + visibleHeight
	if endIdx > len(m.filteredRepos) {
		endIdx = len(m.filteredRepos)
	}

	for i := startIdx; i < endIdx; i++ {
		repo := m.filteredRepos[i]

		var icon string
		if repo.HasWorktrees {
			icon = m.theme.Icons.Worktree
		} else {
			icon = m.theme.Icons.Git
		}

		// Truncate repo name with ellipsis
		truncatedName := truncateWithEllipsis(repo.Name, maxNameWidth)

		var line string
		if i == m.cursor {
			cursor := m.theme.Icons.ChevronRight
			content := cursor + " " + icon + " " + truncatedName
			line = m.theme.SelectedStyle.Render(content)
		} else {
			content := "  " + icon + " " + truncatedName
			line = m.theme.ItemStyle.Render(content)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderPreview renders the preview pane
func (m *RepoSelectModel) renderPreview(layout layoutMetrics) string {
	if len(m.filteredRepos) == 0 || m.cursor >= len(m.filteredRepos) {
		return ""
	}

	repo := m.filteredRepos[m.cursor]
	var lines []string

	// Calculate max width for repo name in preview
	maxNameWidth := layout.previewWidth - 6 // 4 padding + 1 icon + 1 space
	if maxNameWidth < 20 {
		maxNameWidth = 20
	}

	// Name with folder icon (truncated)
	nameLabel := m.theme.LabelStyle.Render(m.theme.Icons.Folder + " ")
	truncatedName := truncateWithEllipsis(repo.Name, maxNameWidth)
	lines = append(lines, nameLabel+m.theme.ValueStyle.Render(truncatedName))

	// Path
	pathStyle := lipgloss.NewStyle().Foreground(m.theme.MutedColor)
	lines = append(lines, pathStyle.Render("  "+repo.Path))
	lines = append(lines, "")

	// Show worktrees OR branches
	if len(repo.Worktrees) > 0 {
		worktreeIcon := m.theme.Icons.Worktree
		lines = append(lines, m.theme.LabelStyle.Render(worktreeIcon+" Worktrees"))
		lines = append(lines, "")

		limit := len(repo.Worktrees)
		if limit > maxPreviewListItems {
			limit = maxPreviewListItems
		}
		for i := 0; i < limit; i++ {
			worktree := repo.Worktrees[i]
			line := "  " + m.theme.Icons.ChevronRight + " " + worktree.Name
			lines = append(lines, m.theme.ValueStyle.Render(line))
		}
		if len(repo.Worktrees) > limit {
			remaining := len(repo.Worktrees) - limit
			moreLine := fmt.Sprintf("  ... %d more worktrees", remaining)
			lines = append(lines, m.theme.HelpStyle.Render(moreLine))
		}
	} else {
		branchIcon := m.theme.Icons.Branch
		lines = append(lines, m.theme.LabelStyle.Render(branchIcon+" Branches"))
		lines = append(lines, "")

		limit := len(repo.Branches)
		if limit > maxPreviewListItems {
			limit = maxPreviewListItems
		}
		for i := 0; i < limit; i++ {
			branch := repo.Branches[i]
			line := "  "
			if branch.IsCurrent {
				line += m.theme.Icons.Check
			} else {
				line += m.theme.Icons.ChevronRight
			}
			line += " " + branch.Name
			var extras []string
			if m.settings.Repo.ShowBranchStatus {
				extras = append(extras, formatStatusInline(m.theme, branch.Status))
			}
			if m.settings.Repo.ShowBranchDiff && branch.Status != "untracked" {
				extras = append(extras, formatDiffText(m.theme, branch.DiffAdded, branch.DiffRemoved))
			}
			if len(extras) > 0 {
				line += " " + strings.Join(extras, " ")
			}
			if branch.IsCurrent {
				lines = append(lines, m.theme.SelectedStyle.Render(line))
			} else {
				lines = append(lines, m.theme.ValueStyle.Render(line))
			}
		}
		if len(repo.Branches) > limit {
			remaining := len(repo.Branches) - limit
			moreLine := fmt.Sprintf("  ... %d more branches", remaining)
			lines = append(lines, m.theme.HelpStyle.Render(moreLine))
		}

		if m.settings.Repo.ShowLineDiff {
			lines = append(lines, "")
			lines = append(lines, m.theme.LabelStyle.Render(m.theme.Icons.Modified+" Line Diff"))
			diffLine := "  " + formatDiffText(m.theme, repo.LineDiffAdded, repo.LineDiffRemoved)
			lines = append(lines, diffLine)
		}
	}

	return strings.Join(lines, "\n")
}

// renderHelp renders the help text
func (m *RepoSelectModel) renderHelp() string {
	helpText := "↑/k up • ↓/j down • enter select • q/ctrl+c quit"
	return m.theme.HelpStyle.Render(helpText)
}
