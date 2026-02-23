# 🌿 GitFlow TUI

<p align="center">
  <img src="assets/banner.png" alt="GitFlow TUI Banner" width="800">
</p>

<p align="center">
  <b>Complete Git Management TUI with Beautiful Visualizations</b>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Neovim-Supported-57A143?style=flat-square&logo=neovim" alt="Neovim">
  <img src="https://img.shields.io/badge/VSCode-Supported-007ACC?style=flat-square&logo=visual-studio-code" alt="VSCode">
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=flat-square" alt="License">
</p>

---

## ✨ Features

### 🎨 **Beautiful Color Scheme**
- **Green** `#00D9A5` - Primary actions
- **Teal** `#00B4A6` - Secondary elements
- **Blue** `#0091EA` - Tertiary highlights
- **Firozi/Cyan** `#00E5FF` - Accent color
- **Orange** `#FF6D00` - Warnings and highlights

### 📊 **Git Graph Visualization**
- ASCII, Unicode, and Compact graph styles
- Interactive commit navigation
- Branch visualization
- Merge commit highlighting

### 🖱️ **Full Mouse Support**
- Click to navigate tabs
- Click to select commits/branches
- Scroll through history
- Context menus

### ⌨️ **Complete Git Commands**
| Command | Description |
|---------|-------------|
| `commit` | Create commits with message editor |
| `push` | Push to remote repositories |
| `pull` | Pull from remote with rebase option |
| `fetch` | Fetch remote changes |
| `checkout` | Switch branches |
| `merge` | Merge branches with conflict resolution |
| `rebase` | Interactive rebase support |
| `cherry-pick` | Pick specific commits |
| `stash` | Save/apply stashes |
| `tag` | Create and manage tags |
| `reset` | Soft/mixed/hard reset |
| `revert` | Revert commits |
| `diff` | View file diffs |

### 🔐 **Authentication Support**
- SSH key management
- HTTPS with credential helper
- Personal Access Tokens
- OAuth (GitHub, GitLab, Bitbucket)

### 🔌 **Editor Integration**
- **Neovim** - Full Lua plugin with Telescope integration
- **VSCode** - Complete extension with sidebar and status bar

---

## 🚀 Installation

### Prerequisites
- Go 1.21+ (for building from source)
- Git 2.20+

### From Source

```bash
# Clone the repository
git clone https://github.com/gitflow/tui.git
cd tui

# Build the binary
make build

# Install globally
make install

# Or install to custom location
make install PREFIX=/usr/local
```

### Using Go Install

```bash
go install github.com/gitflow/tui/cmd/gitflow-tui@latest
```

### Package Managers

```bash
# Homebrew (macOS/Linux)
brew tap gitflow/tui
brew install gitflow-tui

# AUR (Arch Linux)
yay -S gitflow-tui

# Scoop (Windows)
scoop bucket add gitflow https://github.com/gitflow/scoop-bucket
scoop install gitflow-tui
```

---

## 🎮 Usage

### Terminal

```bash
# Open GitFlow TUI in current directory
gitflow-tui

# Open in specific directory
gitflow-tui --cwd /path/to/repo

# Show version
gitflow-tui --version

# Show ASCII banner
gitflow-tui --banner
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `←/h` | Move left |
| `→/l` | Move right |
| `Tab` | Next tab |
| `Shift+Tab` | Previous tab |
| `Enter` | Select/Confirm |
| `Space` | Stage/Unstage file |
| `c` | Commit |
| `p` | Push |
| `P` | Pull |
| `f` | Fetch |
| `b` | Checkout branch |
| `m` | Merge |
| `R` | Rebase |
| `r` | Refresh |
| `?` | Help |
| `q` | Quit |

---

## 🔌 Editor Integration

### Neovim

#### Using [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  'gitflow/tui',
  dependencies = {
    'nvim-telescope/telescope.nvim', -- optional
  },
  config = function()
    require('gitflow').setup({
      -- Configuration
      keymaps = {
        open = '<leader>gg',
      },
      theme = {
        primary = '#00D9A5',
        secondary = '#00B4A6',
        tertiary = '#0091EA',
        accent = '#00E5FF',
        highlight = '#FF6D00',
      },
    })
  end,
}
```

#### Using [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  'gitflow/tui',
  config = function()
    require('gitflow').setup()
  end
}
```

#### Commands

```vim
:GitFlow          " Open GitFlow TUI
:GitFlowToggle    " Toggle GitFlow TUI
:GitFlowClose     " Close GitFlow TUI
```

### VSCode

1. Install from VSCode Marketplace or manually:

```bash
cd editors/vscode
npm install
vsce package
code --install-extension gitflow-tui-*.vsix
```

2. Open command palette (`Ctrl+Shift+P` / `Cmd+Shift+P`):
   - `GitFlow: Open GitFlow TUI`

3. Default keybinding:
   - `Ctrl+Shift+G G` / `Cmd+Shift+G G` - Toggle GitFlow TUI

---

## ⚙️ Configuration

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

## 🔐 Authentication

### SSH

```bash
# Generate SSH key
gitflow-tui auth ssh generate --email "your@email.com"

# Add to SSH agent
gitflow-tui auth ssh add

# Copy public key to clipboard
gitflow-tui auth ssh copy
```

### HTTPS

```bash
# Configure HTTPS credentials
gitflow-tui auth https --host github.com --username yourusername
```

### Token

```bash
# Configure personal access token
gitflow-tui auth token --host github.com --token YOUR_TOKEN
```

### OAuth

```bash
# Start OAuth flow
gitflow-tui auth oauth --provider github
```

---

## 📸 Screenshots

<p align="center">
  <img src="assets/screenshot-dashboard.png" alt="Dashboard" width="800">
  <br>
  <em>Dashboard View</em>
</p>

<p align="center">
  <img src="assets/screenshot-graph.png" alt="Git Graph" width="800">
  <br>
  <em>Git Graph Visualization</em>
</p>

<p align="center">
  <img src="assets/screenshot-status.png" alt="Status" width="800">
  <br>
  <em>Working Tree Status</em>
</p>

---

## 🛠️ Development

### Building

```bash
# Build binary
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
make lint
```

### Project Structure

```
gitflow-tui/
├── cmd/
│   └── gitflow-tui/        # Main entry point
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
├── assets/                 # Assets and screenshots
└── docs/                   # Documentation
```

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📜 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [go-git](https://github.com/go-git/go-git) - Git operations

---

<p align="center">
  Made with 💚 by the GitFlow Team
</p>

<p align="center">
  <a href="https://github.com/gitflow/tui">⭐ Star us on GitHub</a>
</p>
