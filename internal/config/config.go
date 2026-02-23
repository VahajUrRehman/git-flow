package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ThemeColors defines the color palette
type ThemeColors struct {
	Primary     string // Green
	Secondary   string // Teal
	Tertiary    string // Blue
	Accent      string // Firozi/Cyan
	Highlight   string // Orange
	Background  string
	Foreground  string
	Success     string
	Warning     string
	Error       string
	Muted       string
	Border      string
}

// Default theme colors - Green Teal Blue Firozi Orange
type Theme struct {
	Name   string
	Colors ThemeColors
}

var DefaultTheme = Theme{
	Name: "gitflow",
	Colors: ThemeColors{
		Primary:     "#00D9A5", // Green
		Secondary:   "#00B4A6", // Teal
		Tertiary:    "#0091EA", // Blue
		Accent:      "#00E5FF", // Firozi/Cyan
		Highlight:   "#FF6D00", // Orange
		Background:  "#0D1117", // Dark background
		Foreground:  "#E6EDF3", // Light text
		Success:     "#3FB950", // Success green
		Warning:     "#FFA500", // Warning orange
		Error:       "#F85149", // Error red
		Muted:       "#8B949E", // Muted gray
		Border:      "#30363D", // Border color
	},
}

// Config holds all application configuration
type Config struct {
	Theme           Theme    `json:"theme"`
	GitPath         string   `json:"git_path"`
	Editor          string   `json:"editor"`
	DefaultBranch   string   `json:"default_branch"`
	ShowGraph       bool     `json:"show_graph"`
	GraphStyle      string   `json:"graph_style"` // ascii, unicode, compact
	MouseEnabled    bool     `json:"mouse_enabled"`
	Animations      bool     `json:"animations"`
	AuthMethod      string   `json:"auth_method"` // ssh, https, token
	RecentRepos     []string `json:"recent_repos"`
	MaxRecentRepos  int      `json:"max_recent_repos"`
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Theme:          DefaultTheme,
		GitPath:        "git",
		Editor:         os.Getenv("EDITOR"),
		DefaultBranch:  "main",
		ShowGraph:      true,
		GraphStyle:     "unicode",
		MouseEnabled:   true,
		Animations:     true,
		AuthMethod:     "ssh",
		RecentRepos:    []string{},
		MaxRecentRepos: 10,
	}
}

// Load loads configuration from file
func Load() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "gitflow-tui", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := Default()
			_ = cfg.Save()
			return cfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(configDir, "gitflow-tui")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(appDir, "config.json")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetConfigDir returns the configuration directory
func GetConfigDir() string {
	configDir, _ := os.UserConfigDir()
	return filepath.Join(configDir, "gitflow-tui")
}
