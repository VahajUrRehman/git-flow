package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Command represents a UI command
type Command struct {
	Name        string
	Description string
	Action      func(*Model) tea.Cmd
	Key         string
}

// AvailableCommands returns all available commands
func AvailableCommands() []Command {
	return []Command{
		{
			Name:        "commit",
			Description: "Create a new commit",
			Key:         "c",
			Action:      cmdCommit,
		},
		{
			Name:        "push",
			Description: "Push to remote",
			Key:         "p",
			Action:      cmdPush,
		},
		{
			Name:        "pull",
			Description: "Pull from remote",
			Key:         "P",
			Action:      cmdPull,
		},
		{
			Name:        "fetch",
			Description: "Fetch from remote",
			Key:         "f",
			Action:      cmdFetch,
		},
		{
			Name:        "checkout",
			Description: "Checkout branch",
			Key:         "b",
			Action:      cmdCheckout,
		},
		{
			Name:        "merge",
			Description: "Merge branch",
			Key:         "m",
			Action:      cmdMerge,
		},
		{
			Name:        "rebase",
			Description: "Rebase branch",
			Key:         "R",
			Action:      cmdRebase,
		},
		{
			Name:        "stash",
			Description: "Stash changes",
			Key:         "S",
			Action:      cmdStash,
		},
		{
			Name:        "pop",
			Description: "Pop stash",
			Key:         "O",
			Action:      cmdStashPop,
		},
		{
			Name:        "tag",
			Description: "Create tag",
			Key:         "t",
			Action:      cmdTag,
		},
		{
			Name:        "reset",
			Description: "Reset to commit",
			Key:         "X",
			Action:      cmdReset,
		},
		{
			Name:        "cherry-pick",
			Description: "Cherry-pick commit",
			Key:         "C",
			Action:      cmdCherryPick,
		},
	}
}

// cmdCommit handles commit command
func cmdCommit(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.inputMode = "commit"
		m.input.Placeholder = "Enter commit message..."
		m.input.SetValue("")
		m.input.Focus()
		m.currentView = ViewInput
		m.inputCallback = func(value string) {
			if value != "" {
				err := m.git.Commit(value, false)
				if err != nil {
					m.errorMsg = err.Error()
				} else {
					m.successMsg = "Committed: " + value
					m.loadData()
				}
			}
		}
		return nil
	}
}

// cmdPush handles push command
func cmdPush(m *Model) tea.Cmd {
	return func() tea.Msg {
		// Get current branch
		branch, err := m.git.GetCurrentBranch()
		if err != nil {
			m.errorMsg = err.Error()
			return nil
		}

		// Find origin remote
		var remote string
		for _, r := range m.remotes {
			if r.Name == "origin" {
				remote = r.Name
				break
			}
		}
		if remote == "" && len(m.remotes) > 0 {
			remote = m.remotes[0].Name
		}

		err = m.git.Push(remote, branch, false)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = fmt.Sprintf("Pushed to %s/%s", remote, branch)
		}
		return nil
	}
}

// cmdPull handles pull command
func cmdPull(m *Model) tea.Cmd {
	return func() tea.Msg {
		branch, err := m.git.GetCurrentBranch()
		if err != nil {
			m.errorMsg = err.Error()
			return nil
		}

		var remote string
		for _, r := range m.remotes {
			if r.Name == "origin" {
				remote = r.Name
				break
			}
		}
		if remote == "" && len(m.remotes) > 0 {
			remote = m.remotes[0].Name
		}

		err = m.git.Pull(remote, branch, false)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = fmt.Sprintf("Pulled from %s/%s", remote, branch)
			m.loadData()
		}
		return nil
	}
}

// cmdFetch handles fetch command
func cmdFetch(m *Model) tea.Cmd {
	return func() tea.Msg {
		var remote string
		for _, r := range m.remotes {
			if r.Name == "origin" {
				remote = r.Name
				break
			}
		}

		err := m.git.Fetch(remote)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			if remote != "" {
				m.successMsg = "Fetched from " + remote
			} else {
				m.successMsg = "Fetched all remotes"
			}
			m.loadData()
		}
		return nil
	}
}

// cmdCheckout handles checkout command
func cmdCheckout(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.inputMode = "checkout"
		m.input.Placeholder = "Enter branch name..."
		m.input.SetValue("")
		m.input.Focus()
		m.currentView = ViewInput
		m.inputCallback = func(value string) {
			if value != "" {
				// Check if branch exists
				exists := false
				for _, b := range m.branches {
					if b.Name == value {
						exists = true
						break
					}
				}

				err := m.git.Checkout(value, !exists)
				if err != nil {
					m.errorMsg = err.Error()
				} else {
					if exists {
						m.successMsg = "Switched to " + value
					} else {
						m.successMsg = "Created and switched to " + value
					}
					m.loadData()
				}
			}
		}
		return nil
	}
}

