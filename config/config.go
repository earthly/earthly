package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidTransport occurs when a URL transport type is invalid
	ErrInvalidTransport = fmt.Errorf("invalid transport")

	// ErrInvalidAuth occurs when the auth type is invalid
	ErrInvalidAuth = fmt.Errorf("invalid auth")
)

// GlobalConfig contains global config values
type GlobalConfig struct {
	DisableAnalytics         bool     `yaml:"disable_analytics" help:"Controls Earthly telemetry."`
	BuildkitCacheSizeMb      int      `yaml:"cache_size_mb" help:"Size of the buildkit cache in Megabytes."`
	BuildkitImage            string   `yaml:"buildkit_image" help:"Choose a specific image for your buildkitd."`
	DebuggerPort             int      `yaml:"debugger_port" help:"What port should the debugger (and other interactive sessions) use to communicate."`
	BuildkitRestartTimeoutS  int      `yaml:"buildkit_restart_timeout_s" help:"How long to wait for buildkit to (re)start, in seconds."`
	BuildkitAdditionalArgs   []string `yaml:"buildkit_additional_args" help:"Additional args to pass to buildkit when it starts. Useful for custom certs, or user namespace complications."`
	BuildkitAdditionalConfig string   `yaml:"buildkit_additional_config" help:"Additional config to use when starting the buildkit container; like mounting volumes."`

	// Obsolete.
	CachePath string `yaml:"cache_path"`
}

// GitConfig contains git-specific config values
type GitConfig struct {
	// these are used for global config
	GitURLInsteadOf string `yaml:"url_instead_of"`

	// these are used for git vendors (e.g. github, gitlab)
	Pattern    string `yaml:"pattern"  help:"A pattern for how to match parts of your git URL"`
	Substitute string `yaml:"substitute" help:"Substitute?"`
	Suffix     string `yaml:"suffix" help:"Suffix?"`                                                 // .git
	Auth       string `yaml:"auth" help:"What authentication method do you use? (http, https, ssh)"` // http, https, ssh
	User       string `yaml:"user" help:"The user you use to connect to your git server."`
	Password   string `yaml:"password" help:"The password for your user."`
	KeyScan    string `yaml:"serverkey" help:"KeyScan?"`
}

