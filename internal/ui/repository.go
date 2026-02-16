package ui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"switcheroo/internal/git"
)

// RepoDisplay is an adapter between git.Repository and UI needs
type RepoDisplay struct {
	Name            string
	Path            string
	Worktrees       []WorktreeDisplay
	Branches        []BranchDisplay
	CurrentBranch   string
	HasWorktrees    bool
	LineDiffAdded   int
	LineDiffRemoved int
	DefaultBranch   string
}

type BranchDisplay struct {
	Name        string
	IsCurrent   bool
	Status      string
	Upstream    string
	DiffAdded   int
	DiffRemoved int
}

// WorktreeDisplay represents worktree data needed by the UI
type WorktreeDisplay struct {
	Name        string
	Path        string
	Branch      string
	Locked      bool
	Status      string
	DiffAdded   int
	DiffRemoved int
}

// FromGitRepositories converts git.Repository slice to RepoDisplay slice
func FromGitRepositories(repos []git.Repository) []RepoDisplay {
	displays := make([]RepoDisplay, len(repos))

	for i, repo := range repos {
		displays[i] = RepoDisplay{
			Name:          repo.Name,
			Path:          repo.Path,
			HasWorktrees:  repo.HasWorktrees,
			DefaultBranch: repo.DefaultBranch,
		}

		if repo.HasWorktrees {
			worktrees := make([]WorktreeDisplay, len(repo.Worktrees))
			for j, wt := range repo.Worktrees {
				worktrees[j] = WorktreeDisplay{
					Name:        wt.Name,
					Path:        wt.Path,
					Branch:      wt.Branch,
					Locked:      wt.IsLocked,
					Status:      wt.Status,
					DiffAdded:   wt.DiffAdded,
					DiffRemoved: wt.DiffRemoved,
				}
			}
			displays[i].Worktrees = worktrees
		} else {
			branches := make([]BranchDisplay, len(repo.Branches))
			for j, branch := range repo.Branches {
				branches[j] = BranchDisplay{
					Name:        branch.Name,
					IsCurrent:   branch.IsCurrent,
					Status:      branch.Status,
					Upstream:    branch.Upstream,
					DiffAdded:   branch.DiffAdded,
					DiffRemoved: branch.DiffRemoved,
				}
				if branch.IsCurrent {
					displays[i].CurrentBranch = branch.Name
				}
			}
			sort.Slice(branches, func(a, b int) bool {
				return strings.ToLower(branches[a].Name) < strings.ToLower(branches[b].Name)
			})
			displays[i].Branches = branches
		}
	}

	return displays
}

// Selectable interface implementation for RepoDisplay

func (r RepoDisplay) GetID() string {
	return r.Path
}

func (r RepoDisplay) GetName() string {
	return r.Name
}

func (r RepoDisplay) GetPath() string {
	return r.Path
}

func (r RepoDisplay) GetBranch() string {
	return r.CurrentBranch
}

func (r RepoDisplay) HasChildren() bool {
	return r.HasWorktrees
}

func (r RepoDisplay) OnSelect() string {
	if r.HasWorktrees {
		return "" // Has children, need to drill in
	}
	return r.Path
}

func (r RepoDisplay) RenderHeader(theme Theme) string {
	nameLabel := theme.LabelStyle.Render(theme.Icons.Folder + " ")
	maxNameWidth := 40 // reasonable default
	truncatedName := truncateWithEllipsis(r.Name, maxNameWidth)
	return nameLabel + theme.ValueStyle.Render(truncatedName)
}

func (r RepoDisplay) RenderChildren(theme Theme, settings PreviewSettings) []string {
	var lines []string

	if r.HasWorktrees && len(r.Worktrees) > 0 {
		worktreeIcon := theme.Icons.Worktree
		lines = append(lines, theme.LabelStyle.Render(worktreeIcon+" Worktrees"))
		lines = append(lines, "")

		limit := len(r.Worktrees)
		if limit > maxPreviewListItems {
			limit = maxPreviewListItems
		}
		for i := 0; i < limit; i++ {
			worktree := r.Worktrees[i]
			line := "  " + theme.Icons.ChevronRight + " " + worktree.Name
			lines = append(lines, theme.ValueStyle.Render(line))
		}
		if len(r.Worktrees) > limit {
			remaining := len(r.Worktrees) - limit
			moreLine := "  ... " + fmt.Sprintf("%d more worktrees", remaining)
			lines = append(lines, theme.HelpStyle.Render(moreLine))
		}
	} else if len(r.Branches) > 0 {
		branchIcon := theme.Icons.Branch
		lines = append(lines, theme.LabelStyle.Render(branchIcon+" Branches"))
		lines = append(lines, "")

		limit := len(r.Branches)
		if limit > maxPreviewListItems {
			limit = maxPreviewListItems
		}
		for i := 0; i < limit; i++ {
			branch := r.Branches[i]
			line := "  "
			if branch.IsCurrent {
				line += theme.Icons.Check
			} else {
				line += theme.Icons.ChevronRight
			}
			line += " " + branch.Name
			var extras []string
			if settings.Repo.ShowBranchStatus {
				extras = append(extras, formatStatusInline(theme, branch.Status))
			}
			if settings.Repo.ShowBranchDiff && branch.Status != "untracked" {
				extras = append(extras, formatDiffText(theme, branch.DiffAdded, branch.DiffRemoved))
			}
			if len(extras) > 0 {
				line += " " + strings.Join(extras, " ")
			}
			if branch.IsCurrent {
				lines = append(lines, theme.SelectedStyle.Render(line))
			} else {
				lines = append(lines, theme.ValueStyle.Render(line))
			}
		}
		if len(r.Branches) > limit {
			remaining := len(r.Branches) - limit
			moreLine := "  ... " + fmt.Sprintf("%d more branches", remaining)
			lines = append(lines, theme.HelpStyle.Render(moreLine))
		}
	}

	return lines
}

