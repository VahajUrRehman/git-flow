package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Repository represents a Git repository
type Repository struct {
	Path       string
	Branch     string
	RemoteURL  string
	IsBare     bool
	Worktree   string
}

// Commit represents a Git commit
type Commit struct {
	Hash      string
	ShortHash string
	Message   string
	Author    string
	Email     string
	Date      time.Time
	Refs      []string
	Parents   []string
}

// Branch represents a Git branch
type Branch struct {
	Name       string
	Current    bool
	Remote     string
	Ahead      int
	Behind     int
	LastCommit time.Time
}

// Status represents the working tree status
type Status struct {
	Staged    []FileStatus
	Unstaged  []FileStatus
	Untracked []string
	Conflict  []string
}

// FileStatus represents a file's status
type FileStatus struct {
	Path   string
	Status string // M, A, D, R, C, U
	Score  int    // For rename/copy
}

// Remote represents a Git remote
type Remote struct {
	Name string
	URL  string
	Type string // fetch, push
}

// Stash represents a stash entry
type Stash struct {
	Index   int
	Message string
	Branch  string
}

// Tag represents a Git tag
type Tag struct {
	Name    string
	Message string
	Hash    string
}

// Git is the main Git operations handler
type Git struct {
	repoPath string
}

// New creates a new Git handler
func New(repoPath string) *Git {
	return &Git{repoPath: repoPath}
}

// FindRepository finds the Git repository starting from the given path
func FindRepository(startPath string) (*Repository, error) {
	path := startPath
	for {
		gitDir := filepath.Join(path, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			return &Repository{Path: path}, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return nil, fmt.Errorf("not a git repository")
		}
		path = parent
	}
}

// Execute runs a git command and returns output
func (g *Git) Execute(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoPath
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, errOut.String())
	}
	return out.String(), nil
}

// GetCurrentBranch returns the current branch name
func (g *Git) GetCurrentBranch() (string, error) {
	out, err := g.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// GetBranches returns all branches
func (g *Git) GetBranches() ([]Branch, error) {
	out, err := g.Execute("branch", "-vv")
	if err != nil {
		return nil, err
	}

	var branches []Branch
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		current := line[0] == '*'
		name := strings.TrimSpace(line[1:])
		parts := strings.Fields(name)
		if len(parts) == 0 {
			continue
		}

		branch := Branch{
			Name:    parts[0],
			Current: current,
		}

		// Parse ahead/behind info
		for _, part := range parts {
			if strings.Contains(part, "[") {
				info := strings.Trim(part, "[]")
				if strings.Contains(info, "ahead") {
					fmt.Sscanf(info, "ahead %d", &branch.Ahead)
				}
				if strings.Contains(info, "behind") {
					fmt.Sscanf(info, "behind %d", &branch.Behind)
				}
			}
		}

		branches = append(branches, branch)
	}

	return branches, nil
}

// GetCommits returns commit history
func (g *Git) GetCommits(limit int) ([]Commit, error) {
	format := "%H|%h|%s|%an|%ae|%ai|%D|%P"
	out, err := g.Execute("log", fmt.Sprintf("-%d", limit), fmt.Sprintf("--pretty=format:%s", format))
	if err != nil {
		return nil, err
	}

	var commits []Commit
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			continue
		}

		date, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[5])
		
		commit := Commit{
			Hash:      parts[0],
			ShortHash: parts[1],
			Message:   parts[2],
			Author:    parts[3],
			Email:     parts[4],
			Date:      date,
			Refs:      strings.Split(parts[6], ", "),
			Parents:   strings.Split(parts[7], " "),
		}

		commits = append(commits, commit)
	}

	return commits, nil
}

// GetStatus returns the working tree status
func (g *Git) GetStatus() (*Status, error) {
	out, err := g.Execute("status", "--porcelain", "-u")
	if err != nil {
		return nil, err
	}

	status := &Status{}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}

		staged := line[0]
		unstaged := line[1]
		path := line[3:]

		// Handle renamed files
		if strings.Contains(path, " -> ") {
			parts := strings.Split(path, " -> ")
			path = parts[1]
		}

		switch staged {
		case 'M', 'A', 'D', 'R', 'C':
			status.Staged = append(status.Staged, FileStatus{Path: path, Status: string(staged)})
		case 'U':
			status.Conflict = append(status.Conflict, path)
		}

		switch unstaged {
		case 'M', 'D':
			status.Unstaged = append(status.Unstaged, FileStatus{Path: path, Status: string(unstaged)})
		}

		if staged == '?' && unstaged == '?' {
			status.Untracked = append(status.Untracked, path)
		}
	}

	return status, nil
}

// GetRemotes returns all remotes
func (g *Git) GetRemotes() ([]Remote, error) {
	out, err := g.Execute("remote", "-v")
	if err != nil {
		return nil, err
	}

	var remotes []Remote
	seen := make(map[string]bool)
	
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		key := parts[0] + parts[2]
		if seen[key] {
			continue
		}
		seen[key] = true

		remotes = append(remotes, Remote{
			Name: parts[0],
			URL:  parts[1],
			Type: strings.Trim(parts[2], "()"),
		})
	}

	return remotes, nil
}

