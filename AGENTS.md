# GitFlow TUI - AI Agent Guide

## Project Overview

GitFlow TUI is a **terminal-based Git management application** written in Go. It provides a beautiful, interactive Text User Interface (TUI) for performing Git operations with features like visual commit graphs, full mouse support, and editor integrations.

### Key Features
- Git graph visualization (ASCII, Unicode, Compact styles)
- Full mouse support with clickable tabs and selections
- Complete Git command support (commit, push, pull, merge, rebase, etc.)
- Authentication support (SSH, HTTPS, Token, OAuth)
- Editor integrations for Neovim and VSCode
- Customizable theme with distinctive color scheme

### Technology Stack
- **Language**: Go 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm Architecture for Go)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Git Operations**: go-git library + Git CLI
- **Authentication**: golang.org/x/term for secure password input

---

## Project Structure

```
gitflow-tui/
├── cmd/gitflow-tui/          # Application entry point
│   └── main.go               # Main function, ASCII banner, CLI args
│
├── internal/                 # Internal packages (not exposed)
│   ├── git/                  # Git operations layer
│   │   └── commands.go       # All Git command implementations
│   ├── ui/                   # Terminal UI layer
│   │   ├── model.go          # Bubble Tea model, views, key bindings
│   │   └── commands.go       # UI command handlers
│   ├── config/               # Configuration management
│   │   └── config.go         # Theme, settings, config file I/O
│   └── auth/                 # Authentication management
│       └── auth.go           # SSH/HTTPS/Token/OAuth handling
│
├── pkg/                      # Public packages (can be imported)
│   └── graph/                # Graph visualization package
│       └── graph.go          # ASCII/Unicode graph rendering
│
├── editors/                  # Editor integrations
│   ├── nvim/                 # Neovim plugin
│   │   └── lua/gitflow/
│   │       └── init.lua      # Lua plugin implementation
│   └── vscode/               # VSCode extension
│       ├── package.json      # Extension manifest
│       ├── tsconfig.json     # TypeScript configuration
│       └── src/
│           ├── extension.ts  # Extension entry point
│           ├── terminal.ts   # Terminal handler
│           ├── provider.ts   # Tree view provider
│           └── statusbar.ts  # Status bar component
│
├── docs/                     # Documentation
│   └── INSTALLATION.md       # Installation guide
│
├── .github/workflows/        # CI/CD
│   └── release.yml           # GitHub Actions release workflow
│
├── go.mod                    # Go module definition
├── Makefile                  # Build automation
├── Dockerfile                # Container image
├── install.sh                # Installation script
├── README.md                 # Main documentation
├── QUICKSTART.md             # Quick start guide
├── CONTRIBUTING.md           # Contribution guidelines
├── CHANGELOG.md              # Version history
├── PROJECT_STRUCTURE.md      # Architecture documentation
└── LICENSE                   # MIT License
```

---

## Build Commands

### Prerequisites
- Go 1.21+
- Git 2.20+
- Make
- (Optional) Node.js 16+ for VSCode extension
- (Optional) golangci-lint for linting

### Build Commands (via Makefile)

```bash
# Build binary for current platform
make build

# Build for all platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)
make build-all

# Run tests
make test

# Run tests with coverage report
make test-coverage

# Run linter (golangci-lint preferred, falls back to go vet)
make lint

# Format code
make fmt

# Install dependencies
make deps

# Install binary to system (/usr/local/bin)
make install

# Uninstall binary
make uninstall

# Clean build artifacts
make clean

# Build VSCode extension
make build-vscode

# Package VSCode extension
make package-vscode

# Build Neovim plugin
make build-nvim

# Install Neovim plugin locally
make install-nvim

# Run the application
make run

# Run with debug mode
make debug

# Create full release (clean, test, lint, package)
make release

# Show all available targets
make help
```

### Direct Go Commands

```bash
# Build binary
go build -o build/gitflow-tui ./cmd/gitflow-tui

# Run tests
go test -v ./...

# Install dependencies
go mod download
go mod tidy

# Install globally
go install github.com/gitflow/tui/cmd/gitflow-tui@latest
```

### VSCode Extension Development

```bash
cd editors/vscode
npm install
npm run compile    # Compile TypeScript
npm run watch      # Watch mode
npm run lint       # Run ESLint
npm run package    # Create .vsix package
```

---

## Code Style Guidelines

### Go Code Style

1. **Formatting**: Use `gofmt` for all Go code
   ```bash
   gofmt -w .
   ```

2. **Linting**: Follow `golint` and `golangci-lint` rules
   ```bash
   golangci-lint run
   ```

3. **Comments**: Add comments for all exported functions
   ```go
   // GetCurrentBranch returns the name of the current branch.
   func (g *Git) GetCurrentBranch() (string, error) {
       // Implementation
   }
   ```

4. **Function Size**: Keep functions focused and small

5. **Error Handling**: Always check errors and return descriptive messages

6. **Follow Effective Go**: https://golang.org/doc/effective_go

### Lua Code Style (Neovim Plugin)

- Use 2 spaces for indentation (no tabs)
- Follow Lua style guide
- Add type annotations where helpful

### TypeScript Code Style (VSCode Extension)

- Use 2 spaces for indentation
- Enable strict mode (configured in tsconfig.json)
- Add JSDoc comments for functions

---

