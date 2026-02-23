# GitFlow TUI - Quick Start Guide

Get up and running with GitFlow TUI in 5 minutes!

---

## ⚡ Installation (1 minute)

Choose your platform:

<table>
<tr>
<th>Linux/macOS</th>
<th>Windows (PowerShell)</th>
</tr>
<tr>
<td>

```bash
curl -sSL https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.sh | bash
```

</td>
<td>

```powershell
irm https://raw.githubusercontent.com/VahajUrRehman/git-flow/main/install.ps1 | iex
```

</td>
</tr>
</table>

<details>
<summary><b>📦 Other installation methods</b></summary>

**Homebrew (macOS/Linux):**
```bash
brew tap vahaj/gitflow
brew install gitflow-tui
```

**Scoop (Windows):**
```powershell
scoop bucket add gitflow https://github.com/VahajUrRehman/git-flow
scoop install gitflow-tui
```

**Go Install:**
```bash
go install github.com/VahajUrRehman/git-flow/cmd/gitflow-tui@latest
```

**Download manually:** See [releases page](https://github.com/VahajUrRehman/git-flow/releases)

</details>

---

## 🚀 First Launch (1 minute)

Navigate to any Git repository and run:

```bash
cd my-project
gitflow-tui
```

You'll see:
1. 🎨 **Animated splash screen** with ASCII banner
2. 📊 **Dashboard view** showing repository overview
3. 🖱️ **Interactive interface** with tabs and lists

---

## 🎮 Basic Navigation (2 minutes)

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` `↓` `←` `→` | Navigate lists |
| `k` `j` `h` `l` | Vim-style navigation |
| `Tab` / `Shift+Tab` | Switch tabs |
| `Enter` | Select / Open |
| `Space` | Stage/Unstage file |
| `r` | Refresh data |
| `?` | Show help |
| `q` / `Ctrl+C` | Quit |

### Mouse Support

- **Click tabs** to switch views
- **Click items** to select
- **Double-click** to open/execute

---

## 📊 Understanding the Views

### 1. Dashboard
Overview of your repository:
- Current branch
- Recent commits (colored!)
- Working tree status

### 2. Graph View 🌈
Visual commit history:
- Each branch has unique color
- Green → Teal → Blue → Firozi → Orange
- Commit hashes in blue
- Messages in white

### 3. Status View 📝
Working tree changes:
- 🟢 **Green** = Staged files
- 🟠 **Orange** = Unstaged files
- ⚪ **Gray** = Untracked files

### 4. Branch View 🌿
All branches with info:
- ● Current branch (orange)
- ↑ ahead count (green)
- ↓ behind count (red)

---

## ⚡ Common Operations

### Stage Files
```
1. Go to Status view (Tab)
2. Select file with ↑/↓
3. Press Space to stage/unstage
```

### Commit Changes
```
1. Stage your files (Space)
2. Press : for command mode
3. Type :commit
4. Enter message
5. Press Enter
```

### Switch Branch
```
1. Go to Branch view (Tab)
2. Select branch with ↑/↓
3. Press Enter to checkout
```

### Push to Remote
```
1. Press : for command mode
2. Type :push
3. Press Enter
```

### Pull Changes
```
1. Press : for command mode
2. Type :pull
3. Press Enter
```

---

## 🎨 The Color Theme

GitFlow uses a distinctive 5-color palette:

| Color | Hex | Used For |
|-------|-----|----------|
| 🟢 Green | `#00D9A5` | Success, primary actions |
| 🔵 Teal | `#00B4A6` | Secondary, staged files |
| 🔷 Blue | `#0091EA` | Tertiary, commit hashes |
| 💎 Firozi | `#00E5FF` | Accent, selected items |
| 🟠 Orange | `#FF6D00` | Highlight, warnings, current branch |

---

## ⚙️ Configuration

Config file: `~/.config/gitflow-tui/config.json`

```json
{
  "theme": {
    "name": "gitflow",
    "colors": {
      "primary": "#00D9A5",
      "secondary": "#00B4A6",
      "tertiary": "#0091EA",
      "accent": "#00E5FF",
      "highlight": "#FF6D00"
    }
  },
  "mouse_enabled": true,
  "animations": true,
  "graph_style": "unicode"
}
```

---

## 🆘 Need Help?

| Resource | Link |
|----------|------|
| 📖 Full Documentation | [README.md](README.md) |
| 🔧 Build from Source | [BUILD_GUIDE.md](BUILD_GUIDE.md) |
| 🤝 Contributing | [CONTRIBUTING.md](CONTRIBUTING.md) |
| 🐛 Report Issues | [GitHub Issues](https://github.com/VahajUrRehman/git-flow/issues) |

---

## 💡 Pro Tips

1. **Use mouse and keyboard together** - Mouse for quick clicks, keyboard for fast navigation

2. **Press ? anytime** for context-sensitive help

3. **Use command mode (`:`)** for quick git operations without leaving the TUI

4. **Enable animations** in config for smoother experience

5. **Customize colors** to match your terminal theme

---

<div align="center">

**Enjoy using GitFlow TUI!** 🚀

[⭐ Star on GitHub](https://github.com/VahajUrRehman/git-flow) • [🐛 Report Bug](https://github.com/VahajUrRehman/git-flow/issues)

</div>
