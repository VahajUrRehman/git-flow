# Quick Start Guide

Get up and running with GitFlow TUI in minutes!

## Installation (30 seconds)

```bash
# macOS/Linux
curl -sSL https://raw.githubusercontent.com/gitflow/tui/main/install.sh | bash

# Or with Go
go install github.com/gitflow/tui/cmd/gitflow-tui@latest
```

## First Run (1 minute)

```bash
# Navigate to any git repository
cd your-project

# Launch GitFlow TUI
gitflow-tui
```

## Basic Navigation

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Switch tabs |
| `↑/k` / `↓/j` | Navigate |
| `Enter` | Select |
| `Space` | Stage/unstage file |
| `?` | Help |
| `q` | Quit |

## Common Tasks

### Make a Commit
1. Press `s` to go to Status tab
2. Use `Space` to stage files
3. Press `c` to commit
4. Type message and press `Enter`

### Switch Branches
1. Press `b` to go to Branches tab
2. Use `↑/↓` to select branch
3. Press `Enter` to checkout

### Push Changes
1. Press `p` to push
2. Confirm remote and branch

### View Git Graph
1. Press `g` to see commit graph
2. Use `↑/↓` to navigate commits
3. Press `Enter` to see details

## Editor Integration

### Neovim
```lua
-- Add to your config
require('gitflow').setup()

-- Keybinding: <leader>gg
```

### VSCode
```
Ctrl+Shift+P → "GitFlow: Open GitFlow TUI"
```

## Next Steps

- Read the [full documentation](docs/)
- Customize your [theme](docs/THEMES.md)
- Set up [authentication](docs/AUTH.md)
- Explore [advanced features](docs/ADVANCED.md)

## Need Help?

- `?` in the TUI shows all shortcuts
- Check [Troubleshooting](docs/INSTALLATION.md#troubleshooting)
- Open an issue on [GitHub](https://github.com/gitflow/tui/issues)

Happy coding! 🌿
