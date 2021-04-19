package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	// ErrInvalidTransport occurs when a URL transport type is invalid
	ErrInvalidTransport = errors.Errorf("invalid transport")

	// ErrInvalidAuth occurs when the auth type is invalid
	ErrInvalidAuth = errors.Errorf("invalid auth")
)

// GlobalConfig contains global config values
type GlobalConfig struct {
	DisableAnalytics         bool     `yaml:"disable_analytics"          help:"Controls Earthly telemetry."`
	BuildkitCacheSizeMb      int      `yaml:"cache_size_mb"              help:"Size of the buildkit cache in Megabytes."`
	BuildkitImage            string   `yaml:"buildkit_image"             help:"Choose a specific image for your buildkitd."`
	DebuggerPort             int      `yaml:"debugger_port"              help:"What port should the debugger (and other interactive sessions) use to communicate."`
	BuildkitRestartTimeoutS  int      `yaml:"buildkit_restart_timeout_s" help:"How long to wait for buildkit to (re)start, in seconds."`
	BuildkitAdditionalArgs   []string `yaml:"buildkit_additional_args"   help:"Additional args to pass to buildkit when it starts. Useful for custom/self-signed certs, or user namespace complications."`
	BuildkitAdditionalConfig string   `yaml:"buildkit_additional_config" help:"Additional config to use when starting the buildkit container; like using custom/self-signed certificates."`
	CniMtu                   uint16   `yaml:"cni_mtu"                    help:"Override auto-detection of the default interface MTU, for all containers within buildkit"`

	// Obsolete.
	CachePath string `yaml:"cache_path"`
}

// GitConfig contains git-specific config values
type GitConfig struct {
	// these are used for global config
	GitURLInsteadOf string `yaml:"url_instead_of"`

	// these are used for git vendors (e.g. github, gitlab)
	Pattern    string `yaml:"pattern"    help:"A regular expression defined to match git URLs, defaults to the regex: <site>/([^/]+)/([^/]+). For example if the site is github.com, then the default pattern will match github.com/<user>/<repo>."`
	Substitute string `yaml:"substitute" help:"If specified, a regular expression substitution will be preformed to determine which URL is cloned by git. Values like $1, $2, ... will be replaced with matched subgroup data. If no substitute is given, a URL will be created based on the requested SSH authentication mode."`
	Suffix     string `yaml:"suffix"     help:"The git repository suffix, like .git."`                                       // .git
	Auth       string `yaml:"auth"       help:"What authentication method do you use? Valid options are: http, https, ssh."` // http, https, ssh
	User       string `yaml:"user"       help:"The https username to use when auth is set to https. This setting is ignored when auth is ssh."`
	Password   string `yaml:"password"   help:"The https password to use when auth is set to https. This setting is ignored when auth is ssh."`
	KeyScan    string `yaml:"serverkey"  help:"SSH fingerprints, like you would add in your known hosts file, or get from ssh-keyscan."`
}

// Config contains user's configuration values from ~/earthly/config.yml
type Config struct {
	Global GlobalConfig         `yaml:"global" help:"Global configuration object. Requires YAML literal to set directly."`
	Git    map[string]GitConfig `yaml:"git" help:"Git configuration object. Requires YAML literal to set directly."`
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
		return []byte{}, errors.Errorf("cannot set %s to %v, as the types are incompatible", path, value)
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

	return nil, "", errors.Errorf("no path for %s", strings.Join(path, "."))
}

func valueToYaml(value string) (*yaml.Node, error) {
	valueNode := &yaml.Node{}
	if err := yaml.Unmarshal([]byte(value), valueNode); err != nil {
		return nil, errors.Errorf("%s is not a valid YAML value", value)
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
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configPath, data, 0644)
}
