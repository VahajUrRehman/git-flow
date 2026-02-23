package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gitflow/tui/internal/config"
	"github.com/gitflow/tui/internal/git"
	"github.com/gitflow/tui/pkg/graph"
)

// ViewState represents the current UI view
type ViewState int

const (
	ViewDashboard ViewState = iota
	ViewGraph
	ViewBranches
	ViewStatus
	ViewDiff
	ViewCommit
	ViewStash
	ViewRemote
	ViewTags
	ViewHelp
	ViewConfirm
	ViewInput
)

// Tab represents a tab in the UI
type Tab struct {
	Title    string
	View     ViewState
	Key      string
	Shortcut rune
}

// Tabs available in the UI
var Tabs = []Tab{
	{"Dashboard", ViewDashboard, "dashboard", 'd'},
	{"Graph", ViewGraph, "graph", 'g'},
	{"Branches", ViewBranches, "branches", 'b'},
	{"Status", ViewStatus, "status", 's'},
	{"Stash", ViewStash, "stash", 'S'},
	{"Remotes", ViewRemote, "remotes", 'r'},
	{"Tags", ViewTags, "tags", 't'},
}

// Model represents the main UI model
type Model struct {
	// Configuration
	config *config.Config

	// Git repository
	repo     *git.Repository
	git      *git.Git
	repoPath string

	// State
	currentView   ViewState
	activeTab     int
	width         int
	height        int
	ready         bool
	loading       bool
	errorMsg      string
	successMsg    string

	// Data
	commits   []git.Commit
	branches  []git.Branch
	status    *git.Status
	remotes   []git.Remote
	stashes   []git.Stash
	tags      []git.Tag
	currentBranch string

	// UI Components
	help         help.Model
	keys         keyMap
	viewport     viewport.Model
	list         list.Model
	input        textinput.Model
	textArea     textarea.Model
	commitList   list.Model
	branchList   list.Model
	fileList     list.Model
	stashList    list.Model
	remoteList   list.Model
	tagList      list.Model

	// Input state
	inputMode     string
	inputCallback func(string)
	confirmCallback func(bool)

	// Selection
	selectedCommit int
	selectedBranch int
	selectedFile   int
	selectedStash  int

	// Diff view
	diffContent   string
	diffStaged    bool

	// Graph
	graphRenderer *graph.Graph
}

// keyMap defines keyboard shortcuts
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	Enter    key.Binding
	Esc      key.Binding
	Space    key.Binding
	Refresh  key.Binding
}

// Default key bindings
var defaultKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "stage/unstage"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Space, k.Tab, k.ShiftTab},
		{k.Refresh, k.Esc, k.Help, k.Quit},
	}
}

// New creates a new UI model
func New(cfg *config.Config) *Model {
	// Find repository
	repoPath, _ := git.FindRepository(".")
	if repoPath == nil {
		// Try current directory
		repoPath = &git.Repository{Path: "."}
	}

	g := git.New(repoPath.Path)

	// Initialize input
	input := textinput.New()
	input.Placeholder = "Type here..."
	input.CharLimit = 256

	// Initialize textarea for commit messages
	ta := textarea.New()
	ta.Placeholder = "Enter commit message..."
	ta.SetWidth(50)
	ta.SetHeight(5)

	// Initialize help
	h := help.New()

	return &Model{
		config:     cfg,
		repo:       repoPath,
		git:        g,
		repoPath:   repoPath.Path,
		currentView: ViewDashboard,
		help:       h,
		keys:       defaultKeys,
		input:      input,
		textArea:   ta,
	}
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadData(),
		tea.EnterAltScreen,
	)
}

// loadData loads all git data
func (m *Model) loadData() tea.Cmd {
	return func() tea.Msg {
		var err error
		
		// Load commits
		m.commits, err = m.git.GetCommits(50)
		if err != nil {
			return errMsg{err: err}
		}

		// Load branches
		m.branches, err = m.git.GetBranches()
		if err != nil {
			return errMsg{err: err}
		}

		// Load status
		m.status, err = m.git.GetStatus()
		if err != nil {
			return errMsg{err: err}
		}

		// Load remotes
		m.remotes, err = m.git.GetRemotes()
		if err != nil {
			return errMsg{err: err}
		}

		// Load stashes
		m.stashes, err = m.git.GetStash()
		if err != nil {
			return errMsg{err: err}
		}

		// Load tags
		m.tags, err = m.git.GetTags()
		if err != nil {
			return errMsg{err: err}
		}

		// Get current branch
		m.currentBranch, _ = m.git.GetCurrentBranch()

		return dataLoadedMsg{}
	}
}

// Message types
type errMsg struct {
	err error
}

