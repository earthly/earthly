package config

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

var ErrInvalidTransport = fmt.Errorf("invalid transport")

type GlobalConfig struct {
	//TODO add support for global config as needed
}

type GitConfig struct {
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
	var config Config

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

	// iterate over map in a consistent order otherwise it will cause the buildkitd image to restart
	// due to the settings hash being different
	keys := []string{}
	for k := range config.Git {
		if k != "default" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// default
	if v, ok := config.Git["default"]; ok {
		if v.Auth == "https" {
			lines = append(lines, fmt.Sprintf("[credential]"))
			lines = append(lines, fmt.Sprintf("  username=%q", v.User))
			lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", cred_i))
			credentials = append(credentials, fmt.Sprintf("echo password=%q", v.Password))
			cred_i++

			// use https instead of ssh://git@....
			lines = append(lines, fmt.Sprintf("[url \"https://\"]"))
			lines = append(lines, fmt.Sprintf("  insteadOf = ssh://git@"))
		}
	}

	for _, k := range keys {
		v := config.Git[k]
		if v.Auth == "https" {
			url, err := ensureTransport(k, "https")
			if err != nil {
				return "", nil, err
			}
			lines = append(lines, fmt.Sprintf("[credential %q]", url))
			lines = append(lines, fmt.Sprintf("  username=%q", v.User))
			lines = append(lines, fmt.Sprintf("  helper=/usr/bin/git_credentials_%d", cred_i))
			credentials = append(credentials, fmt.Sprintf("echo password=%q", v.Password))
			cred_i++

			// use https instead of ssh://git@....
			lines = append(lines, fmt.Sprintf("[url %q]", url))
			lines = append(lines, fmt.Sprintf("  insteadOf = ssh://git@%s", url[8:]))
		}
	}

	lines = append(lines, "")
	gitConfig := strings.Join(lines, "\n")

	return gitConfig, credentials, nil
}