## Testing Instructions

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration
```

### Test Structure

Tests should follow Go conventions with `*_test.go` files:

```go
// internal/git/commands_test.go example
func TestGetCurrentBranch(t *testing.T) {
    g := New(".")
    branch, err := g.GetCurrentBranch()
    
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if branch == "" {
        t.Error("Expected branch name, got empty string")
    }
}
```

### Testing Strategy

- **Unit Tests**: Test individual functions in `internal/*/*_test.go`
- **Integration Tests**: Test Git operations in `tests/integration/`
- **E2E Tests**: Full application tests in `tests/e2e/`

---

## Development Conventions

### Git Workflow

1. **Branch Naming**:
   - `feature/description` - New features
   - `bugfix/description` - Bug fixes
   - `docs/description` - Documentation
   - `refactor/description` - Code refactoring

2. **Commit Messages**: Follow Conventional Commits
   ```
   <type>(<scope>): <subject>
   
   <body>
   
   <footer>
   ```
   
   Types:
   - `feat`: New feature
   - `fix`: Bug fix
   - `docs`: Documentation
   - `style`: Formatting
   - `refactor`: Code restructuring
   - `test`: Tests
   - `chore`: Maintenance

3. **Example**:
   ```
   feat(ui): add dark theme support
   
   Implement dark theme with configurable colors.
   Add theme toggle in settings.
   
   Closes #123
   ```

### Release Process

1. Update version in:
   - `cmd/gitflow-tui/main.go`
   - `editors/vscode/package.json`

2. Update `CHANGELOG.md`

3. Create git tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

4. GitHub Actions will build and release automatically

---

## Configuration

### Configuration File

Location: `~/.config/gitflow-tui/config.json`

```json
{
  "theme": {
    "name": "gitflow",
    "colors": {
      "primary": "#00D9A5",
      "secondary": "#00B4A6",
      "tertiary": "#0091EA",
      "accent": "#00E5FF",
      "highlight": "#FF6D00",
      "background": "#0D1117",
      "foreground": "#E6EDF3",
      "success": "#3FB950",
      "warning": "#FFA500",
      "error": "#F85149",
      "muted": "#8B949E",
      "border": "#30363D"
    }
  },
  "git_path": "git",
  "editor": "vim",
  "default_branch": "main",
  "show_graph": true,
  "graph_style": "unicode",
  "mouse_enabled": true,
  "animations": true,
  "auth_method": "ssh",
  "recent_repos": [],
  "max_recent_repos": 10
}
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `GITFLOW_CONFIG` | Path to config file |
| `GITFLOW_THEME` | Override theme |
| `GITFLOW_EDITOR` | Override editor |

---

## Dependencies

### Main Dependencies (go.mod)

```go
require (
    github.com/charmbracelet/bubbles v0.18.0      // UI components
    github.com/charmbracelet/bubbletea v0.25.0    // TUI framework
    github.com/charmbracelet/lipgloss v0.9.1      // Styling
    github.com/charmbracelet/log v0.3.1           // Logging
    github.com/go-git/go-git/v5 v5.11.0           // Git operations
    github.com/awesome-gocui/gocui v1.1.0         // Alternative TUI
    github.com/atotto/clipboard v0.1.4            // Clipboard operations
    golang.org/x/term v0.15.0                     // Terminal utilities
)
```

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        GitFlow TUI                           │
├─────────────────────────────────────────────────────────────┤
│  UI Layer (Bubble Tea)                                      │
│  ├── Dashboard View                                         │
│  ├── Graph View                                             │
│  ├── Branches View                                          │
│  ├── Status View                                            │
│  ├── Stash View                                             │
│  ├── Remotes View                                           │
│  └── Tags View                                              │
├─────────────────────────────────────────────────────────────┤
│  Git Operations Layer                                       │
│  ├── Commit, Push, Pull                                     │
│  ├── Branch, Merge, Rebase                                  │
│  ├── Stash, Tag                                             │
│  └── Diff, Log                                              │
├─────────────────────────────────────────────────────────────┤
│  Core Services                                              │
│  ├── Config Manager                                         │
│  ├── Auth Manager                                           │
│  └── Graph Renderer                                         │
├─────────────────────────────────────────────────────────────┤
│  External Interfaces                                        │
│  ├── Git CLI                                                │
│  ├── SSH/HTTPS Auth                                         │
│  └── Editor APIs                                            │
└─────────────────────────────────────────────────────────────┘
```

### Module Dependencies

```
cmd/gitflow-tui
    ├── internal/ui
    │   ├── internal/git
    │   ├── internal/config
    │   └── pkg/graph
    ├── internal/auth
    └── internal/config
```

---

## Security Considerations

1. **Credential Storage**: Credentials are stored in `~/.config/gitflow-tui/credentials` with restrictive permissions (0600)

2. **Password Input**: Uses `golang.org/x/term` for secure password input without echo

3. **SSH Keys**: Supports SSH key-based authentication with proper key generation

4. **No Hardcoded Secrets**: Never commit credentials or secrets to the repository

---

## Common Tasks

### Adding a New Git Command

1. Add command implementation in `internal/git/commands.go`
2. Add UI handler in `internal/ui/commands.go`
3. Add key binding in `internal/ui/model.go`
4. Update help text

### Adding a New View

1. Add view state constant in `internal/ui/model.go`
2. Add rendering method
3. Add key handlers
4. Update tab list if needed

### Modifying the Theme

1. Update color constants in `internal/config/config.go`
2. Update all references in `internal/ui/model.go`
3. Update editor plugin themes accordingly

---

## Troubleshooting

### Build Issues

```bash
# Clean and rebuild
make clean
make deps
make build
```

### Test Failures

```bash
# Run with verbose output
go test -v ./...

# Run specific test
go test -v -run TestFunctionName ./internal/git
```

### VSCode Extension Issues

```bash
cd editors/vscode
rm -rf node_modules out
npm install
npm run compile
```

---

## License

MIT License - See [LICENSE](LICENSE) for details.

---

## Additional Resources

- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
- [go-git Documentation](https://github.com/go-git/go-git)
- [Effective Go](https://golang.org/doc/effective_go)
