package auth

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// AuthMethod represents the authentication method
type AuthMethod string

const (
	SSH      AuthMethod = "ssh"
	HTTPS    AuthMethod = "https"
	Token    AuthMethod = "token"
	OAuth    AuthMethod = "oauth"
)

// Credential represents stored credentials
type Credential struct {
	Host       string     `json:"host"`
	Username   string     `json:"username"`
	Password   string     `json:"password"` // Encrypted
	Token      string     `json:"token"`    // For token-based auth
	Method     AuthMethod `json:"method"`
	SSHKeyPath string     `json:"ssh_key_path"`
}

// Manager handles authentication
type Manager struct {
	configDir string
	credsFile string
}

// New creates a new auth manager
func New() (*Manager, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	appDir := filepath.Join(configDir, "gitflow-tui")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return nil, err
	}

	return &Manager{
		configDir: appDir,
		credsFile: filepath.Join(appDir, "credentials"),
	}, nil
}

// LoadCredentials loads stored credentials
func (m *Manager) LoadCredentials() (map[string]Credential, error) {
	creds := make(map[string]Credential)

	data, err := os.ReadFile(m.credsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return creds, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return creds, nil
}

// SaveCredentials saves credentials
func (m *Manager) SaveCredentials(creds map[string]Credential) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	// Set restrictive permissions
	return os.WriteFile(m.credsFile, data, 0600)
}

// AddCredential adds a new credential
func (m *Manager) AddCredential(cred Credential) error {
	creds, err := m.LoadCredentials()
	if err != nil {
		return err
	}

	creds[cred.Host] = cred
	return m.SaveCredentials(creds)
}

// RemoveCredential removes a credential
func (m *Manager) RemoveCredential(host string) error {
	creds, err := m.LoadCredentials()
	if err != nil {
		return err
	}

	delete(creds, host)
	return m.SaveCredentials(creds)
}

// GetCredential gets a credential for a host
func (m *Manager) GetCredential(host string) (*Credential, error) {
	creds, err := m.LoadCredentials()
	if err != nil {
		return nil, err
	}

	if cred, ok := creds[host]; ok {
		return &cred, nil
	}

	return nil, fmt.Errorf("no credentials found for %s", host)
}

// SetupSSH sets up SSH authentication
func (m *Manager) SetupSSH(keyPath string) error {
	if keyPath == "" {
		home, _ := os.UserHomeDir()
		keyPath = filepath.Join(home, ".ssh", "id_rsa")
	}

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("SSH key not found at %s", keyPath)
	}

	// Test SSH connection
	cmd := exec.Command("ssh", "-T", "git@github.com")
	cmd.Env = append(os.Environ(), fmt.Sprintf("SSH_AUTH_SOCK=%s", os.Getenv("SSH_AUTH_SOCK")))
	
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "successfully authenticated") {
		return fmt.Errorf("SSH test failed: %s", string(output))
	}

	return nil
}

// GenerateSSHKey generates a new SSH key
func (m *Manager) GenerateSSHKey(email, keyPath string) error {
	if keyPath == "" {
		home, _ := os.UserHomeDir()
		keyPath = filepath.Join(home, ".ssh", "gitflow_tui")
	}

	// Ensure .ssh directory exists
	sshDir := filepath.Dir(keyPath)
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		return err
	}

	// Generate key
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-C", email, "-f", keyPath, "-N", "")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to generate SSH key: %s", string(output))
	}

	return nil
}