// cmdMerge handles merge command
func cmdMerge(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.selectedBranch >= len(m.branches) {
			m.errorMsg = "No branch selected"
			return nil
		}

		branch := m.branches[m.selectedBranch]
		err := m.git.Merge(branch.Name, false)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = "Merged " + branch.Name
			m.loadData()
		}
		return nil
	}
}

// cmdRebase handles rebase command
func cmdRebase(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.selectedBranch >= len(m.branches) {
			m.errorMsg = "No branch selected"
			return nil
		}

		branch := m.branches[m.selectedBranch]
		err := m.git.Rebase(branch.Name, false)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = "Rebased onto " + branch.Name
			m.loadData()
		}
		return nil
	}
}

// cmdStash handles stash command
func cmdStash(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.inputMode = "stash"
		m.input.Placeholder = "Enter stash message (optional)..."
		m.input.SetValue("")
		m.input.Focus()
		m.currentView = ViewInput
		m.inputCallback = func(value string) {
			err := m.git.StashSave(value)
			if err != nil {
				m.errorMsg = err.Error()
			} else {
				m.successMsg = "Changes stashed"
				m.loadData()
			}
		}
		return nil
	}
}

// cmdStashPop handles stash pop command
func cmdStashPop(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.selectedStash >= len(m.stashes) {
			m.errorMsg = "No stash selected"
			return nil
		}

		stash := m.stashes[m.selectedStash]
		err := m.git.StashPop(stash.Index)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = "Popped stash@{" + fmt.Sprintf("%d", stash.Index) + "}"
			m.loadData()
		}
		return nil
	}
}

// cmdTag handles tag command
func cmdTag(m *Model) tea.Cmd {
	return func() tea.Msg {
		m.inputMode = "tag"
		m.input.Placeholder = "Enter tag name..."
		m.input.SetValue("")
		m.input.Focus()
		m.currentView = ViewInput
		m.inputCallback = func(value string) {
			if value != "" {
				err := m.git.CreateTag(value, "")
				if err != nil {
					m.errorMsg = err.Error()
				} else {
					m.successMsg = "Created tag " + value
					m.loadData()
				}
			}
		}
		return nil
	}
}

// cmdReset handles reset command
func cmdReset(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.selectedCommit >= len(m.commits) {
			m.errorMsg = "No commit selected"
			return nil
		}

		commit := m.commits[m.selectedCommit]
		
		// Show confirmation
		m.inputMode = "reset"
		m.input.Placeholder = "Type 'soft', 'mixed', or 'hard' for reset mode..."
		m.input.SetValue("")
		m.input.Focus()
		m.currentView = ViewInput
		m.inputCallback = func(value string) {
			mode := strings.ToLower(value)
			if mode != "soft" && mode != "mixed" && mode != "hard" {
				m.errorMsg = "Invalid reset mode"
				return
			}

			err := m.git.Reset("--"+mode, commit.Hash)
			if err != nil {
				m.errorMsg = err.Error()
			} else {
				m.successMsg = fmt.Sprintf("Reset (%s) to %s", mode, commit.ShortHash)
				m.loadData()
			}
		}
		return nil
	}
}

// cmdCherryPick handles cherry-pick command
func cmdCherryPick(m *Model) tea.Cmd {
	return func() tea.Msg {
		if m.selectedCommit >= len(m.commits) {
			m.errorMsg = "No commit selected"
			return nil
		}

		commit := m.commits[m.selectedCommit]
		err := m.git.CherryPick(commit.Hash)
		if err != nil {
			m.errorMsg = err.Error()
		} else {
			m.successMsg = "Cherry-picked " + commit.ShortHash
			m.loadData()
		}
		return nil
	}
}

// ExecuteCommand executes a command by name
func (m *Model) ExecuteCommand(name string) tea.Cmd {
	for _, cmd := range AvailableCommands() {
		if cmd.Name == name {
			return cmd.Action(m)
		}
	}
	return nil
}

// GetCommandHelp returns help text for all commands
func GetCommandHelp() string {
	var lines []string
	lines = append(lines, "Available Commands:")
	lines = append(lines, "")

	for _, cmd := range AvailableCommands() {
		lines = append(lines, fmt.Sprintf("  %s - %s", cmd.Key, cmd.Description))
	}

	return strings.Join(lines, "\n")
}
