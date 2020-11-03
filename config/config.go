package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	// ErrInvalidTransport occurs when a URL transport type is invalid
	ErrInvalidTransport = fmt.Errorf("invalid transport")

	// ErrInvalidAuth occurs when the auth type is invalid
	ErrInvalidAuth = fmt.Errorf("invalid auth")
)

// GlobalConfig contains global config values
type GlobalConfig struct {
	RunPath                 string `yaml:"run_path"`
	DisableAnalytics        bool   `yaml:"disable_analytics"`
	BuildkitCacheSizeMb     int    `yaml:"cache_size_mb"`
	BuildkitImage           string `yaml:"buildkit_image"`
	DebuggerPort            int    `yaml:"debugger_port"`
	BuildkitRestartTimeoutS int    `yaml:"buildkit_restart_timeout_s"`

	// Obsolete.
	CachePath string `yaml:"cache_path"`
}

// GitConfig contains git-specific config values
type GitConfig struct {
	// these are used for global config
	GitURLInsteadOf string `yaml:"url_instead_of"`

	// these are used for git vendors (e.g. github, gitlab)
	Auth     string `yaml:"auth"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Config contains user's configuration values from ~/earthly/config.yml
type Config struct {
	Global GlobalConfig         `yaml:"global"`
	Git    map[string]GitConfig `yaml:"git"`
}

func ensureTransport(s, transport string) (string, error) {
	parts := strings.SplitN(s, "://", 2)
	if len(parts) == 2 {
		if parts[0] != transport {
			return "", ErrInvalidTransport
		}
	}
	return transport + "://" + s, nil
}

// ParseConfigFile parse config data
func ParseConfigFile(yamlData []byte) (*Config, error) {
	// pre-populate defaults
	config := Config{
		Global: GlobalConfig{
			RunPath:                 defaultRunPath(),
			BuildkitCacheSizeMb:     0,
			DebuggerPort:            8373,
			BuildkitRestartTimeoutS: 60,
		},
	}

	err := yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// CreateGitConfig returns the contents of the /root/.gitconfig file and a list of corresponding
// password credentials (the passwords are stored as env variables rather than written to disk)
func CreateGitConfig(config *Config) (string, []string, error) {
	credentials := []string{}
	lines := []string{}
	credIndex := 0

	// automatically add default auth=ssh for known sites
	defaultSites := []string{"github.com", "gitlab.com"}
	for _, k := range defaultSites {
		if _, ok := config.Git[k]; !ok {
			config.Git[k] = GitConfig{
				Auth: "ssh",
			}
		}
	}

	// iterate over map in a consistent order otherwise it will cause the buildkitd image to restart
	// due to the settings hash being different
	keys := []string{}
	for k := range config.Git {
		if k != "default" && k != "global" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// TODO figure out how to get the URL rewritting working for the generic case for all URLs
	// default
	//if v, ok := config.Git["default"]; ok {
	//	if v.Auth == "https" {
	//		lines = append(lines, fmt.Sprintf("[credential]"))
	//		lines = append(lines, fmt.Sprintf("  username=%q", v.User))
	//		lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", credIndex))
	//		credentials = append(credentials, fmt.Sprintf("echo password=%q", v.Password))
	//		credIndex++

	//		// use https instead of ssh://git@....
	//		lines = append(lines, fmt.Sprintf("[url \"https://\"]"))
	//		lines = append(lines, fmt.Sprintf("  insteadOf = ssh://git@"))
	//	}
	//}

	for _, k := range keys {
		v := config.Git[k]

		url, err := ensureTransport(k, "https")
		if err != nil {
			return "", nil, err
		}

		switch v.Auth {
		case "https":
			lines = append(lines, fmt.Sprintf("[credential %q]", url))
			lines = append(lines, fmt.Sprintf("  username=%q", v.User))
			lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", credIndex))
			credentials = append(credentials, v.Password)
			credIndex++

			// use https instead of ssh://git@....
			lines = append(lines, fmt.Sprintf("[url %q]", url+"/"))
			lines = append(lines, fmt.Sprintf("  insteadOf = git@%s:", url[8:]))
		case "ssh":
			// use git@... instead of https://...
			lines = append(lines, fmt.Sprintf("[url %q]", "git@"+url[8:]+":"))
			lines = append(lines, fmt.Sprintf("  insteadOf = %s:", url+"/"))
		default:
			return "", nil, errors.Wrapf(ErrInvalidAuth, "unsupported auth %s for site %s", v.Auth, k)
		}
	}

	lines = append(lines, "")
	gitConfig := strings.Join(lines, "\n")

	return gitConfig, credentials, nil
}

func defaultRunPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homeDir, ".earthly/run")
}