// GetStash returns stash list
func (g *Git) GetStash() ([]Stash, error) {
	out, err := g.Execute("stash", "list", "--format=%gd|%s")
	if err != nil {
		return nil, err
	}

	var stashes []Stash
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 2)
		if len(parts) < 2 {
			continue
		}

		var index int
		fmt.Sscanf(parts[0], "stash@{%d}", &index)
		
		stashes = append(stashes, Stash{
			Index:   index,
			Message: parts[1],
		})
	}

	return stashes, nil
}

// GetTags returns all tags
func (g *Git) GetTags() ([]Tag, error) {
	out, err := g.Execute("tag", "-l", "-n1")
	if err != nil {
		return nil, err
	}

	var tags []Tag
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		tag := Tag{Name: parts[0]}
		if len(parts) > 1 {
			tag.Message = strings.Join(parts[1:], " ")
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// Stage adds files to staging area
func (g *Git) Stage(paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := g.Execute(args...)
	return err
}

// Unstage removes files from staging area
func (g *Git) Unstage(paths ...string) error {
	args := append([]string{"reset", "HEAD"}, paths...)
	_, err := g.Execute(args...)
	return err
}

// Commit creates a new commit
func (g *Git) Commit(message string, amend bool) error {
	args := []string{"commit", "-m", message}
	if amend {
		args = append(args, "--amend")
	}
	_, err := g.Execute(args...)
	return err
}

// Push pushes to remote
func (g *Git) Push(remote, branch string, force bool) error {
	args := []string{"push", remote, branch}
	if force {
		args = append(args, "--force")
	}
	_, err := g.Execute(args...)
	return err
}

// Pull pulls from remote
func (g *Git) Pull(remote, branch string, rebase bool) error {
	args := []string{"pull", remote, branch}
	if rebase {
		args = append(args, "--rebase")
	}
	_, err := g.Execute(args...)
	return err
}

// Fetch fetches from remote
func (g *Git) Fetch(remote string) error {
	args := []string{"fetch"}
	if remote != "" {
		args = append(args, remote)
	}
	_, err := g.Execute(args...)
	return err
}

// Checkout switches branches
func (g *Git) Checkout(branch string, create bool) error {
	args := []string{"checkout"}
	if create {
		args = append(args, "-b")
	}
	args = append(args, branch)
	_, err := g.Execute(args...)
	return err
}

// Merge merges a branch
func (g *Git) Merge(branch string, noFF bool) error {
	args := []string{"merge"}
	if noFF {
		args = append(args, "--no-ff")
	}
	args = append(args, branch)
	_, err := g.Execute(args...)
	return err
}

// Rebase starts a rebase
func (g *Git) Rebase(branch string, interactive bool) error {
	args := []string{"rebase"}
	if interactive {
		args = append(args, "-i")
	}
	args = append(args, branch)
	_, err := g.Execute(args...)
	return err
}

// CherryPick cherry-picks a commit
func (g *Git) CherryPick(hash string) error {
	_, err := g.Execute("cherry-pick", hash)
	return err
}

// Reset resets to a commit
func (g *Git) Reset(mode, hash string) error {
	_, err := g.Execute("reset", mode, hash)
	return err
}

// Revert reverts a commit
func (g *Git) Revert(hash string, noEdit bool) error {
	args := []string{"revert"}
	if noEdit {
		args = append(args, "--no-edit")
	}
	args = append(args, hash)
	_, err := g.Execute(args...)
	return err
}

// StashSave saves changes to stash
func (g *Git) StashSave(message string) error {
	args := []string{"stash", "push"}
	if message != "" {
		args = append(args, "-m", message)
	}
	_, err := g.Execute(args...)
	return err
}

// StashPop pops stash
func (g *Git) StashPop(index int) error {
	_, err := g.Execute("stash", "pop", fmt.Sprintf("stash@{%d}", index))
	return err
}

// StashApply applies stash
func (g *Git) StashApply(index int) error {
	_, err := g.Execute("stash", "apply", fmt.Sprintf("stash@{%d}", index))
	return err
}

// StashDrop drops stash
func (g *Git) StashDrop(index int) error {
	_, err := g.Execute("stash", "drop", fmt.Sprintf("stash@{%d}", index))
	return err
}

// CreateTag creates a new tag
func (g *Git) CreateTag(name, message string) error {
	args := []string{"tag"}
	if message != "" {
		args = append(args, "-a", "-m", message)
	}
	args = append(args, name)
	_, err := g.Execute(args...)
	return err
}

// DeleteTag deletes a tag
func (g *Git) DeleteTag(name string) error {
	_, err := g.Execute("tag", "-d", name)
	return err
}

// GetDiff returns diff for files
func (g *Git) GetDiff(staged bool, paths ...string) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--cached")
	}
	args = append(args, paths...)
	return g.Execute(args...)
}

// GetLog returns formatted log
func (g *Git) GetLog(format string, limit int) (string, error) {
	return g.Execute("log", fmt.Sprintf("-%d", limit), fmt.Sprintf("--pretty=format:%s", format))
}
