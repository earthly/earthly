package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidTransport = fmt.Errorf("invalid transport")
	ErrInvalidAuth      = fmt.Errorf("invalid auth")
)

type GlobalConfig struct {
	CachePath           string `yaml:"cache_path"`
	DisableLoopDevice   bool   `yaml:"no_loop_device"`
	BuildkitCacheSizeMb int    `yaml:"cache_size_mb"`
	BuildkitImage       string `yaml:"buildkit_image"`
}

type GitConfig struct {
	// these are used for global config
	GitURLInsteadOf string `yaml:"url_instead_of"`

	// these are used for git vendors (e.g. github, gitlab)
	Auth     string `yaml:"auth"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

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

func ParseConfigFile(yamlData []byte) (*Config, error) {
	config := Config{
		Global: GlobalConfig{
			CachePath:           "/var/cache/earthly",
			DisableLoopDevice:   false,
			BuildkitCacheSizeMb: 10000,
		},
	}

	err := yaml.Unmarshal(yamlData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func CreateGitConfig(config *Config) (string, []string, error) {
	credentials := []string{}
	lines := []string{}
	cred_i := 0

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
	//		lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", cred_i))
	//		credentials = append(credentials, fmt.Sprintf("echo password=%q", v.Password))
	//		cred_i++

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
			lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", cred_i))
			credentials = append(credentials, v.Password)
			cred_i++

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
