package sshutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// GetSSHAuthSock returns the SSH auth socket path by checking:
// 1. The SSH_AUTH_SOCK environment variable
// 2. The IdentityAgent directive from ~/.ssh/config (Host *)
func GetSSHAuthSock() string {
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		return sock
	}
	return identityAgentFromConfig()
}

// identityAgentFromConfig reads ~/.ssh/config and returns the IdentityAgent
// value from the Host * block (or the first IdentityAgent found outside any
// host block), with ~ and environment variables expanded.
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
	var globalAgent string

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "host ") {
			pattern := strings.TrimSpace(line[5:])
			inWildcardHost = pattern == "*"
			continue
		}
		if strings.HasPrefix(lower, "match ") {
			inWildcardHost = false
			continue
		}

		if strings.HasPrefix(lower, "identityagent ") {
			value := strings.TrimSpace(line[len("identityagent "):])
			// Also handle = syntax: IdentityAgent = /path
			value = strings.TrimPrefix(value, "=")
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}
			expanded := expandPath(value, home)
			if inWildcardHost {
				return expanded
			}
			// Store first seen global (before any Host block) as fallback
			if globalAgent == "" {
				globalAgent = expanded
			}
		}
	}

	return globalAgent
}

// expandPath expands ~ to the home directory and $VAR / ${VAR} environment variables.
func expandPath(path string, home string) string {
	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		path = home
	}
	return os.ExpandEnv(path)
}