// Config contains user's configuration values from ~/earthly/config.yml
type Config struct {
	Global GlobalConfig         `yaml:"global" help:"Global configuration object. Requires YAML literal to set directly."`
	Git    map[string]GitConfig `yaml:"git" help:"Git configuration object. Requires YAML literal to set directly."`
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
			BuildkitCacheSizeMb:     0,
			DebuggerPort:            8373,
			BuildkitRestartTimeoutS: 60,
			BuildkitAdditionalArgs:  []string{},
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
				Auth:   "ssh",
				Suffix: ".git",
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

func keyAndValueCompatible(key reflect.Type, value *yaml.Node) bool {
	var val interface{}
	switch key.Kind() {
	// add other types as needed as they are introduced in the config struct
	case reflect.Map:
		val = reflect.MakeMap(key).Interface()
	default:
		val = reflect.New(key).Interface()
	}

	err := value.Decode(val)

	return err == nil
}

// UpsertConfig adds or modifies the key to be the specified value.
// This is saved to disk in your earthly config file.
func UpsertConfig(config []byte, path, value string) ([]byte, error) {
	base := &yaml.Node{}
	yaml.Unmarshal(config, base)

	if base.IsZero() {
		// Empty file, or a simple comment results in a null document.
		// Not handled well, so manufacture somewhat acceptable document
		fullDoc := string(config) + "\n---"
		yaml.Unmarshal([]byte(fullDoc), base)
		base.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}

	pathParts := splitPath(path)

	t, help, err := validatePath(reflect.TypeOf(Config{}), pathParts)
	if err != nil {
		return []byte{}, errors.Wrap(err, "path is not valid")
	}

	if value == "--help" {
		fmt.Printf("(%s): %s\n", t.Kind(), help)
		return []byte{}, nil
	}

	yamlValue, err := valueToYaml(value)
	if err != nil {
		return []byte{}, errors.Wrap(err, "could not parse value")
	}

	if !keyAndValueCompatible(t, yamlValue) {
		return []byte{}, fmt.Errorf("Cannot set %s to %v, as the types are incompatible", path, value)
	}

	setYamlValue(base, pathParts, yamlValue)

	newConfig, err := yaml.Marshal(base)
	if err != nil {
		return []byte{}, err
	}

	return newConfig, nil
}

func splitPath(path string) []string {
	// Allow quotes to group keys, since git repos are keys and have periods... this is why we dont just strings.Split
	// If you screw up the quotes you will get a weird invalid path later.
	re := regexp.MustCompile(`[^\."']+|"([^"]*)"|'([^']*)`)
	pathParts := re.FindAllString(path, -1)

	for i := 0; i < len(pathParts); i++ {
		// If we did have a quoted string we need to prune it
		pathParts[i] = strings.Trim(pathParts[i], `"`)
	}

	return pathParts
}

func validatePath(t reflect.Type, path []string) (reflect.Type, string, error) {
	if len(path) == 0 {
		return nil, "", errors.New("No path present")
	}

	if t.Kind() == reflect.Map {
		// Maps are only for git repos. Grab the kind on the other side of the map
		// and advance; to validate the path on the other side of the repo name

		// path is a git."some.repo", so we can't advance
		if len(path) == 1 {
			// base case. I am not happy with this. Will need to change if we get more than one map in the config.
			return t.Elem(), "Git repository. Quote names with dots in them, like this: git.\"github.com\". Requires YAML literal to set directly.", nil
		}

		return validatePath(t.Elem(), path[1:])
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")
		helpTag := field.Tag.Get("help")

		if yamlTag == path[0] {
			if len(path) == 1 {
				// base case
				return field.Type, helpTag, nil
			}

			return validatePath(field.Type, path[1:])
		}
	}

	return nil, "", fmt.Errorf("no path for %s", strings.Join(path, "."))
}

func valueToYaml(value string) (*yaml.Node, error) {
	valueNode := &yaml.Node{}
	if err := yaml.Unmarshal([]byte(value), valueNode); err != nil {
		return nil, fmt.Errorf("%s is not a valid YAML value", value)
	}

	// Unfold all the yaml so its not mixed inline and flow styles in the final document
	var fixStyling func(node *yaml.Node)
	fixStyling = func(node *yaml.Node) {
		node.Style = 0

		for _, n := range node.Content {
			fixStyling(n)
		}
	}
	fixStyling(valueNode)

	return valueNode.Content[0], nil
}

func pathToYaml(path []string, value *yaml.Node) []*yaml.Node {
	yamlNodes := []*yaml.Node{}

	var last *yaml.Node

	for i, seg := range path {
		key := &yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: seg,
		}

		mapping := &yaml.Node{
			Kind: yaml.MappingNode,
		}

		if i == len(path)-1 {
			// Last node should assign path as the value, not another mapping node
			// Otherwise we would need to dig it up again.

			if last == nil {
				// Single depth special case
				yamlNodes = append(yamlNodes, key, value)
				continue
			}

			last.Content = append(last.Content, key, value)
		} else if last == nil {
			// First, top level mapping node
			yamlNodes = append(yamlNodes, key, mapping)
			last = mapping
		} else {
			// Middle of the road regular case
			last.Content = append(last.Content, key, mapping)
			last = mapping
		}
	}

	return yamlNodes
}

func setYamlValue(node *yaml.Node, path []string, value *yaml.Node) []string {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, c := range node.Content {
			path = setYamlValue(c, path, value)
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			// Keys/Values are inline. Count by twos to get it right.
			key := node.Content[i]
			val := node.Content[i+1]

			if len(path) > 0 && key.Value == path[0] {
				path = path[1:]

				if len(path) == 0 {
					node.Content[i+1] = value
					return []string{}
				}

				path = setYamlValue(val, path, value)
			}
		}

	default: // Sequence, Scalar nodes get skipped
		return path
	}

	// If we get here, we have consumed all the path possible.
	// Build YAML and add it from where we are at.
	if len(path) > 0 {
		yamlMap := pathToYaml(path, value)
		node.Content = append(node.Content, yamlMap...)
	}

	return []string{}
}

// ReadConfigFile reads in the config file from the disk, into a byte slice.
func ReadConfigFile(configPath string, contextSet bool) ([]byte, error) {
	yamlData, err := ioutil.ReadFile(configPath)
	if os.IsNotExist(err) && !contextSet {
		return []byte{}, nil
	} else if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to read from %s", configPath)
	}

	return yamlData, nil
}

// WriteConfigFile writes the config file to disk with preset permission 0644
func WriteConfigFile(configPath string, data []byte) error {
	err := os.MkdirAll(path.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, 0644)
}
