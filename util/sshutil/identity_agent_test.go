package sshutil

import (
	"os"
	"path/filepath"
	"testing"
)

// writeSSHConfig creates a temp HOME, writes configContent to ~/.ssh/config,
// sets HOME, and clears SSH_AUTH_SOCK for the test. Returns the tmpDir so
// callers can build expected expanded paths.
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
		input, expected string
	}{
		{"~/some/path", filepath.Join(home, "some/path")},
		{"~", home},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}
	for _, tt := range tests {
		if got := expandPath(tt.input, home); got != tt.expected {
			t.Errorf("expandPath(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestExpandPathEnvVar(t *testing.T) {
	t.Setenv("TEST_SSH_VAR", "myvalue")
	got := expandPath("$TEST_SSH_VAR/agent.sock", "/home/u")
	if got != "myvalue/agent.sock" {
		t.Errorf("expandPath env var = %q", got)
	}
}

func TestSSHAuthSockTakesPrecedence(t *testing.T) {
	t.Setenv("SSH_AUTH_SOCK", "/tmp/existing.sock")
	if got := GetSSHAuthSock(); got != "/tmp/existing.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want SSH_AUTH_SOCK value", got)
	}
}

func TestNoConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	t.Setenv("SSH_AUTH_SOCK", "")
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty", got)
	}
}

func TestIdentityAgentWildcardHost(t *testing.T) {
	tmpDir := writeSSHConfig(t, `Host github.com
    HostName github.com
    User git

Host *
    IdentityAgent ~/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock
    AddKeysToAgent yes
`)
	want := filepath.Join(tmpDir, "Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock")
	if got := GetSSHAuthSock(); got != want {
		t.Errorf("GetSSHAuthSock() = %q, want %q", got, want)
	}
}

func TestIdentityAgentQuotedPath(t *testing.T) {
	tmpDir := writeSSHConfig(t, `Host *
    IdentityAgent "~/Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock"
`)
	want := filepath.Join(tmpDir, "Library/Group Containers/2BUA8C4S2C.com.1password/t/agent.sock")
	if got := GetSSHAuthSock(); got != want {
		t.Errorf("GetSSHAuthSock() = %q, want %q", got, want)
	}
}

func TestIdentityAgentMultiPatternHost(t *testing.T) {
	writeSSHConfig(t, `Host * !excluded.example.com
    IdentityAgent /tmp/wildcard.sock
`)
	if got := GetSSHAuthSock(); got != "/tmp/wildcard.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/wildcard.sock", got)
	}
}

func TestIdentityAgentGlobalScope(t *testing.T) {
	writeSSHConfig(t, `IdentityAgent /tmp/global.sock

Host github.com
    User git
`)
	if got := GetSSHAuthSock(); got != "/tmp/global.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/global.sock", got)
	}
}

func TestIdentityAgentHostSpecificIsIgnored(t *testing.T) {
	writeSSHConfig(t, `Host github.com
    IdentityAgent /tmp/github-only.sock
`)
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty for host-specific config", got)
	}
}

func TestIdentityAgentMatchScopeIsIgnored(t *testing.T) {
	writeSSHConfig(t, `Match host github.com
    IdentityAgent /tmp/match-only.sock
`)
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty for match-only config", got)
	}
}

func TestIdentityAgentScopedDomainNotUniversal(t *testing.T) {
	// "*.example.com" is host-scoped, not universal.
	writeSSHConfig(t, `Host *.example.com
    IdentityAgent /tmp/example.sock
`)
	if got := GetSSHAuthSock(); got != "" {
		t.Errorf("GetSSHAuthSock() = %q, want empty for *.example.com scope", got)
	}
}

func TestIdentityAgentTabSeparated(t *testing.T) {
	writeSSHConfig(t, "Host\t*\n\tIdentityAgent\t/tmp/tab.sock\n")
	if got := GetSSHAuthSock(); got != "/tmp/tab.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/tab.sock", got)
	}
}

func TestIdentityAgentEqualsSyntax(t *testing.T) {
	writeSSHConfig(t, "Host=*\nIdentityAgent=/tmp/eq.sock\n")
	if got := GetSSHAuthSock(); got != "/tmp/eq.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/eq.sock", got)
	}
}

func TestIdentityAgentCommentsAndBlankLines(t *testing.T) {
	writeSSHConfig(t, `# Comment line

Host *
    # IdentityAgent /tmp/commented-out.sock
    IdentityAgent /tmp/real.sock
`)
	if got := GetSSHAuthSock(); got != "/tmp/real.sock" {
		t.Errorf("GetSSHAuthSock() = %q, want /tmp/real.sock", got)
	}
}

func TestStripQuotes(t *testing.T) {
	tests := []struct{ in, want string }{
		{`"hello"`, "hello"},
		{`'hello'`, "hello"},
		{`hello`, "hello"},
		{`"hello`, `"hello`},
		{`"`, `"`},
		{``, ``},
	}
	for _, tt := range tests {
		if got := stripQuotes(tt.in); got != tt.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestIsUniversalWildcard(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"*", true},
		{"* !excluded.com", true},
		{"!a *", true},
		{"*.example.com", false},
		{"github.com", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isUniversalWildcard(tt.in); got != tt.want {
			t.Errorf("isUniversalWildcard(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestSplitDirective(t *testing.T) {
	tests := []struct {
		in, kw, val string
	}{
		{"Host *", "Host", "*"},
		{"Host\t*", "Host", "*"},
		{"IdentityAgent=/tmp/a.sock", "IdentityAgent", "/tmp/a.sock"},
		{"IdentityAgent  =  /tmp/a.sock", "IdentityAgent", "/tmp/a.sock"},
		{"loner", "loner", ""},
	}
	for _, tt := range tests {
		kw, val := splitDirective(tt.in)
		if kw != tt.kw || val != tt.val {
			t.Errorf("splitDirective(%q) = (%q,%q), want (%q,%q)", tt.in, kw, val, tt.kw, tt.val)
		}
	}
}
