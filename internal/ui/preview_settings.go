package ui

import "strings"

// PreviewSettings stores which preview fields are enabled
type PreviewSettings struct {
	Repo struct {
		ShowBranchStatus bool
		ShowBranchDiff   bool
		ShowLineDiff     bool
	}
	Worktree struct {
		ShowStatus   bool
		ShowLineDiff bool
	}
}

// NewPreviewSettings constructs settings from config lists
func NewPreviewSettings(repoFields, worktreeFields []string) PreviewSettings {
	var settings PreviewSettings
	for _, field := range repoFields {
		switch strings.ToLower(field) {
		case "branch_status":
			settings.Repo.ShowBranchStatus = true
		case "branch_diff":
			settings.Repo.ShowBranchDiff = true
		case "line_diff":
			settings.Repo.ShowLineDiff = true
		}
	}
	for _, field := range worktreeFields {
		switch strings.ToLower(field) {
		case "status":
			settings.Worktree.ShowStatus = true
		case "line_diff":
			settings.Worktree.ShowLineDiff = true
		}
	}
	return settings
}