// GetSSHPublicKey returns the SSH public key
func (m *Manager) GetSSHPublicKey(keyPath string) (string, error) {
	pubKeyPath := keyPath + ".pub"
	data, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PromptPassword prompts for password securely
func PromptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// PromptInput prompts for input
func PromptInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ConfigureHTTPS configures HTTPS authentication
func (m *Manager) ConfigureHTTPS(host, username, password string) error {
	cred := Credential{
		Host:     host,
		Username: username,
		Password: password, // Should be encrypted in production
		Method:   HTTPS,
	}

	return m.AddCredential(cred)
}

// ConfigureToken configures token authentication
func (m *Manager) ConfigureToken(host, token string) error {
	cred := Credential{
		Host:   host,
		Token:  token, // Should be encrypted in production
		Method: Token,
	}

	return m.AddCredential(cred)
}

// GetAuthForRemote returns authentication for a remote URL
func (m *Manager) GetAuthForRemote(remoteURL string) (*Credential, error) {
	// Parse remote URL to extract host
	host := extractHostFromURL(remoteURL)
	if host == "" {
		return nil, fmt.Errorf("could not extract host from URL")
	}

	return m.GetCredential(host)
}

// extractHostFromURL extracts host from Git URL
func extractHostFromURL(url string) string {
	// Handle SSH format: git@github.com:user/repo.git
	if strings.HasPrefix(url, "git@") {
		parts := strings.Split(url, ":")
		if len(parts) >= 2 {
			host := strings.TrimPrefix(parts[0], "git@")
			return host
		}
	}

	// Handle HTTPS format: https://github.com/user/repo.git
	if strings.HasPrefix(url, "https://") {
		url = strings.TrimPrefix(url, "https://")
		parts := strings.Split(url, "/")
		if len(parts) >= 1 {
			return parts[0]
		}
	}

	return ""
}

// TestAuth tests authentication with a remote
func (m *Manager) TestAuth(remoteURL string) error {
	// Try to fetch from remote
	cmd := exec.Command("git", "ls-remote", remoteURL)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("authentication failed: %s", string(output))
	}
	return nil
}

// SetupGitCredentialHelper sets up Git credential helper
func (m *Manager) SetupGitCredentialHelper() error {
	// Configure Git to use the credential helper
	cmd := exec.Command("git", "config", "--global", "credential.helper", "cache")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to setup credential helper: %s", string(output))
	}
	return nil
}

// Providers for OAuth
type OAuthProvider string

const (
	GitHub    OAuthProvider = "github"
	GitLab    OAuthProvider = "gitlab"
	Bitbucket OAuthProvider = "bitbucket"
)

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	Provider    OAuthProvider
	ClientID    string
	Secret      string
	RedirectURL string
	Scopes      []string
}

// StartOAuthFlow starts OAuth authentication flow
func (m *Manager) StartOAuthFlow(config OAuthConfig) (string, error) {
	// This is a simplified version
	// In production, you'd implement the full OAuth flow
	switch config.Provider {
	case GitHub:
		return m.startGitHubOAuth(config)
	case GitLab:
		return m.startGitLabOAuth(config)
	default:
		return "", fmt.Errorf("unsupported OAuth provider: %s", config.Provider)
	}
}

func (m *Manager) startGitHubOAuth(config OAuthConfig) (string, error) {
	// GitHub OAuth URL
	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s",
		config.ClientID,
		config.RedirectURL,
		strings.Join(config.Scopes, "%20"),
	)

	fmt.Printf("Please visit this URL to authorize: %s\n", authURL)
	code, err := PromptInput("Enter the authorization code: ")
	if err != nil {
		return "", err
	}

	// Exchange code for token
	return m.exchangeGitHubToken(config, code)
}

func (m *Manager) startGitLabOAuth(config OAuthConfig) (string, error) {
	// Similar implementation for GitLab
	return "", fmt.Errorf("GitLab OAuth not yet implemented")
}

func (m *Manager) exchangeGitHubToken(config OAuthConfig, code string) (string, error) {
	// Make HTTP request to exchange code for token
	// This is a placeholder - implement actual token exchange
	return "", fmt.Errorf("token exchange not implemented")
}

// ListConfiguredHosts lists all configured authentication hosts
func (m *Manager) ListConfiguredHosts() ([]string, error) {
	creds, err := m.LoadCredentials()
	if err != nil {
		return nil, err
	}

	var hosts []string
	for host := range creds {
		hosts = append(hosts, host)
	}

	return hosts, nil
}
