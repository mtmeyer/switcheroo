package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RepoSelectModel handles repository selection
type RepoSelectModel struct {
	repos         []RepoDisplay
	filteredRepos []RepoDisplay
	cursor        int
	scrollOffset  int
	searchQuery   string
	theme         Theme
	width         int
	height        int
}

// NewRepoSelectModel creates a new repository selection model
func NewRepoSelectModel(repos []RepoDisplay, theme Theme) RepoSelectModel {
	return RepoSelectModel{
		repos:         repos,
		filteredRepos: repos,
		cursor:        0,
		scrollOffset:  0,
		searchQuery:   "",
		theme:         theme,
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
			}
		case "down", "j":
			if m.cursor < len(m.filteredRepos)-1 {
				m.cursor++
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

// renderTitle renders the header
func (m *RepoSelectModel) renderTitle() string {
	return m.theme.TitleStyle.Render("Select Repository")
}

// renderSearch renders the search input
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
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	if hasSize && layout.contentWidth > 0 {
		inputStyle = inputStyle.Width(layout.contentWidth)
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
	leftContent := m.renderList()
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
	rightContent := m.renderPreview()
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
func (m *RepoSelectModel) renderList() string {
	if len(m.filteredRepos) == 0 {
		return m.theme.LabelStyle.Render("No repositories found")
	}

	visibleHeight := paneHeight - 2 // Account for padding
	var lines []string

	// Show repos that fit in the visible area
	for i := 0; i < len(m.filteredRepos) && i < visibleHeight; i++ {
		repo := m.filteredRepos[i]

		var icon string
		if repo.HasWorktrees {
			icon = m.theme.Icons.Worktree
		} else {
			icon = m.theme.Icons.Git
		}

		var line string
		if i == m.cursor {
			cursor := m.theme.Icons.ChevronRight
			content := cursor + " " + icon + " " + repo.Name
			line = m.theme.SelectedStyle.Render(content)
		} else {
			content := "  " + icon + " " + repo.Name
			line = m.theme.ItemStyle.Render(content)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// renderPreview renders the preview pane
func (m *RepoSelectModel) renderPreview() string {
	if len(m.filteredRepos) == 0 || m.cursor >= len(m.filteredRepos) {
		return ""
	}

	repo := m.filteredRepos[m.cursor]
	var lines []string

	// Name with folder icon
	nameLabel := m.theme.LabelStyle.Render(m.theme.Icons.Folder + " ")
	lines = append(lines, nameLabel+m.theme.ValueStyle.Render(repo.Name))

	// Path
	pathStyle := lipgloss.NewStyle().Foreground(m.theme.MutedColor)
	lines = append(lines, pathStyle.Render("  "+repo.Path))
	lines = append(lines, "")

	// Show worktrees OR branches
	if len(repo.Worktrees) > 0 {
		worktreeIcon := m.theme.Icons.Worktree
		lines = append(lines, m.theme.LabelStyle.Render(worktreeIcon+" Worktrees"))
		lines = append(lines, "")

		for _, worktree := range repo.Worktrees {
			line := "  " + m.theme.Icons.ChevronRight + " " + worktree.Name
			lines = append(lines, m.theme.ValueStyle.Render(line))
		}
	} else {
		branchIcon := m.theme.Icons.Branch
		lines = append(lines, m.theme.LabelStyle.Render(branchIcon+" Branches"))
		lines = append(lines, "")

		for _, branch := range repo.Branches {
			if branch == repo.CurrentBranch {
				line := "  " + m.theme.Icons.Check + " " + branch
				lines = append(lines, m.theme.SelectedStyle.Render(line))
			} else {
				line := "  " + m.theme.Icons.ChevronRight + " " + branch
				lines = append(lines, m.theme.ValueStyle.Render(line))
			}
		}
	}

	return strings.Join(lines, "\n")
}

// renderHelp renders the help text
func (m *RepoSelectModel) renderHelp() string {
	helpText := "↑/k up • ↓/j down • enter select • q/ctrl+c quit"
	return m.theme.HelpStyle.Render(helpText)
}
