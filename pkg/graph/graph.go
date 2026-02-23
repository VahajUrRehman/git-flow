package graph

import (
	"fmt"
	"strings"

	"github.com/gitflow/tui/internal/git"
)

// GraphStyle represents the style of graph rendering
type GraphStyle int

const (
	ASCII GraphStyle = iota
	Unicode
	Compact
	Detailed
)

// branchInfo tracks branch column and activity
type branchInfo struct {
	column int
	active bool
}

// Graph represents a Git commit graph
type Graph struct {
	commits []git.Commit
	style   GraphStyle
	width   int
}

// New creates a new graph
func New(commits []git.Commit, style GraphStyle) *Graph {
	return &Graph{
		commits: commits,
		style:   style,
		width:   80,
	}
}

// SetWidth sets the graph width
func (g *Graph) SetWidth(width int) {
	g.width = width
}

// GraphChar represents characters used for drawing
type GraphChar struct {
	Vertical   string
	Horizontal string
	Corner     string
	Branch     string
	Merge      string
	Commit     string
	Space      string
}

// CharSets for different styles
var (
	UnicodeChars = GraphChar{
		Vertical:   "│",
		Horizontal: "─",
		Corner:     "└",
		Branch:     "├",
		Merge:      "◉",
		Commit:     "●",
		Space:      " ",
	}
	ASCIIChars = GraphChar{
		Vertical:   "|",
		Horizontal: "-",
		Corner:     "\\",
		Branch:     "|",
		Merge:      "*",
		Commit:     "o",
		Space:      " ",
	}
	CompactChars = GraphChar{
		Vertical:   " ",
		Horizontal: "-",
		Corner:     " ",
		Branch:     "|",
		Merge:      "*",
		Commit:     "•",
		Space:      " ",
	}
)

// Render renders the commit graph
func (g *Graph) Render() string {
	var chars GraphChar
	switch g.style {
	case ASCII:
		chars = ASCIIChars
	case Compact:
		chars = CompactChars
	default:
		chars = UnicodeChars
	}

	if len(g.commits) == 0 {
		return "No commits to display"
	}

	// Build commit map for parent lookup
	commitMap := make(map[string]int)
	for i, c := range g.commits {
		commitMap[c.Hash] = i
	}

	// Track active branches
	branches := make(map[string]*branchInfo)
	nextColumn := 0

	var lines []string
	for i, commit := range g.commits {
		line := g.renderCommitLine(commit, i, branches, &nextColumn, commitMap, chars)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (g *Graph) renderCommitLine(commit git.Commit, index int, branches map[string]*branchInfo, nextColumn *int, commitMap map[string]int, chars GraphChar) string {
	var graphPart strings.Builder
	var msgPart strings.Builder

	// Determine column for this commit
	var col int
	if info, exists := branches[commit.Hash]; exists {
		col = info.column
	} else {
		col = *nextColumn
		branches[commit.Hash] = &branchInfo{column: col, active: true}
		*nextColumn++
	}

	// Build graph visualization
	columns := make([]string, *nextColumn+1)
	for i := range columns {
		columns[i] = chars.Space
	}

	// Mark current commit position
	columns[col] = chars.Commit

	// Draw connections to parents
	for _, parentHash := range commit.Parents {
		if parentHash == "" {
			continue
		}
		if parentInfo, exists := branches[parentHash]; exists {
			// Parent already has a column
			if parentInfo.column > col {
				// Draw horizontal line
				for i := col + 1; i < parentInfo.column; i++ {
					if columns[i] == chars.Space {
						columns[i] = chars.Horizontal
					}
				}
				columns[parentInfo.column] = chars.Corner
			}
		} else {
			// Create new branch for parent
			parentCol := *nextColumn
			branches[parentHash] = &branchInfo{column: parentCol, active: true}
			*nextColumn++
			
			// Extend columns
			for len(columns) <= parentCol {
				columns = append(columns, chars.Space)
			}
			
			// Draw branch line
			for i := col + 1; i < parentCol; i++ {
				if columns[i] == chars.Space {
					columns[i] = chars.Horizontal
				}
			}
			columns[parentCol] = chars.Corner
		}
	}

	// Build graph string
	for _, c := range columns {
		graphPart.WriteString(c)
		graphPart.WriteString(chars.Space)
	}

	// Build message part with refs
	refs := ""
	if len(commit.Refs) > 0 && commit.Refs[0] != "" {
		refList := []string{}
		for _, ref := range commit.Refs {
			if ref != "HEAD" && ref != "" {
				refList = append(refList, ref)
			}
		}
		if len(refList) > 0 {
			refs = "(" + strings.Join(refList, ", ") + ") "
		}
	}

	msgPart.WriteString(refs)
	msgPart.WriteString(commit.ShortHash)
	msgPart.WriteString(" ")
	msgPart.WriteString(commit.Message)

	// Combine parts
	result := graphPart.String() + " " + msgPart.String()
	
	// Truncate if too long
	if len(result) > g.width {
		result = result[:g.width-3] + "..."
	}

	return result
}

// RenderWithDetails renders detailed graph with author and date
func (g *Graph) RenderWithDetails() string {
	var lines []string
	chars := UnicodeChars

	for _, commit := range g.commits {
		// Main commit line
		commitLine := fmt.Sprintf("%s %s %s",
			chars.Commit,
			commit.ShortHash,
			commit.Message)

		// Details line
		details := fmt.Sprintf("   Author: %s <%s>", commit.Author, commit.Email)
		date := fmt.Sprintf("   Date: %s", commit.Date.Format("Mon Jan 2 15:04:05 2006"))

		lines = append(lines, commitLine, details, date, "")
	}

	return strings.Join(lines, "\n")
}

// RenderBranchGraph renders a simplified branch-only graph
func (g *Graph) RenderBranchGraph(branches []git.Branch, currentBranch string) string {
	var lines []string
	chars := UnicodeChars

	for _, branch := range branches {
		prefix := chars.Space + chars.Space
		if branch.Current {
			prefix = chars.Commit + chars.Space
		}

		line := prefix + branch.Name
		
		if branch.Ahead > 0 || branch.Behind > 0 {
			line += fmt.Sprintf(" [%d ahead, %d behind]", branch.Ahead, branch.Behind)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// ColorizeGraph adds color codes to graph output (for terminal)
func ColorizeGraph(graph string, colorFunc func(string, string) string) string {
	lines := strings.Split(graph, "\n")
	var result []string

	for _, line := range lines {
		// Color commit hashes
		if len(line) > 10 {
			hashStart := strings.IndexFunc(line, func(r rune) bool {
				return r == '●' || r == 'o' || r == '•'
			})
			if hashStart >= 0 && hashStart+8 < len(line) {
				potentialHash := strings.TrimSpace(line[hashStart+1 : hashStart+9])
				if len(potentialHash) == 7 {
					// This looks like a short hash
					colored := line[:hashStart+1] + 
						colorFunc(potentialHash, "#00E5FF") + 
						line[hashStart+9:]
					result = append(result, colored)
					continue
				}
			}
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}
