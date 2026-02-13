package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WorktreeSelectModel handles worktree selection for a repo
type WorktreeSelectModel struct {
	repoName          string
	worktrees         []WorktreeDisplay
	filteredWorktrees []WorktreeDisplay
	cursor            int
	searchQuery       string
	theme             Theme
	width             int
	height            int
}

// NewWorktreeSelectModel creates a new model for a repo's worktrees
func NewWorktreeSelectModel(repo RepoDisplay, theme Theme) WorktreeSelectModel {
	return WorktreeSelectModel{
		repoName:          repo.Name,
		worktrees:         repo.Worktrees,
		filteredWorktrees: repo.Worktrees,
		theme:             theme,
	}
}

func (m WorktreeSelectModel) Init() tea.Cmd {
	return nil
}

func (m WorktreeSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.filteredWorktrees)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.filteredWorktrees) > 0 {
				selected := m.filteredWorktrees[m.cursor]
				return m, worktreeSelectedCmd(selected)
			}
		case "esc":
			return m, backToReposCmd()
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filterWorktrees()
			}
		default:
			if len(msg.String()) == 1 {
				m.searchQuery += msg.String()
				m.filterWorktrees()
			}
		}
	}
	return m, nil
}

func (m *WorktreeSelectModel) filterWorktrees() {
	if m.searchQuery == "" {
		m.filteredWorktrees = m.worktrees
	} else {
		var filtered []WorktreeDisplay
		query := strings.ToLower(m.searchQuery)
		for _, wt := range m.worktrees {
			if strings.Contains(strings.ToLower(wt.Name), query) || strings.Contains(strings.ToLower(wt.Branch), query) {
				filtered = append(filtered, wt)
			}
		}
		m.filteredWorktrees = filtered
	}
	if m.cursor >= len(m.filteredWorktrees) {
		m.cursor = 0
	}
}

func (m WorktreeSelectModel) View() string {
	layout := calculateLayout(m.width)
	hasSize := m.width > 0 && m.height > 0

	title := m.renderTitle()
	search := m.renderSearch(layout, hasSize)
	content := m.renderContent(layout)
	help := m.renderHelp()

	sections := []string{title, search, content, help}
	inner := strings.Join(sections, "\n")

	containerStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.BorderColor).
		Padding(1, containerPadding)

	if hasSize {
		containerStyle = containerStyle.Width(layout.containerWidth)
	}

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

func (m WorktreeSelectModel) renderTitle() string {
	return m.theme.TitleStyle.Render("Select Worktree: " + m.repoName)
}

func (m WorktreeSelectModel) renderSearch(layout layoutMetrics, hasSize bool) string {
	queryText := m.searchQuery
	if queryText == "" {
		queryText = "Search worktrees..."
	}

	cursor := "█"
	inputContent := m.theme.InputPromptStyle.Render(m.theme.Icons.Search+" ") + queryText + cursor

	inputStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	if hasSize && layout.contentWidth > 0 {
		inputStyle = inputStyle.Width(layout.contentWidth)
	}

	if m.theme.InputBgColor != lipgloss.Color("") {
		inputStyle = inputStyle.Background(m.theme.InputBgColor)
	}

	return inputStyle.Render(inputContent)
}

func (m WorktreeSelectModel) renderContent(layout layoutMetrics) string {
	listPaneWidth := layout.listWidth
	previewPaneWidth := layout.previewWidth

	if listPaneWidth < 1 {
		listPaneWidth = 1
	}
	if previewPaneWidth < 1 {
		previewPaneWidth = 1
	}

	leftContent := m.renderList()
	leftStyle := lipgloss.NewStyle().
		Width(listPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	if m.theme.ListBgColor != lipgloss.Color("") {
		leftStyle = leftStyle.Background(m.theme.ListBgColor)
	}

	leftPane := leftStyle.Render(leftContent)

	rightContent := m.renderPreview()
	rightStyle := lipgloss.NewStyle().
		Width(previewPaneWidth).
		Height(paneHeight).
		Padding(1, 2).
		BorderLeft(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(m.theme.BorderColor)

	if m.theme.PreviewBgColor != lipgloss.Color("") {
		rightStyle = rightStyle.Background(m.theme.PreviewBgColor)
	}

	rightPane := rightStyle.Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, strings.Repeat(" ", gapWidth), rightPane)
}

func (m WorktreeSelectModel) renderList() string {
	if len(m.filteredWorktrees) == 0 {
		return m.theme.LabelStyle.Render("No worktrees found")
	}

	visibleHeight := paneHeight - 2
	var lines []string

	for i := 0; i < len(m.filteredWorktrees) && i < visibleHeight; i++ {
		wt := m.filteredWorktrees[i]
		icon := m.theme.Icons.Worktree
		lineContent := icon + " " + wt.Name
		if wt.Branch != "" {
			lineContent += " (" + wt.Branch + ")"
		}

		if i == m.cursor {
			cursorSymbol := m.theme.Icons.ChevronRight
			line := cursorSymbol + " " + lineContent
			lines = append(lines, m.theme.SelectedStyle.Render(line))
		} else {
			line := "  " + lineContent
			lines = append(lines, m.theme.ItemStyle.Render(line))
		}
	}

	return strings.Join(lines, "\n")
}

func (m WorktreeSelectModel) renderPreview() string {
	if len(m.filteredWorktrees) == 0 || m.cursor >= len(m.filteredWorktrees) {
		return ""
	}

	wt := m.filteredWorktrees[m.cursor]
	var lines []string

	nameLine := m.theme.LabelStyle.Render(m.theme.Icons.Worktree+" ") + m.theme.ValueStyle.Render(wt.Name)
	lines = append(lines, nameLine)

	if wt.Branch != "" {
		branchLine := m.theme.LabelStyle.Render(m.theme.Icons.Branch+" Branch ") + m.theme.ValueStyle.Render(wt.Branch)
		lines = append(lines, branchLine)
	}

	if wt.Path != "" {
		pathStyle := lipgloss.NewStyle().Foreground(m.theme.MutedColor)
		lines = append(lines, "")
		lines = append(lines, pathStyle.Render(wt.Path))
	}

	if wt.Locked {
		warningStyle := lipgloss.NewStyle()
		if m.theme.WarningColor != "" {
			warningStyle = warningStyle.Foreground(m.theme.WarningColor)
		}
		lines = append(lines, warningStyle.Render(m.theme.Icons.Warning+" Locked"))
	}

	return strings.Join(lines, "\n")
}

func (m WorktreeSelectModel) renderHelp() string {
	return m.theme.HelpStyle.Render("↑/k up • ↓/j down • enter select • esc back • q/ctrl+c quit")
}
