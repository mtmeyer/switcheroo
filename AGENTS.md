# AGENTS.md

This document provides guidance for AI assistants working on the Switcheroo project.

## Project Overview

**Switcheroo** is a terminal-based Git repository and worktree switcher built with Go and Bubble Tea. It provides a fuzzy-findable interface for navigating repositories, with special support for git worktrees.

## Architecture

### Technology Stack
- **Language**: Go 1.22+
- **TUI Framework**: Bubble Tea (Elm Architecture)
- **Styling**: Lipgloss
- **Icons**: Nerd Fonts (unicode escapes)
- **Configuration**: JSON

### Project Structure
```
switcheroo/
├── cmd/switcheroo/          # Main entry point
│   └── main.go
├── internal/
│   ├── config/              # Configuration loading
│   │   ├── config.go
│   │   └── defaults.go
│   ├── git/                 # Git operations
│   │   ├── repository.go    # Repo discovery
│   │   ├── worktree.go      # Worktree types
│   │   ├── branch.go        # Branch types
│   │   ├── metadata.go      # Metadata types
│   │   └── git.go           # Git commands
│   ├── ui/                  # User interface
│   │   ├── model.go         # Bubble Tea model
│   │   ├── repo_select.go   # Repository selection view
│   │   ├── repo_display.go  # Display adapter
│   │   ├── theme.go         # Styling and colors
│   │   └── worktree_select.go # Worktree selection (future)
│   └── plugin/              # Plugin system (future)
│       ├── interface.go
│       └── hooks.go
├── DESIGN.md                # Comprehensive design document
└── README.md                # User documentation
```

## Design Principles

### UI Design
1. **No full borders** - Use left borders and background colors for depth
2. **Background colors** create visual hierarchy:
   - Input: Darkest (`#234`)
   - List: Medium (`#235`)
   - Preview: Lightest (`#236`)
3. **Left borders** in accent color (cyan) for structure
4. **Nerd Font icons** for visual distinction
5. **Fixed heights** prevent layout shifts
6. **Scrolling** for overflow content in both list and preview

### Git Logic
1. **Worktree detection** uses `git worktree list --porcelain`
2. **Filter out main worktree** - only show additional worktrees
3. **Worktrees XOR branches** - never show both in preview
4. **Alphabetical sorting** for branches
5. **Current branch highlighting** with check icon

## Coding Conventions

### Go Style
- Use standard Go formatting (gofmt)
- Error handling: check errors explicitly, don't ignore
- Use `filepath.Join()` for cross-platform paths
- Prefer `strings.Contains()` over regex for simple matching

### Bubble Tea Patterns
```go
// Model should hold all state
type Model struct {
    cursor        int
    scrollOffset  int
    previewScroll int // For preview pane scrolling
    // ... other fields
}

// Update handles all message types
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle key presses
    case tea.WindowSizeMsg:
        // Handle terminal resize
    }
    return m, nil
}
```

### UI Constants
```go
const (
    maxWidth  = 120  // Max UI width
    maxHeight = 30   // Max UI height
    listPaneWidth    = 40
    previewPaneWidth = 60
    gap = 2
)
```

## Nerd Font Icons

Use unicode escapes for Nerd Font icons:
```go
Icons.Folder = "\uf07c"       // nf-fa-folder
Icons.Git = "\ue725"          // nf-dev-git_branch
Icons.Branch = "\ue725"       // nf-oct-git_branch
Icons.Worktree = "\uf115"     // nf-fa-folder_open
Icons.Check = "\uf00c"        // nf-fa-check
Icons.Warning = "\uf071"      // nf-fa-exclamation_triangle
Icons.Modified = "\uf040"     // nf-fa-pencil
Icons.Search = "\uf002"       // nf-fa-search
Icons.ChevronRight = "\uf054" // nf-fa-chevron_right
```

## Configuration

