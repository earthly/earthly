package sshutil

import (
	"os"
	"path/filepath"
	"testing"
)

// writeSSHConfig creates a temp HOME directory, writes configContent to
// ~/.ssh/config, sets HOME and clears SSH_AUTH_SOCK for the test. Returns the
// tmpDir path so callers can construct expected expanded paths.
func writeSSHConfig(t *testing.T, configContent string) string {
	t.Helper()
	tmpDir := t.TempDir()
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sshDir, "config"), []byte(configContent), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", tmpDir)
	t.Setenv("SSH_AUTH_SOCK", "")
	return tmpDir
}

func TestExpandPath(t *testing.T) {
	home := "/home/testuser"
	tests := []struct {
		input    string
		expected string
	}{
		{"~/some/path", filepath.Join(home, "some/path")},
		{"~", home},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}
	for _, tt := range tests {
		got := expandPath(tt.input, home)
		if got != tt.expected {
			t.Errorf("expandPath(%q, %q) = %q, want %q", tt.input, home, got, tt.expected)
		}
	}
}

func TestExpandPathEnvVar(t *testing.T) {
	t.Setenv("TEST_SSH_VAR", "myvalue")
	home := "/home/testuser"
	got := expandPath("$TEST_SSH_VAR/agent.sock", home)
	if got != "myvalue/agent.sock" {
		t.Errorf("expandPath with env var = %q, want %q", got, "myvalue/agent.sock")
	}
}

func TestIdentityAgentFromConfig(t *testing.T) {
	tmpDir := writeSSHConfig(t, `Host github.com
    HostName github.com
    User git

Host *
    IdentityAgent ~/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock
    AddKeysToAgent yes
`)
	sock := GetSSHAuthSock()
	expected := filepath.Join(tmpDir, "Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock")
	if sock != expected {
		t.Errorf("GetSSHAuthSock() = %q, want %q", sock, expected)
	}
}

func TestIdentityAgentSSHAuthSockTakesPrecedence(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/existing.sock")
	sock := GetSSHAuthSock()
	if sock != "/tmp/existing.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want SSH_AUTH_SOCK value", sock)
	}
}

func TestIdentityAgentNoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("SSH_AUTH_SOCK", "")

	sock := GetSSHAuthSock()
	if sock != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty", sock)
	}
}

func TestIdentityAgentQuotedPath(t *testing.T) {
	tmpDir := writeSSHConfig(t, `Host *
    IdentityAgent "~/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"
`)
	sock := GetSSHAuthSock()
	expected := filepath.Join(tmpDir, "Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock")
	if sock != expected {
		t.Errorf("GetSSHAuthSock() with quoted path = %q, want %q", sock, expected)
	}
}

func TestIdentityAgentMultiPatternHost(t *testing.T) {
	writeSSHConfig(t, `Host * !excluded.example.com
    IdentityAgent /tmp/wildcard-agent.sock
`)
	sock := GetSSHAuthSock()
	if sock != "/tmp/wildcard-agent.sock" {
		t.Errorf("GetSSHAuthSock() with multi-pattern Host = %q, want /tmp/wildcard-agent.sock", sock)
	}
}

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"~/path/to/agent.sock"`, "~/path/to/agent.sock"},
		{`'~/path/to/agent.sock'`, "~/path/to/agent.sock"},
		{"~/path/to/agent.sock", "~/path/to/agent.sock"},
		{`"mismatched'`, `"mismatched'`},
		{"", ""},
	}
	for _, tt := range tests {
		got := stripQuotes(tt.input)
		if got != tt.expected {
			t.Errorf("stripQuotes(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsUniversalWildcard(t *testing.T) {
	tests := []struct {
		pattern  string
		expected bool
	}{
		{"*", true},
		{"* !excluded.example.com", true},
		// *.example.com is a domain-scoped pattern, NOT a universal wildcard
		{"*.example.com", false},
		{"github.com", false},
		{"!excluded.example.com", false},
		{"github.com bitbucket.org", false},
	}
	for _, tt := range tests {
		got := isUniversalWildcard(tt.pattern)
		if got != tt.expected {
			t.Errorf("isUniversalWildcard(%q) = %v, want %v", tt.pattern, got, tt.expected)
		}
	}
}

func TestIdentityAgentHostSpecificIsIgnored(t *testing.T) {
	writeSSHConfig(t, `Host github.com
    IdentityAgent /tmp/github-only.sock
`)
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty for host-specific-only config", got)
	}
}

func TestIdentityAgentMatchScopeIsIgnored(t *testing.T) {
	writeSSHConfig(t, `Match host github.com
    IdentityAgent /tmp/match-only.sock
`)
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty for match-scoped-only config", got)
	}
}

func TestIdentityAgentGlobalScope(t *testing.T) {
	writeSSHConfig(t, `IdentityAgent /tmp/global-agent.sock

Host github.com
    HostName github.com
`)
	sock := GetSSHAuthSock()
	if sock != "/tmp/global-agent.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/global-agent.sock", sock)
	}
}
