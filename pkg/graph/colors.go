package graph

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/gitflow/tui/internal/config"
	"github.com/gitflow/tui/internal/git"
)

// ColoredGraph renders a colorful git commit graph
type ColoredGraph struct {
	commits []git.Commit
	style   GraphStyle
	width   int
	colors  config.ThemeColors
}

// NewColored creates a new colored graph
func NewColored(commits []git.Commit, style GraphStyle, colors config.ThemeColors) *ColoredGraph {
	return &ColoredGraph{
		commits: commits,
		style:   style,
		width:   80,
		colors:  colors,
	}
}

// SetWidth sets the graph width
func (g *ColoredGraph) SetWidth(width int) {
	g.width = width
}

// Render renders the colorful commit graph
func (g *ColoredGraph) Render() string {
	if len(g.commits) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(g.colors.Muted)).
			Render("No commits to display")
	}

	var lines []string
	for i, commit := range g.commits {
		line := g.renderColoredLine(commit, i)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (g *ColoredGraph) renderColoredLine(commit git.Commit, index int) string {
	// Color palette for branches (cycles through theme colors)
	branchColors := []string{
		g.colors.Primary,   // Green
		g.colors.Secondary, // Teal
		g.colors.Tertiary,  // Blue
		g.colors.Accent,    // Firozi
		g.colors.Highlight, // Orange
	}

	// Get color for this commit based on index
	graphColor := branchColors[index%len(branchColors)]

	// Styles
	graphStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(graphColor))
	hashStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(g.colors.Tertiary))
	msgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(g.colors.Foreground))
	refStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(g.colors.Highlight)).
		Bold(true)
	mutedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(g.colors.Muted))

	// Build graph part with unicode characters
	indent := strings.Repeat("  ", index%4)
	connector := "│"
	commitDot := "●"
	
	if index == 0 {
		connector = " "
	}

	// Graph line
	graphPart := graphStyle.Render(indent + connector + "─" + commitDot + "─")

	// Build message part
	var parts []string

	// Add refs (branches/tags) with highlight color
	if len(commit.Refs) > 0 {
		var refs []string
		for _, ref := range commit.Refs {
			if ref != "HEAD" && ref != "" {
				refs = append(refs, ref)
			}
		}
		if len(refs) > 0 {
			refStr := fmt.Sprintf("(%s)", strings.Join(refs, ", "))
			parts = append(parts, refStyle.Render(refStr))
		}
	}

	// Add hash
	parts = append(parts, hashStyle.Render(commit.ShortHash))

	// Add message
	parts = append(parts, msgStyle.Render(commit.Message))

	// Add author and date in muted style
	meta := mutedStyle.Render(fmt.Sprintf("<%s> %s", commit.Author, commit.Date.Format("Jan 2")))

	// Combine all parts
	msgPart := strings.Join(parts, " ")
	
	// Full line
	result := graphPart + " " + msgPart + " " + meta

	// Truncate if too long
	if lipgloss.Width(result) > g.width {
		available := g.width - lipgloss.Width(graphPart+" "+meta) - 4
		if available > 10 {
			truncated := msgPart[:available] + "..."
			result = graphPart + " " + truncated + " " + meta
		}
	}

	return result
}

// RenderCompact renders a compact colorful graph
func (g *ColoredGraph) RenderCompact() string {
	if len(g.commits) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(g.colors.Muted)).
			Render("No commits")
	}

	var lines []string
	colors := []string{
		g.colors.Primary,
		g.colors.Secondary,
		g.colors.Tertiary,
		g.colors.Accent,
	}

	for i, commit := range g.commits {
		color := colors[i%len(colors)]
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		
		line := fmt.Sprintf("%s %s %s",
			style.Render("●"),
			lipgloss.NewStyle().Foreground(lipgloss.Color(g.colors.Tertiary)).Render(commit.ShortHash),
			lipgloss.NewStyle().Foreground(lipgloss.Color(g.colors.Foreground)).Render(commit.Message))
		
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// RenderBranchGraph renders colorful branch list
func (g *ColoredGraph) RenderBranchGraph(branches []git.Branch, currentBranch string) string {
	if len(branches) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(g.colors.Muted)).
			Render("No branches")
	}

	var lines []string
	
	for _, branch := range branches {
		var line string
		
		if branch.Current {
			// Current branch - highlighted
			currentStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(g.colors.Highlight)).
				Bold(true)
			line = currentStyle.Render("● " + branch.Name)
		} else {
			// Other branches - muted
			otherStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(g.colors.Muted))
			line = otherStyle.Render("  " + branch.Name)
		}

		// Add ahead/behind info
		if branch.Ahead > 0 || branch.Behind > 0 {
			infoParts := []string{}
			if branch.Ahead > 0 {
				aheadStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(g.colors.Success))
				infoParts = append(infoParts, aheadStyle.Render(fmt.Sprintf("↑%d", branch.Ahead)))
			}
			if branch.Behind > 0 {
				behindStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color(g.colors.Error))
				infoParts = append(infoParts, behindStyle.Render(fmt.Sprintf("↓%d", branch.Behind)))
			}
			line += " " + strings.Join(infoParts, " ")
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// RenderStatusGraph renders colorful file status
func RenderStatusGraph(status *git.Status, colors config.ThemeColors) string {
	if status == nil {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Muted)).
			Render("No status available")
	}

	var lines []string
	
	// Staged files - Green
	if len(status.Staged) > 0 {
		stagedHeader := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Success)).
			Bold(true).
			Render(fmt.Sprintf("Staged (%d):", len(status.Staged)))
		lines = append(lines, stagedHeader)
		
		for _, f := range status.Staged {
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Success)).
				Render(fmt.Sprintf("  + %s [%s]", f.Path, f.Status))
			lines = append(lines, line)
		}
	}

	// Unstaged files - Orange
	if len(status.Unstaged) > 0 {
		unstagedHeader := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Highlight)).
			Bold(true).
			Render(fmt.Sprintf("\nUnstaged (%d):", len(status.Unstaged)))
		lines = append(lines, unstagedHeader)
		
		for _, f := range status.Unstaged {
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Highlight)).
				Render(fmt.Sprintf("  ~ %s [%s]", f.Path, f.Status))
			lines = append(lines, line)
		}
	}

	// Untracked files - Muted
	if len(status.Untracked) > 0 {
		untrackedHeader := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Muted)).
			Bold(true).
			Render(fmt.Sprintf("\nUntracked (%d):", len(status.Untracked)))
		lines = append(lines, untrackedHeader)
		
		for _, f := range status.Untracked {
			line := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Muted)).
				Render(fmt.Sprintf("  ? %s", f))
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Success)).
			Render("✓ Working tree clean")
	}

	return strings.Join(lines, "\n")
}
