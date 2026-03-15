package sshutil

import (
	"os"
	"path/filepath"
	"testing"
)

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
	// Create a temp dir with a fake ssh config
	tmpDir := t.TempDir()
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}

	configContent := `Host github.com
    HostName github.com
    User git

Host *
    IdentityAgent ~/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock
    AddKeysToAgent yes
`
	if err := os.WriteFile(filepath.Join(sshDir, "config"), []byte(configContent), 0600); err != nil {
		t.Fatal(err)
	}

	// Override HOME so identityAgentFromConfig finds our temp config
	t.Setenv("HOME", tmpDir)
	// Clear SSH_AUTH_SOCK so GetSSHAuthSock falls through
	t.Setenv("SSH_AUTH_SOCK", "")

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

func TestIdentityAgentGlobalScope(t *testing.T) {
	tmpDir := t.TempDir()
	sshDir := filepath.Join(tmpDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		t.Fatal(err)
	}

	// IdentityAgent before any Host block
	configContent := `IdentityAgent /tmp/global-agent.sock

Host github.com
    HostName github.com
`
	if err := os.WriteFile(filepath.Join(sshDir, "config"), []byte(configContent), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("HOME", tmpDir)
	t.Setenv("SSH_AUTH_SOCK", "")

	sock := GetSSHAuthSock()
	if sock != "/tmp/global-agent.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/global-agent.sock", sock)
	}
}
