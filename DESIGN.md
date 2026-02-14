# Switcheroo v2 - Design Document

## Overview

Switcheroo is a terminal-based Git repository and worktree switcher built with Go and Bubble Tea. It provides a fuzzy-findable interface for navigating repositories and their worktrees, with configurable preview metadata and output options.

### Core Purpose

- **Primary use case**: Navigate Git repositories that use worktrees for parallel branch development
- **Secondary use case**: Navigate regular Git repositories without worktrees
- **Output flexibility**: Either output a path for shell scripts or execute a configured command (e.g., open Zellij/tmux session)

## Design Principles

1. **Simplicity**: Go is chosen for its simplicity and approachability for contributors
2. **Performance**: Fuzzy finding and Git operations must be fast
3. **Flexibility**: Support both worktree-based and regular Git workflows
4. **Extensibility**: Architecture designed to support future plugin system without current implementation
5. **User Control**: Preview metadata is fully configurable

## Architecture Overview

### Technology Stack

- **Language**: Go 1.22+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm Architecture)
- **Fuzzy Finding**: fzf library or go-fuzzyfinder
- **Future Plugin System**: [Goja](https://github.com/dop251/goja) (JavaScript in Go)

### Project Structure

```
switcheroo/
├── cmd/
│   └── switcheroo/
│       └── main.go              # Entry point and main flow
├── internal/
│   ├── config/
│   │   ├── config.go            # Configuration parsing and types
│   │   └── defaults.go          # Default configuration values
│   ├── git/
│   │   ├── repository.go        # Repository discovery and operations
│   │   ├── worktree.go          # Worktree detection and parsing
│   │   ├── branch.go            # Branch listing and operations
│   │   ├── metadata.go          # Metadata collection (status, commits, etc.)
│   │   └── git.go               # Low-level git command execution
│   ├── ui/
│   │   ├── model.go             # Bubble Tea root model
│   │   ├── repo_select.go       # Repository selection view
│   │   ├── worktree_select.go   # Worktree selection view
│   │   ├── preview.go           # Preview pane rendering
│   │   └── styles.go            # Terminal styling (colors, formatting)
│   └── plugin/                  # Future plugin system (interfaces only)
│       ├── interface.go         # Plugin interface definitions
│       └── hooks.go             # Hook point definitions (empty implementations)
├── go.mod
├── go.sum
├── README.md
└── DESIGN.md                    # This document
```

## Configuration

### Configuration File Location

- Primary: `$XDG_CONFIG_HOME/switcheroo/config.json`
- Fallback: `~/.config/switcheroo/config.json`

### Configuration Schema

```json
{
  "directory": "/Users/username/repos",
  "output": {
    "type": "command",
    "value": "zellij attach -c {{path}}"
  },
  "preview": {
    "enabled": true,
    "repo_fields": [
      "name",
      "path",
      "worktree_count",
      "branches",
      "current_branch",
      "status",
      "last_commit",
      "commit_hash",
      "remote_url",
      "ahead_behind"
    ],
    "worktree_fields": [
      "name",
      "path",
      "branch",
      "status",
      "last_commit",
      "commit_hash",
      "ahead_behind",
      "is_locked"
    ]
  }
}
```

### Configuration Fields

#### Root Level

- `directory` (string, required): Absolute path to directory containing all repositories
- `output` (object, default: `{ "type": "path" }`): Controls what happens after selection
  - `type` (string): `"path"` prints the selected path; `"command"` executes `value`
  - `value` (string, required when type is `"command"`): Command template supporting `{{path}}`
- `preview` (object): Preview pane configuration

#### Preview Configuration

- `enabled` (boolean, default: true): Whether to show preview pane
- `repo_fields` (array of strings): Metadata fields to display for repositories
- `worktree_fields` (array of strings): Metadata fields to display for worktrees

### Default Configuration

If no config file exists, defaults are:

```json
{
  "directory": "",
  "output": {
    "type": "path"
  },
  "preview": {
    "enabled": true,
    "repo_fields": ["name", "branches", "status", "last_commit"],
    "worktree_fields": ["branch", "status", "last_commit", "ahead_behind"]
  }
}
```

Note: `directory` must be explicitly set by the user.

## Core Data Types

### Repository

Represents a Git repository.

```go
type Repository struct {
    Name          string      // Directory name
    Path          string      // Absolute path to repository
    HasWorktrees  bool        // Whether this repo has additional worktrees
    Worktrees     []Worktree  // List of worktrees (if any)
    Branches      []Branch    // All local branches
    CurrentBranch string      // Currently checked out branch
    RemoteURL     string      // Git remote origin URL
}
```

### Worktree

Represents a Git worktree.

```go
type Worktree struct {
    Name       string  // Directory name
    Path       string  // Absolute path to worktree
    Branch     string  // Branch checked out in this worktree
    IsLocked   bool    // Whether worktree is locked
    IsBare     bool    // Whether this is a bare worktree
    CommitHash string  // Current commit hash (short)
}
```

### Branch

Represents a Git branch.

```go
type Branch struct {
    Name      string  // Branch name
    IsCurrent bool    // Whether this is the currently checked out branch
}
```

### Metadata

Runtime metadata collected about repositories and worktrees.

```go
type Metadata struct {
    Status      string          // "clean", "modified", "untracked"
    LastCommit  CommitInfo      // Most recent commit information
    AheadBehind AheadBehindInfo // Sync status with remote
}

type CommitInfo struct {
    Hash    string    // Short SHA
    Message string    // Commit message (first line)
    Author  string    // Commit author name
    Date    time.Time // Commit timestamp
}

type AheadBehindInfo struct {
    Ahead  int  // Commits ahead of remote
    Behind int  // Commits behind remote
}
```

## Application Flow

### Main Execution Flow

```
1. Load configuration from config file
   ├─ Parse JSON
   ├─ Apply defaults for missing fields
   └─ Validate required fields

2. Discover repositories in configured directory
   ├─ Scan directory for subdirectories
   ├─ Filter for Git repositories (.git exists)
   ├─ For each repository:
   │  ├─ Detect worktrees (see Worktree Detection)
   │  ├─ Get branch list
   │  └─ Collect metadata based on repo_fields config
   └─ Return list of Repository objects

3. Display repository selection UI
   ├─ Fuzzy finder list of repository names
   ├─ Preview pane showing configured metadata
   └─ User selects a repository

4. Check if selected repository has worktrees
   ├─ If HasWorktrees == true:
   │  ├─ Display worktree selection UI
   │  ├─ Fuzzy finder list of worktree names
   │  ├─ Preview pane showing configured metadata
   │  ├─ User selects a worktree
   │  └─ selectedPath = worktree.Path
   └─ If HasWorktrees == false:
      └─ selectedPath = repository.Path

5. Handle output based on configuration
   ├─ If output == "path":
   │  └─ Print selectedPath to stdout
   └─ If output == "command":
      ├─ Replace {{path}} in command template
      └─ Execute command
```

### Worktree Detection Logic

Switcheroo supports two methods of worktree detection with a priority order:

```
1. Directory-based detection (PRIORITY 1)
   ├─ Check if {repo_path}/worktrees/ directory exists
   ├─ If exists:
   │  ├─ Scan for subdirectories
   │  ├─ For each subdirectory:
   │  │  ├─ Verify it's a valid worktree (has .git file)
   │  │  └─ Create Worktree object
   │  └─ If worktrees found, return them
   └─ If no worktrees found, continue to next method

2. Git command detection (FALLBACK)
   ├─ Execute: git -C {repo_path} worktree list --porcelain
   ├─ Parse output to extract worktree information
   ├─ Filter out main worktree (the repository itself)
   └─ Return additional worktrees only

3. Result
   ├─ If worktrees found: HasWorktrees = true
   └─ If no worktrees found: HasWorktrees = false
```

**Important**: The main repository checkout is NOT included in the worktree list. Only additional worktrees are shown in the worktree selection step.

### Branch Listing

```
1. Execute: git -C {repo_path} branch --format=%(refname:short)|%(HEAD)
2. Parse output:
   ├─ Split by newlines
   ├─ For each line:
   │  ├─ Split by "|" delimiter
   │  ├─ Extract branch name
   │  ├─ Check if current (marked with "*")
   │  └─ Create Branch object
   └─ Sort branches alphabetically by name
3. Return sorted branch list with current branch flagged
```

## Preview Metadata Fields

### Repository Preview Fields

Available fields for `preview.repo_fields`:

| Field | Description | Example Output |
|-------|-------------|----------------|
| `name` | Repository directory name | `my-project` |
| `path` | Full absolute path | `/Users/me/repos/my-project` |
| `worktree_count` | Number of additional worktrees | `Worktrees: 3` or `Worktrees: None` (disabled) |
| `branches` | List of all local branches, current highlighted | See Branch Display Format |
| `branch_status` | Shows `(status)` after each branch | `feature/login (ahead)` / `feature/new-ui (untracked)` |
| `branch_diff` | Shows per-branch diff summary when upstream is configured | `feature/login +14 -9` |
| `line_diff` | Global diff summary for current branch (falls back to default branch when untracked) | `+14 additions / -9 deletions` |
| `current_branch` | Just the current branch name | `Current: feature/new-ui` |
| `status` | Working directory status | `Clean` / `Modified (3 files)` / `Untracked files` |
| `last_commit` | Most recent commit | `feat: add feature (2h ago) - John Doe` |
| `commit_hash` | Short commit SHA | `a3f5b2c` |
| `remote_url` | Git remote origin URL | `git@github.com:user/repo.git` |
| `ahead_behind` | Sync status with remote | `↑2 ↓1` (2 ahead, 1 behind) |

Branches without an upstream are labeled `(untracked)` in the preview and omit per-branch diff counts. The global `line_diff` summary compares the current branch to its upstream or, if none exists, to the repository's default branch.

### Worktree Preview Fields

Available fields for `preview.worktree_fields`:

| Field | Description | Example Output |
|-------|-------------|----------------|
| `name` | Worktree directory name | `feature-branch` |
| `path` | Full absolute path | `/Users/me/repos/project/worktrees/feature-branch` |
| `branch` | Branch checked out | `feature/new-ui` |
| `status` | Working directory status | `Clean` / `Modified (3 files)` |
| `line_diff` | Diff summary for the worktree | `+8 / -3` |
| `last_commit` | Most recent commit | `feat: add feature (2h ago) - John Doe` |
| `commit_hash` | Short commit SHA | `a3f5b2c` |
| `ahead_behind` | Sync status with remote | `↑2 ↓1` |
| `is_locked` | Whether worktree is locked | `🔒 Locked` / (omitted if not locked) |

### Branch Display Format

When `branches` field is included in preview, it renders as:

```
Branches:
* feature/new-ui
  fix/bug-123
  main
  staging
```

- Branches sorted alphabetically
- Current branch prefixed with `*` and highlighted (bold, different color)
- No limit on number of branches displayed

### Worktree Count Display

When `worktree_count` field is included:

- If `HasWorktrees == true`: Display `Worktrees: N` (where N > 0)
- If `HasWorktrees == false`: Display `Worktrees: None` as **disabled/grayed out text**
  - This line should be visually distinct (dimmed color)
  - Cannot be interacted with
  - Indicates this repository does not use worktrees

## User Interface

### Bubble Tea Model Structure

The application uses the Elm Architecture pattern via Bubble Tea:

```go
type Model struct {
    state         AppState        // Current state (RepoSelect, WorktreeSelect, Complete)
    config        *config.Config  // Loaded configuration
    repositories  []Repository    // All discovered repositories
    selectedRepo  *Repository     // Currently selected repository
    worktrees     []Worktree      // Worktrees for selected repo
    selectedWorktree *Worktree    // Selected worktree (if applicable)
    fuzzyModel    tea.Model       // Fuzzy finder component
    previewModel  tea.Model       // Preview pane component
    err           error           // Error state
}

type AppState int

const (
    StateRepoSelect AppState = iota
    StateWorktreeSelect
    StateComplete
    StateError
)
```

### Repository Selection View

**Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Select Repository                                        │
├─────────────────────┬───────────────────────────────────┤
│                     │                                   │
│ > my-project        │ Name: my-project                  │
│   another-repo      │ Path: /Users/me/repos/my-project  │
│   third-project     │                                   │
│   util-library      │ Worktrees: 3                      │
│                     │                                   │
│                     │ Branches:                         │
│                     │ * feature/new-ui                  │
│                     │   main                            │
│                     │   staging                         │
│                     │                                   │
│                     │ Status: Clean                     │
│                     │ Last Commit: feat: add...         │
└─────────────────────┴───────────────────────────────────┘
```

- Left pane: Fuzzy-findable list of repositories
- Right pane: Preview metadata based on `repo_fields` config
- Navigation: Arrow keys, Vim keys (j/k), or fuzzy search
- Selection: Enter key

### Worktree Selection View

**Layout:**
```
┌─────────────────────────────────────────────────────────┐
│ Select Worktree: my-project                              │
├─────────────────────┬───────────────────────────────────┤
│                     │                                   │
│ > feature-new-ui    │ Name: feature-new-ui              │
│   fix-bug-123       │ Path: .../worktrees/feature-...   │
│   main              │                                   │
│                     │ Branch: feature/new-ui            │
│                     │ Status: Modified (2 files)        │
│                     │ Last Commit: wip: working on...   │
│                     │ Ahead/Behind: ↑3                  │
│                     │                                   │
└─────────────────────┴───────────────────────────────────┘
```

- Left pane: Fuzzy-findable list of worktrees
- Right pane: Preview metadata based on `worktree_fields` config
- Navigation: Same as repository selection
- Selection: Enter key

### Styling Guidelines

Using Bubble Tea's Lip Gloss for styling:

- **Current selection**: Highlighted background, bold text
- **Current branch indicator**: Bold, accent color (e.g., cyan)
- **Disabled text** (e.g., "Worktrees: None"): Dimmed/gray color
- **Preview labels**: Subtle color (gray)
- **Preview values**: Normal text color
- **Borders**: Subtle box drawing characters
- **Status colors**:
  - Clean: Green
  - Modified: Yellow
  - Error: Red

## Git Command Operations

### Commands Used

All Git operations use `git -C {path}` to operate on specific repositories.

| Operation | Command | Notes |
|-----------|---------|-------|
| List branches | `git branch --format=%(refname:short)\|%(HEAD)` | Porcelain output for parsing |
| Current branch | `git branch --show-current` | Returns empty if detached HEAD |
| Worktree list | `git worktree list --porcelain` | Structured output for parsing |
| Git status | `git status --porcelain` | Machine-readable format |
| Last commit | `git log -1 --format=%h\|%s\|%an\|%ar` | Hash, message, author, relative date |
| Ahead/behind | `git rev-list --left-right --count @{u}...HEAD` | Requires tracking branch |
| Remote URL | `git remote get-url origin` | May not exist for all repos |

### Error Handling

Git commands may fail in various scenarios:

- Repository not initialized
- No remote configured
- Detached HEAD state
- No commits yet
- Permission errors

**Strategy**: Gracefully degrade, show partial information
- If command fails, omit that metadata field from preview
- Log errors for debugging but don't crash application
- Show "Unknown" or omit field entirely rather than error messages in preview

## Output Modes

### Path Mode (`output.type = "path"`)

Prints the selected path to stdout and exits.

**Example:**
```bash
$ switcheroo
# User selects repo "my-project" → worktree "feature-a"
/Users/me/repos/my-project/worktrees/feature-a

$ cd $(switcheroo)  # Common usage pattern
```

### Command Mode (`output.type = "command"`)

Executes the configured command with path variable substitution.

**Template variables:**
- `{{path}}`: Replaced with selected path

**Example configuration:**
```json
{
  "output": {
    "type": "command",
    "value": "zellij attach -c {{path}}"
  }
}
```

**Execution:**
```bash
$ switcheroo
# User selects worktree with path: /Users/me/repos/my-project/worktrees/feature-a
# Executes: zellij attach -c /Users/me/repos/my-project/worktrees/feature-a
```

**Implementation:**
See `cmd/switcheroo/main.go` for the current implementation, which replaces `{{path}}` in `output.value`, splits it into args, and executes it via `exec.Command`.

## Future Plugin System

**Status**: Not implemented in initial version, but architecture designed to support it.

### Design Considerations

The codebase includes interface definitions and hook points for a future plugin system without current implementation. This allows plugins to be added later without major refactoring.

### Plugin Interface (Planned)

```go
// internal/plugin/interface.go

// MetadataProvider allows plugins to add custom metadata fields
type MetadataProvider interface {
    Name() string
    GetMetadata(item interface{}) (key string, value string, error)
}

// FilterProvider allows plugins to filter/transform lists
type FilterProvider interface {
    Name() string
    FilterRepositories(repos []Repository) []Repository
    FilterWorktrees(worktrees []Worktree) []Worktree
}

// ActionProvider allows plugins to run on selection
type ActionProvider interface {
    Name() string
    OnRepoSelect(repo Repository) error
    OnWorktreeSelect(worktree Worktree) error
}
```

### Hook Points (Planned)

```go
// internal/plugin/hooks.go

type Hooks struct {
    // Metadata hooks
    MetadataProviders []MetadataProvider
    
    // Filter hooks
    PreRepoSelect     []func(repos []Repository) []Repository
    PostRepoSelect    []func(repo Repository) error
    PreWorktreeSelect []func(worktrees []Worktree) []Worktree
    PostWorktreeSelect []func(worktree Worktree) error
}

// Empty implementations for now
var GlobalHooks = &Hooks{
    MetadataProviders:  []MetadataProvider{},
    PreRepoSelect:      []func([]Repository) []Repository{},
    PostRepoSelect:     []func(Repository) error{},
    PreWorktreeSelect:  []func([]Worktree) []Worktree{},
    PostWorktreeSelect: []func(Worktree) error{},
}
```

### Plugin Language (Planned)

**Chosen**: JavaScript via [Goja](https://github.com/dop251/goja)

**Rationale**:
- Pure Go implementation (no CGO)
- Everyone knows JavaScript
- ~3-5MB binary size increase
- Good sandboxing
- Easy to expose Go functions to plugins

**Example plugin (future):**
```javascript
// ~/.config/switcheroo/plugins/tmux-session.js

Config = {
  name: "Tmux Session Creator",
  type: "action"
};

function OnWorktreeSelect(worktree) {
  const sessionName = worktree.name.replace(/\//g, "-");
  exec(`tmux new-session -s ${sessionName} -c ${worktree.path}`);
}
```

### When to Implement

Plugins should be implemented when:
1. Core functionality is stable
2. Community requests extensibility
3. Common use cases emerge that aren't core features

**Not a priority for v1.0.**

## Implementation Phases

### Phase 1: Core Functionality (MVP)
- [ ] Configuration loading
- [ ] Repository discovery
- [ ] Worktree detection (both methods)
- [ ] Branch listing
- [ ] Basic Bubble Tea UI (repo selection only)
- [ ] Path output mode

### Phase 2: Enhanced UI
- [ ] Worktree selection UI
- [ ] Preview pane rendering
- [ ] Configurable preview fields
- [ ] Styling and polish

### Phase 3: Metadata Collection
- [ ] Git status detection
- [ ] Last commit information
- [ ] Ahead/behind tracking
- [ ] Remote URL detection
- [ ] All preview field implementations

### Phase 4: Polish & Release
- [ ] Command output mode
- [ ] Error handling and graceful degradation
- [ ] CLI flags (if needed)
- [ ] Documentation
- [ ] Binary releases (GitHub Actions)

### Phase 5: Future Enhancements
- [ ] Plugin system implementation
- [ ] Additional output formats
- [ ] Advanced filtering options
- [ ] Custom themes

## Testing Strategy

### Unit Tests

Test coverage for:
- Configuration parsing
- Git command parsing (branch list, worktree list)
- Path manipulation
- Template variable substitution

### Integration Tests

- Repository discovery in test fixtures
- Worktree detection with mock repos
- Git command execution (require git in PATH)

### Manual Testing

- Test with various repository structures:
  - Regular repo (no worktrees)
  - Repo with worktrees/ directory
  - Repo with git worktrees
  - Bare repositories
  - Detached HEAD states
  - No remote configured
  - Many branches (50+)

## Dependencies

```go
// go.mod
module github.com/yourusername/switcheroo

go 1.22

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/charmbracelet/bubbles v0.18.0  // For text input, list components
    github.com/sahilm/fuzzy v0.1.1            // Fuzzy matching algorithm
)
```

**Future plugin dependency** (not included in v1):
```go
github.com/dop251/goja v0.0.0-20231027120936-b396bb4c349d
```

## Performance Considerations

### Repository Scanning

- Scan directory once at startup
- Cache repository metadata during scan
- Avoid re-executing git commands on navigation

### Fuzzy Finding

- Use efficient fuzzy matching algorithm
- Debounce search input to avoid excessive filtering
- For large repo counts (100+), consider virtual scrolling

### Git Operations

- Execute git commands concurrently where possible (goroutines)
- Timeout long-running git operations (e.g., 5 second timeout)
- Cache results that don't change during session (branch list, remote URL)

### Expected Performance

- Repository discovery: <100ms for 50 repos
- UI rendering: 60fps
- Fuzzy search response: <50ms

## Cross-Platform Considerations

### Supported Platforms

- macOS (primary development platform)
- Linux
- Windows (WSL and native)

### Path Handling

- Use `filepath.Join()` for cross-platform path construction
- Accept forward slashes in config, normalize internally
- Support both Unix and Windows path formats

### Terminal Compatibility

- Use Bubble Tea's built-in terminal detection
- Test with common terminals:
  - iTerm2 (macOS)
  - Terminal.app (macOS)
  - Alacritty
  - Kitty
  - Windows Terminal
  - Tmux/screen

### Git Command Compatibility

- Assume Git 2.x+ installed
- Check git version at startup, warn if <2.0
- Use only widely-supported git flags

## Error Scenarios & Handling

| Scenario | Handling Strategy |
|----------|-------------------|
| No config file | Create default config, prompt for directory |
| Invalid config JSON | Show parse error, exit with code 1 |
| Directory doesn't exist | Show error, prompt to update config |
| Directory not readable | Show permission error |
| No Git repos found | Show message: "No repositories found in {dir}" |
| Git not installed | Show error: "Git not found in PATH" |
| Git command fails | Gracefully omit that metadata field |
| No worktrees found | Set HasWorktrees=false, show disabled indicator |
| User cancels selection | Exit with code 130 (standard for Ctrl+C) |
| Command execution fails | Show error, exit with command's exit code |

## Security Considerations

### Command Injection

When executing configured commands, prevent injection:
- Do NOT use shell execution (`/bin/sh -c`)
- Use `exec.Command()` with separate args
- Validate path doesn't contain shell metacharacters if paranoid
- Future: Allowlist of safe commands

### Plugin Sandboxing (Future)

When plugins are implemented:
- Use Goja's sandboxing features
- Don't expose `os.exec` directly
- Provide safe wrapper functions
- Limit filesystem access to repo paths
- Timeout plugin execution (e.g., 5 seconds)

## Documentation Requirements

### README.md

- Installation instructions
- Quick start guide
- Configuration examples
- Common workflows (worktree setup)
- Troubleshooting

### CONTRIBUTING.md

- Development setup
- Code style guide
- Testing requirements
- PR process

### Config Documentation

- All available preview fields
- Output modes
- Template variables
- Example configurations for common setups

## Success Metrics

For v1.0 release:
- [ ] Handles 100+ repositories without lag
- [ ] Supports both worktree detection methods
- [ ] Preview metadata configurable
- [ ] Works on macOS, Linux, Windows
- [ ] Binary size <15MB
- [ ] Zero crashes on valid inputs
- [ ] Clear error messages for invalid inputs

## Open Questions

None at this time. Design is complete and ready for implementation.

## Changelog

- **2024-02-12**: Initial design document created
  - Core architecture defined
  - Technology stack selected (Go + Bubble Tea)
  - Two-step selection flow designed
  - Preview metadata fields defined
  - Plugin system planned (not implemented)
  - Worktree detection strategy finalized
