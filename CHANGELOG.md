# Changelog

All notable changes to GitFlow TUI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of GitFlow TUI
- Complete Git management interface
- Beautiful color scheme (Green, Teal, Blue, Firozi, Orange)
- Git graph visualization (ASCII, Unicode, Compact)
- Full mouse support
- All major Git commands
- Authentication support (SSH, HTTPS, Token, OAuth)
- Neovim plugin integration
- VSCode extension integration
- Terminal UI with tabs
- Working tree status view
- Branch management
- Stash management
- Tag management
- Remote management
- Commit history with graph
- Diff viewer
- Configuration system
- Cross-platform support (macOS, Linux, Windows)

## [1.0.0] - 2024-XX-XX

### Added
- Core TUI framework using Bubble Tea
- Git operations module
- Graph visualization package
- Authentication manager
- Configuration management
- Mouse event handling
- Keyboard shortcuts
- Tab-based navigation
- Dashboard view
- Graph view
- Branches view
- Status view
- Stash view
- Remotes view
- Tags view
- Help view
- Input dialogs
- Confirmation dialogs
- Diff viewer
- Status bar
- Error handling
- Loading states
- Theme system
- Editor plugins (Neovim, VSCode)

### Features
- **Dashboard**: Repository overview with recent commits and status
- **Graph**: Visual commit graph with branch lines
- **Branches**: List and manage branches with ahead/behind info
- **Status**: Working tree status with staged/unstaged files
- **Stash**: View and manage stash entries
- **Remotes**: List and configure remotes
- **Tags**: Create and manage tags
- **Help**: Comprehensive keyboard shortcuts

### Git Commands
- `commit`: Create commits with message input
- `push`: Push to remote repositories
- `pull`: Pull from remote with rebase option
- `fetch`: Fetch remote changes
- `checkout`: Switch branches with create option
- `merge`: Merge branches
- `rebase`: Interactive rebase support
- `cherry-pick`: Pick specific commits
- `stash`: Save/apply/pop stashes
- `tag`: Create and manage tags
- `reset`: Soft/mixed/hard reset
- `revert`: Revert commits
- `diff`: View file diffs

### Authentication
- SSH key generation and management
- HTTPS with credential helper
- Personal Access Tokens
- OAuth flow (GitHub, GitLab, Bitbucket)

### Editor Integration
- **Neovim**: Full Lua plugin with Telescope integration
- **VSCode**: Complete extension with sidebar and status bar

### Themes
- Default theme with Green (#00D9A5), Teal (#00B4A6), Blue (#0091EA), Firozi (#00E5FF), Orange (#FF6D00)
- Configurable colors
- Terminal color support

### Platforms
- macOS (Intel & Apple Silicon)
- Linux (AMD64 & ARM64)
- Windows (AMD64)

[Unreleased]: https://github.com/gitflow/tui/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/gitflow/tui/releases/tag/v1.0.0
