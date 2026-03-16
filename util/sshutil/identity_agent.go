// Package sshutil provides helpers for SSH configuration and authentication.
// Currently it exposes GetSSHAuthSock to resolve the SSH agent socket path
// from the environment or from ~/.ssh/config (IdentityAgent directive).
package sshutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GetSSHAuthSock returns the SSH auth socket path to use. It checks the
// SSH_AUTH_SOCK environment variable first. If that is not set, it falls back
// to reading the IdentityAgent directive from ~/.ssh/config (Host * or global
// scope). Returns an empty string when no socket can be determined.
func GetSSHAuthSock() string {
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		return sock
	}
	return identityAgentFromConfig()
}

// identityAgentFromConfig reads ~/.ssh/config and returns the IdentityAgent
// value from the first matching universal wildcard Host block (i.e. "Host *")
// or the first IdentityAgent directive that appears before any Host block.
// IdentityAgent directives inside specific host blocks (e.g. "Host github.com")
// are intentionally ignored so that host-scoped keys don't pollute the default
// socket path. The returned path has ~ and environment variables expanded.
// Returns empty string when nothing is found or the config file cannot be read.
func identityAgentFromConfig() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configPath := filepath.Join(home, ".ssh", "config")
	f, err := os.Open(configPath)
	if err != nil {
		return ""
	}
	defer f.Close()

	var inWildcardHost bool
	var seenAnyHostBlock bool
	var globalAgent string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "host ") {
			seenAnyHostBlock = true
			inWildcardHost = isUniversalWildcard(strings.TrimSpace(line[5:]))
			continue
		}
		if strings.HasPrefix(lower, "match ") {
			seenAnyHostBlock = true
			inWildcardHost = false
			continue
		}

		if strings.HasPrefix(lower, "identityagent ") {
			value := strings.TrimSpace(line[len("identityagent "):])
			// Handle optional = syntax: IdentityAgent = /path
			value = strings.TrimPrefix(value, "=")
			value = strings.TrimSpace(value)
			// Strip surrounding quotes that SSH config allows for paths with spaces
			value = stripQuotes(value)
			if value == "" {
				continue
			}
			expanded := expandPath(value, home)
			if inWildcardHost {
				return expanded
			}
			// Only treat as global fallback if we haven't entered any Host block yet.
			// An IdentityAgent inside a specific host block (e.g. Host github.com)
			// must not bleed into the default socket path.
			if !seenAnyHostBlock && globalAgent == "" {
				globalAgent = expanded
			}
		}
	}

	return globalAgent
}

// isUniversalWildcard reports whether a Host pattern line represents a
// universal match (i.e. matches all hosts). It requires at least one token
// that is exactly "*". Patterns like "*.example.com" are host-scoped and are
// not treated as universal. Negation tokens (starting with "!") are skipped.
//
// Examples:
//
//	"*"                     -> true   (matches all hosts)
//	"* !excluded.example.com" -> true (universal with exclusion)
//	"*.example.com"         -> false  (scoped to a domain)
//	"github.com"            -> false  (specific host)
func isUniversalWildcard(pattern string) bool {
	for _, token := range strings.Fields(pattern) {
		if strings.HasPrefix(token, "!") {
			continue
		}
		if token == "*" {
			return true
		}
	}
	return false
}

// stripQuotes removes a single pair of matching surrounding quotes (single or
// double) from s if present. It does not remove mismatched or nested quotes.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// expandPath expands a leading ~ to the user home directory and replaces
// $VAR and ${VAR} references using os.ExpandEnv.
func expandPath(path string, home string) string {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		path = home
	}
	return os.ExpandEnv(path)
}