type dataLoadedMsg struct{}
type refreshMsg struct{}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 10
		m.ready = true

	case tea.MouseMsg:
		if m.config.MouseEnabled {
			return m.handleMouse(msg)
		}

	case tea.KeyMsg:
		return m.handleKey(msg)

	case errMsg:
		m.errorMsg = msg.err.Error()
		m.loading = false

	case dataLoadedMsg:
		m.loading = false
		m.updateLists()

	case refreshMsg:
		return m, m.loadData()
	}

	return m, nil
}

// handleMouse handles mouse events
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.MouseLeft:
		// Check if clicked on a tab
		if msg.Y == 1 {
			tabWidth := m.width / len(Tabs)
			clickedTab := msg.X / tabWidth
			if clickedTab < len(Tabs) {
				m.activeTab = clickedTab
				m.currentView = Tabs[clickedTab].View
			}
		}
	}
	return m, nil
}

// handleKey handles keyboard events
func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, m.keys.Help):
		if m.currentView == ViewHelp {
			m.currentView = ViewDashboard
		} else {
			m.currentView = ViewHelp
		}

	case key.Matches(msg, m.keys.Tab):
		m.activeTab = (m.activeTab + 1) % len(Tabs)
		m.currentView = Tabs[m.activeTab].View

	case key.Matches(msg, m.keys.ShiftTab):
		m.activeTab = (m.activeTab - 1 + len(Tabs)) % len(Tabs)
		m.currentView = Tabs[m.activeTab].View

	case key.Matches(msg, m.keys.Refresh):
		return m, m.loadData()

	case key.Matches(msg, m.keys.Esc):
		if m.currentView == ViewInput || m.currentView == ViewConfirm {
			m.currentView = ViewDashboard
		}

	default:
		// View-specific key handling
		switch m.currentView {
		case ViewGraph:
			return m.handleGraphKeys(msg)
		case ViewBranches:
			return m.handleBranchKeys(msg)
		case ViewStatus:
			return m.handleStatusKeys(msg)
		case ViewInput:
			return m.handleInputKeys(msg)
		}
	}

	return m, nil
}

// handleGraphKeys handles graph view keys
func (m *Model) handleGraphKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.selectedCommit > 0 {
			m.selectedCommit--
		}
	case key.Matches(msg, m.keys.Down):
		if m.selectedCommit < len(m.commits)-1 {
			m.selectedCommit++
		}
	case key.Matches(msg, m.keys.Enter):
		if m.selectedCommit < len(m.commits) {
			commit := m.commits[m.selectedCommit]
			m.showCommitDetails(commit)
		}
	}
	return m, nil
}

// handleBranchKeys handles branch view keys
func (m *Model) handleBranchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Up):
		if m.selectedBranch > 0 {
			m.selectedBranch--
		}
	case key.Matches(msg, m.keys.Down):
		if m.selectedBranch < len(m.branches)-1 {
			m.selectedBranch++
		}
	case key.Matches(msg, m.keys.Enter):
		if m.selectedBranch < len(m.branches) {
			branch := m.branches[m.selectedBranch]
			m.showBranchMenu(branch)
		}
	}
	return m, nil
}

// handleStatusKeys handles status view keys
func (m *Model) handleStatusKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Space):
		// Stage/unstage file
		m.toggleStage()
	case key.Matches(msg, m.keys.Enter):
		// Show diff
		m.showFileDiff()
	}
	return m, nil
}

// handleInputKeys handles input view keys
func (m *Model) handleInputKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		if m.inputCallback != nil {
			m.inputCallback(m.input.Value())
			m.input.SetValue("")
			m.currentView = ViewDashboard
		}
	case tea.KeyEsc:
		m.currentView = ViewDashboard
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

// updateLists updates all list components
func (m *Model) updateLists() {
	// Update commit list
	var commitItems []list.Item
	for _, c := range m.commits {
		commitItems = append(commitItems, commitItem{commit: c})
	}
	m.commitList.SetItems(commitItems)

	// Update branch list
	var branchItems []list.Item
	for _, b := range m.branches {
		branchItems = append(branchItems, branchItem{branch: b})
	}
	m.branchList.SetItems(branchItems)

	// Update file list
	var fileItems []list.Item
	if m.status != nil {
		for _, f := range m.status.Staged {
			fileItems = append(fileItems, fileItem{path: f.Path, status: f.Status, staged: true})
		}
		for _, f := range m.status.Unstaged {
			fileItems = append(fileItems, fileItem{path: f.Path, status: f.Status, staged: false})
		}
		for _, f := range m.status.Untracked {
			fileItems = append(fileItems, fileItem{path: f, status: "?", staged: false})
		}
	}
	m.fileList.SetItems(fileItems)

	// Update stash list
	var stashItems []list.Item
	for _, s := range m.stashes {
		stashItems = append(stashItems, stashItem{stash: s})
	}
	m.stashList.SetItems(stashItems)

	// Update remote list
	var remoteItems []list.Item
	for _, r := range m.remotes {
		remoteItems = append(remoteItems, remoteItem{remote: r})
	}
	m.remoteList.SetItems(remoteItems)

	// Update tag list
	var tagItems []list.Item
	for _, t := range m.tags {
		tagItems = append(tagItems, tagItem{tag: t})
	}
	m.tagList.SetItems(tagItems)
}