func (r RepoDisplay) RenderMetadata(theme Theme, settings PreviewSettings) []string {
	var lines []string

	// Path line
	pathStyle := lipgloss.NewStyle().Foreground(theme.MutedColor)
	lines = append(lines, pathStyle.Render("  "+r.Path))

	if settings.Repo.ShowLineDiff {
		lines = append(lines, "")
		lines = append(lines, theme.LabelStyle.Render(theme.Icons.Modified+" Line Diff"))
		diffLine := "  " + formatDiffText(theme, r.LineDiffAdded, r.LineDiffRemoved)
		lines = append(lines, diffLine)
	}

	return lines
}

// Selectable interface implementation for WorktreeDisplay

func (w WorktreeDisplay) GetID() string {
	return w.Path
}

func (w WorktreeDisplay) GetName() string {
	return w.Name
}

func (w WorktreeDisplay) GetPath() string {
	return w.Path
}

func (w WorktreeDisplay) GetBranch() string {
	return w.Branch
}

func (w WorktreeDisplay) HasChildren() bool {
	return false
}

func (w WorktreeDisplay) OnSelect() string {
	return w.Path
}

func (w WorktreeDisplay) RenderHeader(theme Theme) string {
	nameLabel := theme.LabelStyle.Render(theme.Icons.Worktree + " ")
	maxNameWidth := 40
	truncatedName := truncateWithEllipsis(w.Name, maxNameWidth)
	return nameLabel + theme.ValueStyle.Render(truncatedName)
}

func (w WorktreeDisplay) RenderChildren(theme Theme, settings PreviewSettings) []string {
	return nil // Worktrees don't have children
}

func (w WorktreeDisplay) RenderMetadata(theme Theme, settings PreviewSettings) []string {
	var lines []string

	if w.Branch != "" {
		branchLine := theme.LabelStyle.Render(theme.Icons.Branch+" Branch ") + theme.ValueStyle.Render(w.Branch)
		lines = append(lines, branchLine)
	}

	pathStyle := lipgloss.NewStyle().Foreground(theme.MutedColor)
	lines = append(lines, "")
	lines = append(lines, pathStyle.Render(w.Path))

	if settings.Worktree.ShowStatus {
		status := formatStatusValue(theme, w.Status)
		line := theme.LabelStyle.Render(theme.Icons.Warning+" Status ") + " " + status
		lines = append(lines, "")
		lines = append(lines, line)
	}

	if settings.Worktree.ShowLineDiff {
		lines = append(lines, theme.LabelStyle.Render(theme.Icons.Modified+" Line Diff"))
		diff := "  " + formatDiffText(theme, w.DiffAdded, w.DiffRemoved)
		lines = append(lines, diff)
	}

	if w.Locked {
		warningStyle := lipgloss.NewStyle()
		if theme.WarningColor != "" {
			warningStyle = warningStyle.Foreground(theme.WarningColor)
		}
		lines = append(lines, warningStyle.Render(theme.Icons.Warning+" Locked"))
	}

	return lines
}

// LoadRepoMetadata loads expensive metadata for a repo when selected
func LoadRepoMetadata(repo RepoDisplay) tea.Cmd {
	return func() tea.Msg {
		msg := RepoMetadataLoadedMsg{
			RepoPath: repo.Path,
		}

		if repo.HasWorktrees {
			// Load worktree metadata
			worktrees := make([]WorktreeDisplay, len(repo.Worktrees))
			for i, wt := range repo.Worktrees {
				worktrees[i] = wt
				// Load status
				if status, err := git.GetWorktreeStatus(wt.Path); err == nil {
					worktrees[i].Status = status
				}
				// Load diff
				if added, removed, err := git.GetWorktreeDiff(wt.Path); err == nil {
					worktrees[i].DiffAdded = added
					worktrees[i].DiffRemoved = removed
				}
			}
			msg.Worktrees = worktrees
		} else {
			// Load branch diffs
			branches := make([]BranchDisplay, len(repo.Branches))
			for i, branch := range repo.Branches {
				branches[i] = branch
				if branch.Upstream != "" {
					if added, removed, err := git.DiffAgainstUpstream(repo.Path, branch.Name, branch.Upstream); err == nil {
						branches[i].DiffAdded = added
						branches[i].DiffRemoved = removed
					}
				}
			}
			msg.Branches = branches

			// Load line diff for current branch
			for _, branch := range branches {
				if branch.IsCurrent {
					if branch.Status == "untracked" && repo.DefaultBranch != "" && repo.DefaultBranch != branch.Name {
						if added, removed, err := git.DiffBranches(repo.Path, branch.Name, repo.DefaultBranch); err == nil {
							msg.LineDiffAdded = added
							msg.LineDiffRemoved = removed
						}
					} else {
						msg.LineDiffAdded = branch.DiffAdded
						msg.LineDiffRemoved = branch.DiffRemoved
					}
					break
				}
			}
		}

		return msg
	}
}