### Config File Locations (in order)
1. `~/.config/switcheroo/config.json` (XDG standard)
2. `~/.switcheroo/config.json` (fallback)

### Config Schema
```json
{
  "directory": "/path/to/repos",     // Required
  "output": {                         // Optional, default type: "path"
    "type": "command",              // "path" prints, "command" executes
    "value": "zellij attach -c {{path}}" // Required when type == "command"
  },
  "preview": {
    "enabled": true,                  // Optional, default: true
    "repo_fields": ["name", "branches", "status"],
    "worktree_fields": ["branch", "status"]
  }
}
```

### Available Preview Fields

**Repo Fields:**
- `name`, `path`, `branches`, `current_branch`
- `status`, `last_commit`, `commit_hash`
- `remote_url`, `ahead_behind`

**Worktree Fields:**
- `name`, `path`, `branch`
- `status`, `last_commit`, `commit_hash`
- `ahead_behind`, `is_locked`

## Git Commands Used

```bash
# Check if git repo
git -C <path> rev-parse --git-dir

# Get worktrees
git -C <path> worktree list --porcelain

# Get branches
git -C <path> branch --format="%(refname:short)|%(HEAD)"

# Get current branch
git -C <path> branch --show-current
```

## Common Tasks

### Adding a New Preview Field
1. Add field to `internal/git/metadata.go` if needed
2. Add git command to fetch data in `internal/git/git.go`
3. Update `RepoDisplay` in `internal/ui/repo_display.go`
4. Add rendering logic in `internal/ui/repo_select.go` renderPreview()
5. Add to default config in `internal/config/defaults.go`

### Adding a New Keybinding
1. Add case in `Update()` method in `repo_select.go`
2. Update help text in `renderHelp()`
3. Document in README.md

### Changing Theme Colors
1. Modify colors in `internal/ui/theme.go` DefaultTheme()
2. Colors use 256-color palette (e.g., `lipgloss.Color("12")`)
3. Update both color constants and style builders

## Testing

### Manual Testing Checklist
- [ ] Works with repos that have worktrees
- [ ] Works with repos without worktrees
- [ ] Scrolling works for long repo lists
- [ ] Scrolling works for long branch lists in preview
- [ ] Fuzzy search filters correctly
- [ ] Config loads from both locations
- [ ] Missing config shows helpful error
- [ ] Window resize handled gracefully

### Test Commands
```bash
# Build and run
go build -o switcheroo-v2 ./cmd/switcheroo
./switcheroo-v2

# Run with specific config
./switcheroo-v2 --config /path/to/config.json

# Check for linting issues
go vet ./...
```

## Future Enhancements

### Plugin System (Planned)
- JavaScript via Goja
- Hook points: Pre/Post selection, Metadata providers
- Interface defined in `internal/plugin/`

### Worktree Selection (Planned)
- Second step after repo selection
- Shows worktrees for worktree-enabled repos
- Same UI pattern as repo selection

## Debugging Tips

### Common Issues
1. **Nerd Fonts not showing**: Terminal must have Nerd Font installed
2. **Colors look wrong**: Terminal must support 256 colors
3. **Git commands fail**: Git must be in PATH
4. **Config not found**: Check file permissions and JSON validity

### Useful Logs
Add temporary logging:
```go
import "fmt"
// In Update or other methods:
fmt.Fprintf(os.Stderr, "Debug: cursor=%d, repos=%d\n", m.cursor, len(m.repos))
```

## Resources

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [Nerd Fonts Cheat Sheet](https://www.nerdfonts.com/cheat-sheet)
- [Git Worktree Documentation](https://git-scm.com/docs/git-worktree)
- [DESIGN.md](./DESIGN.md) - Detailed design document

## Questions?

If you're an AI assistant and have questions:
1. Check DESIGN.md for architecture details
2. Check existing code for patterns
3. Ask the user for clarification on requirements
4. When in doubt, prefer simplicity over complexity

## License

This project is open source. Contributors welcome!