// View renders the UI
func (m *Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var sections []string

	// Banner
	sections = append(sections, m.renderBanner())

	// Tabs
	sections = append(sections, m.renderTabs())

	// Main content
	sections = append(sections, m.renderContent())

	// Status bar
	sections = append(sections, m.renderStatusBar())

	// Help
	sections = append(sections, m.renderHelp())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderBanner renders the ASCII banner
func (m *Model) renderBanner() string {
	banner := `
╔══════════════════════════════════════════════════════════════════╗
║  🌿 GitFlow TUI - Complete Git Management                        ║
║  Branch: ` + m.currentBranch + strings.Repeat(" ", max(0, 50-len(m.currentBranch))) + `║
╚══════════════════════════════════════════════════════════════════╝`

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.config.Theme.Colors.Primary)).
		Bold(true).
		Render(banner)
}

// renderTabs renders the tab bar
func (m *Model) renderTabs() string {
	var tabs []string
	for i, tab := range Tabs {
		style := lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder())

		if i == m.activeTab {
			style = style.
				Background(lipgloss.Color(m.config.Theme.Colors.Primary)).
				Foreground(lipgloss.Color(m.config.Theme.Colors.Background))
		} else {
			style = style.
				Foreground(lipgloss.Color(m.config.Theme.Colors.Muted))
		}

		tabs = append(tabs, style.Render(tab.Title))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, tabs...)
}

// renderContent renders the main content area
func (m *Model) renderContent() string {
	switch m.currentView {
	case ViewDashboard:
		return m.renderDashboard()
	case ViewGraph:
		return m.renderGraph()
	case ViewBranches:
		return m.renderBranches()
	case ViewStatus:
		return m.renderStatus()
	case ViewStash:
		return m.renderStash()
	case ViewRemote:
		return m.renderRemotes()
	case ViewTags:
		return m.renderTags()
	case ViewHelp:
		return m.renderHelpView()
	case ViewInput:
		return m.renderInput()
	case ViewDiff:
		return m.renderDiff()
	default:
		return m.renderDashboard()
	}
}

// renderDashboard renders the dashboard view
func (m *Model) renderDashboard() string {
	var sections []string

	// Repository info
	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	info := fmt.Sprintf("Repository: %s\nBranch: %s\nCommits: %d",
		m.repoPath, m.currentBranch, len(m.commits))
	sections = append(sections, infoStyle.Render(info))

	// Recent commits
	commitStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var commits string
	for i, c := range m.commits {
		if i >= 5 {
			break
		}
		commits += fmt.Sprintf("%s %s\n", c.ShortHash, c.Message)
	}
	sections = append(sections, commitStyle.Render("Recent Commits:\n"+commits))

	// Working tree status
	statusStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var status string
	if m.status != nil {
		status = fmt.Sprintf("Staged: %d\nUnstaged: %d\nUntracked: %d",
			len(m.status.Staged), len(m.status.Unstaged), len(m.status.Untracked))
	}
	sections = append(sections, statusStyle.Render("Working Tree:\n"+status))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderGraph renders the commit graph
func (m *Model) renderGraph() string {
	if len(m.commits) == 0 {
		return "No commits found"
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var graphStyle graph.GraphStyle
	switch m.config.GraphStyle {
	case "ascii":
		graphStyle = graph.ASCII
	case "compact":
		graphStyle = graph.Compact
	default:
		graphStyle = graph.Unicode
	}

	g := graph.New(m.commits, graphStyle)
	g.SetWidth(m.width - 4)
	
	return style.Render(g.Render())
}

// renderBranches renders the branches view
func (m *Model) renderBranches() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var content strings.Builder
	for i, b := range m.branches {
		prefix := "  "
		if b.Current {
			prefix = lipgloss.NewStyle().
				Foreground(lipgloss.Color(m.config.Theme.Colors.Accent)).
				Render("● ")
		}

		line := prefix + b.Name
		if b.Ahead > 0 || b.Behind > 0 {
			line += fmt.Sprintf(" [%d↑ %d↓]", b.Ahead, b.Behind)
		}

		if i == m.selectedBranch {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color(m.config.Theme.Colors.Tertiary)).
				Render(line)
		}

		content.WriteString(line + "\n")
	}

	return style.Render(content.String())
}

