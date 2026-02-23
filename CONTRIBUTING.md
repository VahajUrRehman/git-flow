# Contributing to GitFlow TUI

Thank you for your interest in contributing to GitFlow TUI! This document provides guidelines and instructions for contributing.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Testing](#testing)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to:
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Respect different viewpoints

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/tui.git
   cd tui
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/gitflow/tui.git
   ```

## Development Setup

### Prerequisites
- Go 1.21+
- Git 2.20+
- Make

### Install Dependencies
```bash
make deps
make install-tools
```

### Build
```bash
make build
```

### Run Tests
```bash
make test
```

## Making Changes

### Branch Naming
- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `docs/description` - Documentation
- `refactor/description` - Code refactoring

Example:
```bash
git checkout -b feature/add-dark-theme
```

### Commit Messages
Follow conventional commits:
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

Example:
```
feat(ui): add dark theme support

Implement dark theme with configurable colors.
Add theme toggle in settings.

Closes #123
```

### Code Structure

```
gitflow-tui/
├── cmd/
│   └── gitflow-tui/        # Entry point
├── internal/
│   ├── git/                # Git operations
│   ├── ui/                 # Terminal UI
│   ├── config/             # Configuration
│   └── auth/               # Authentication
├── pkg/
│   └── graph/              # Graph visualization
├── editors/
│   ├── nvim/               # Neovim plugin
│   └── vscode/             # VSCode extension
└── docs/                   # Documentation
```

## Submitting Changes

1. **Sync with upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push to your fork**:
   ```bash
   git push origin feature/description
   ```

3. **Create Pull Request**:
   - Go to GitHub
   - Click "New Pull Request"
   - Fill in the template
   - Link related issues

### PR Checklist
- [ ] Tests pass
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Commit messages are clear

## Coding Standards

### Go Code
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Use `golint` for linting
- Add comments for exported functions
- Keep functions focused and small

Example:
```go
// GetCurrentBranch returns the name of the current branch.
func (g *Git) GetCurrentBranch() (string, error) {
    // Implementation
}
```

### Lua Code (Neovim)
- Use 2 spaces for indentation
- Follow Lua style guide
- Add type annotations where helpful

### TypeScript Code (VSCode)
- Use 2 spaces for indentation
- Enable strict mode
- Add JSDoc comments

## Testing

### Unit Tests
```bash
make test
```

### Coverage
```bash
make test-coverage
```

### Integration Tests
```bash
make test-integration
```

### Writing Tests
```go
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

## Documentation

### Code Documentation
- Add package comments
- Document exported functions
- Include usage examples

### User Documentation
- Update README.md for features
- Add to docs/ for detailed guides
- Include screenshots for UI changes

## Release Process

1. Update version in:
   - `cmd/gitflow-tui/main.go`
   - `editors/vscode/package.json`

2. Update CHANGELOG.md

3. Create git tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

4. GitHub Actions will build and release

## Questions?

- Open an issue for discussion
- Join our Discord
- Check existing issues

Thank you for contributing! 🌿
