# Project Structure

```
gitflow-tui/
│
├── 📁 cmd/gitflow-tui/          # Application entry point
│   └── main.go                  # Main function with ASCII banner
│
├── 📁 internal/                 # Internal packages
│   ├── 📁 git/                  # Git operations
│   │   └── commands.go          # All Git commands implementation
│   ├── 📁 ui/                   # Terminal UI
│   │   ├── model.go             # Bubble Tea model & views
│   │   └── commands.go          # UI command handlers
│   ├── 📁 config/               # Configuration
│   │   └── config.go            # Theme & settings
│   └── 📁 auth/                 # Authentication
│       └── auth.go              # SSH/HTTPS/Token/OAuth
│
├── 📁 pkg/                      # Public packages
│   └── 📁 graph/                # Graph visualization
│       └── graph.go             # ASCII/Unicode graph renderer
│
├── 📁 editors/                  # Editor integrations
│   ├── 📁 nvim/                 # Neovim plugin
│   │   └── 📁 lua/gitflow/
│   │       └── init.lua         # Lua plugin
│   └── 📁 vscode/               # VSCode extension
│       ├── 📁 src/
│       │   ├── extension.ts     # Extension entry
│       │   ├── terminal.ts      # Terminal handler
│       │   ├── provider.ts      # Tree view provider
│       │   └── statusbar.ts     # Status bar
│       ├── package.json         # Extension manifest
│       └── tsconfig.json        # TypeScript config
│
├── 📁 docs/                     # Documentation
│   └── INSTALLATION.md          # Installation guide
│
├── 📁 .github/workflows/        # CI/CD
│   └── release.yml              # GitHub Actions
│
├── 📄 go.mod                    # Go module
├── 📄 Makefile                  # Build automation
├── 📄 Dockerfile                # Container image
├── 📄 README.md                 # Main documentation
├── 📄 QUICKSTART.md             # Quick start guide
├── 📄 CONTRIBUTING.md           # Contribution guidelines
├── 📄 CHANGELOG.md              # Version history
├── 📄 LICENSE                   # MIT License
└── 📄 .gitignore                # Git ignore rules
```

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

## Color Scheme

```
┌─────────────────────────────────────────────────────────────┐
│                     GitFlow Theme                            │
├─────────────────────────────────────────────────────────────┤
│  Primary    #00D9A5  ████  Green   - Main actions           │
│  Secondary  #00B4A6  ████  Teal    - Secondary elements     │
│  Tertiary   #0091EA  ████  Blue    - Highlights             │
│  Accent     #00E5FF  ████  Firozi  - Accent color           │
│  Highlight  #FF6D00  ████  Orange  - Warnings               │
│  Background #0D1117  ████  Dark    - Background             │
│  Foreground #E6EDF3  ████  Light   - Text                   │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

```
User Input (Keyboard/Mouse)
    ↓
Bubble Tea Framework
    ↓
UI Model (Update/View)
    ↓
Git Operations
    ↓
Git CLI / go-git
    ↓
Repository
```

## Module Dependencies

```
cmd/gitflow-tui
    ├── internal/ui
    │   ├── internal/git
    │   ├── internal/config
    │   └── pkg/graph
    ├── internal/auth
    └── internal/config
```

## Editor Integration Flow

### Neovim
```
User Command (:GitFlow)
    ↓
Lua Plugin
    ↓
Terminal Buffer
    ↓
gitflow-tui Binary
```

### VSCode
```
User Command (Ctrl+Shift+P)
    ↓
Extension Host
    ↓
Terminal API
    ↓
gitflow-tui Binary
```

## Build Process

```
Source Code
    ↓
Go Compiler
    ↓
Binary (per platform)
    ↓
Package (tar.gz/zip)
    ↓
Release
```

## Testing Strategy

```
Unit Tests
    ├── internal/git/*_test.go
    ├── internal/ui/*_test.go
    └── pkg/graph/*_test.go

Integration Tests
    └── tests/integration/

E2E Tests
    └── tests/e2e/
```