// renderStatus renders the status view
func (m *Model) renderStatus() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var content strings.Builder

	// Staged files
	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.config.Theme.Colors.Success)).
		Bold(true).
		Render("Staged:\n"))
	if m.status != nil {
		for _, f := range m.status.Staged {
			content.WriteString(fmt.Sprintf("  + %s\n", f.Path))
		}
	}
	content.WriteString("\n")

	// Unstaged files
	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.config.Theme.Colors.Warning)).
		Bold(true).
		Render("Unstaged:\n"))
	if m.status != nil {
		for _, f := range m.status.Unstaged {
			content.WriteString(fmt.Sprintf("  ~ %s\n", f.Path))
		}
	}
	content.WriteString("\n")

	// Untracked files
	content.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(m.config.Theme.Colors.Muted)).
		Bold(true).
		Render("Untracked:\n"))
	if m.status != nil {
		for _, f := range m.status.Untracked {
			content.WriteString(fmt.Sprintf("  ? %s\n", f))
		}
	}

	return style.Render(content.String())
}

// renderStash renders the stash view
func (m *Model) renderStash() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var content strings.Builder
	for _, s := range m.stashes {
		content.WriteString(fmt.Sprintf("stash@{%d}: %s\n", s.Index, s.Message))
	}

	return style.Render(content.String())
}

// renderRemotes renders the remotes view
func (m *Model) renderRemotes() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var content strings.Builder
	for _, r := range m.remotes {
		content.WriteString(fmt.Sprintf("%s\n  %s (%s)\n\n", r.Name, r.URL, r.Type))
	}

	return style.Render(content.String())
}

// renderTags renders the tags view
func (m *Model) renderTags() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	var content strings.Builder
	for _, t := range m.tags {
		content.WriteString(fmt.Sprintf("%s\n  %s\n\n", t.Name, t.Message))
	}

	return style.Render(content.String())
}

// renderHelpView renders the help view
func (m *Model) renderHelpView() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	help := `
Keyboard Shortcuts:

Navigation:
  ↑/k      Move up
  ↓/j      Move down
  ←/h      Move left
  →/l      Move right
  Tab      Next tab
  Shift+Tab Previous tab

Actions:
  Enter    Select/Confirm
  Space    Stage/Unstage file
  r        Refresh data
  ?        Toggle help
  q        Quit

Git Commands:
  c        Commit
  p        Push
  P        Pull
  f        Fetch
  b        Checkout branch
  m        Merge
  R        Rebase
`

	return style.Render(help)
}

// renderInput renders the input view
func (m *Model) renderInput() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	return style.Render(m.input.View())
}

// renderDiff renders the diff view
func (m *Model) renderDiff() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(m.config.Theme.Colors.Border)).
		Padding(1)

	return style.Render(m.diffContent)
}

// renderStatusBar renders the status bar
func (m *Model) renderStatusBar() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color(m.config.Theme.Colors.Border)).
		Foreground(lipgloss.Color(m.config.Theme.Colors.Foreground)).
		Padding(0, 1)

	status := fmt.Sprintf(" %s | %d commits | %d branches ",
		m.repoPath, len(m.commits), len(m.branches))

	if m.errorMsg != "" {
		status = " Error: " + m.errorMsg
		style = style.Background(lipgloss.Color(m.config.Theme.Colors.Error))
	} else if m.successMsg != "" {
		status = " " + m.successMsg
		style = style.Background(lipgloss.Color(m.config.Theme.Colors.Success))
	}

	return style.Render(status)
}

// renderHelp renders the help bar
func (m *Model) renderHelp() string {
	return m.help.View(m.keys)
}

// Helper methods
func (m *Model) showCommitDetails(commit git.Commit) {
	// Show commit details in a modal or new view
	m.successMsg = fmt.Sprintf("Selected: %s - %s", commit.ShortHash, commit.Message)
}

func (m *Model) showBranchMenu(branch git.Branch) {
	// Show branch actions menu
	m.successMsg = fmt.Sprintf("Branch: %s", branch.Name)
}

func (m *Model) toggleStage() {
	// Toggle stage/unstage for selected file
}

func (m *Model) showFileDiff() {
	// Show diff for selected file
	m.currentView = ViewDiff
}

// List item types
type commitItem struct {
	commit git.Commit
}

func (i commitItem) FilterValue() string { return i.commit.Message }

type branchItem struct {
	branch git.Branch
}

func (i branchItem) FilterValue() string { return i.branch.Name }

type fileItem struct {
	path   string
	status string
	staged bool
}

func (i fileItem) FilterValue() string { return i.path }

type stashItem struct {
	stash git.Stash
}

func (i stashItem) FilterValue() string { return i.stash.Message }

type remoteItem struct {
	remote git.Remote
}

func (i remoteItem) FilterValue() string { return i.remote.Name }

type tagItem struct {
	tag git.Tag
}

func (i tagItem) FilterValue() string { return i.tag.Name }

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
