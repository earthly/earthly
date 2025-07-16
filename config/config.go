package config

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/earthly/earthly/util/cliutil"
)

const (
	// DefaultLocalRegistryPort is the default user-facing port for the local registry used for exports.
	DefaultLocalRegistryPort = 8371

	// DefaultDarwinProxyImage is the image tag used for the Docker Desktop registry proxy on Darwin.
	DefaultDarwinProxyImage = "alpine/socat:1.7.4.4"

	// DefaultDarwinProxyWait is the maximum time to wait for the Darwin registry proxy support container to become available.
	DefaultDarwinProxyWait = 10 * time.Second

	// DefaultBuildkitScheme is the default scheme earthly uses to connect to its buildkitd. tcp or docker-container.
	DefaultBuildkitScheme = "docker-container"

	// DefaultConversionParallelism is the default conversion parallelism that Earthly uses internally to generate LLB for BuildKit to consume.
	DefaultConversionParallelism = 10

	// DefaultBuildkitMaxParallelism is the default max parallelism for buildkit workers.
	DefaultBuildkitMaxParallelism = 20

	// DefaultCACert is the default path to use when looking for a CA cert to use for TLS.
	DefaultCACert = "./certs/ca_cert.pem"

	// DefaultCAKey is the default path to use when looking for a CA key to use for TLS cert generation.
	DefaultCAKey = "./certs/ca_key.pem"

	// DefaultClientTLSCert is the default path to use when looking for the Earthly TLS cert
	DefaultClientTLSCert = "./certs/earthly_cert.pem"

	// DefaultClientTLSKey is the default path to use when looking for the Earthly TLS key
	DefaultClientTLSKey = "./certs/earthly_key.pem"

	// DefaultServerTLSCert is the default path to use when looking for the Buildkit TLS cert
	DefaultServerTLSCert = "./certs/buildkit_cert.pem"

	// DefaultServerTLSKey is the default path to use when looking for the Buildkit TLS key
	DefaultServerTLSKey = "./certs/buildkit_key.pem"

	// DefaultContainerFrontend is the default frontend program or interfacing with the running containers and saved images
	DefaultContainerFrontend = "auto"
)

var (
	// ErrInvalidTransport occurs when a URL transport type is invalid
	ErrInvalidTransport = errors.Errorf("invalid transport")

	// ErrInvalidAuth occurs when the auth type is invalid
	ErrInvalidAuth = errors.Errorf("invalid auth")
)

// GlobalConfig contains global config values
type GlobalConfig struct {
	BuildkitCacheSizeMb        int           `yaml:"cache_size_mb"                  help:"Size of the buildkit cache in Megabytes."`
	BuildkitCacheSizePct       int           `yaml:"cache_size_pct"                 help:"Size of the buildkit cache, as percentage (0-100)."`
	BuildkitCacheKeepDurationS int           `yaml:"buildkit_cache_keep_duration_s" help:"Max age of cache, in seconds. 0 disables age-based cache expiry."`
	BuildkitImage              string        `yaml:"buildkit_image"                 help:"Choose a specific image for your buildkitd."`
	BuildkitRestartTimeoutS    int           `yaml:"buildkit_restart_timeout_s"     help:"How long to wait for buildkit to (re)start, in seconds."`
	BuildkitAdditionalArgs     []string      `yaml:"buildkit_additional_args"       help:"Additional args to pass to buildkit when it starts. Useful for custom/self-signed certs, or user namespace complications."`
	BuildkitAdditionalConfig   string        `yaml:"buildkit_additional_config"     help:"Additional config to use when starting the buildkit container; like using custom/self-signed certificates."`
	BuildkitMaxParallelism     int           `yaml:"buildkit_max_parallelism"       help:"Max parallelism for buildkit workers"`
	ConversionParallelism      int           `yaml:"conversion_parallelism"         help:"Set the conversion parallelism for speeding up the use of IF, WITH, DOCKER --load, FROMDOCKERFILE and others. A value of 0 disables the feature"`
	CniMtu                     uint16        `yaml:"cni_mtu"                        help:"Override auto-detection of the default interface MTU, for all containers within buildkit"`
	BuildkitHost               string        `yaml:"buildkit_host"                  help:"The URL of your buildkit, remote or local."`
	LocalRegistryHost          string        `yaml:"local_registry_host"            help:"The URL of the local registry used for image exports to Docker."`
	DarwinProxyImage           string        `yaml:"darwin_proxy_image"             help:"The container image & tag used for the Docker Desktop registry proxy."`
	DarwinProxyWait            time.Duration `yaml:"darwin_proxy_wait"              help:"The maximum time to wait for the Darwin registry proxy support container to become available."`
	TLSCACert                  string        `yaml:"tlsca"                          help:"The path to the CA cert for verification. Relative paths are interpreted as relative to the config path."`
	TLSCAKey                   string        `yaml:"tlsca_key"                      help:"The path to the CA key for generating any missing certificates. Relative paths are interpreted as relative to the config path."`
	ClientTLSCert              string        `yaml:"tlscert"                        help:"The path to the client cert for verification. Relative paths are interpreted as relative to the config path."`
	ClientTLSKey               string        `yaml:"tlskey"                         help:"The path to the client key for verification. Relative paths are interpreted as relative to the config path."`
	ServerTLSCert              string        `yaml:"buildkitd_tlscert"              help:"The path to the server cert for verification. Relative paths are interpreted as relative to the config path. Only used when Earthly manages buildkit."`
	ServerTLSKey               string        `yaml:"buildkitd_tlskey"               help:"The path to the server key for verification. Relative paths are interpreted as relative to the config path. Only used when Earthly manages buildkit."`
	TLSEnabled                 bool          `yaml:"tls_enabled"                    help:"If TLS should be used to communicate with Buildkit. Only honored when BuildkitScheme is 'tcp'."`
	ContainerFrontend          string        `yaml:"container_frontend"             help:"What program should be used to start and stop buildkitd, save images. Default is 'docker'. Valid options are 'docker' and 'podman' (experimental)."`
	IPTables                   string        `yaml:"ip_tables"                      help:"Which iptables binary to use. Valid values are iptables-legacy or iptables-nft. Bypasses any autodetection."`
	SecretProvider             string        `yaml:"secret_provider"                help:"Command to execute to retrieve secret."`
	GitImage                   string        `yaml:"git_image"                      help:"Image used to resolve git repositories"`

	// Obsolete.
	CachePath      string `yaml:"cache_path"         help:" *Deprecated* The path to keep Earthly's cache."`
	BuildkitScheme string `yaml:"buildkit_transport" help:" *Deprecated* Change how Earthly communicates with its buildkit daemon. Valid options are: docker-container, tcp. TCP is experimental."`
}

// GitConfig contains git-specific config values
type GitConfig struct {
	// these are used for git vendors (e.g. github, gitlab)
	Pattern               string `yaml:"pattern"                      help:"A regular expression defined to match git URLs, defaults to the regex: <site>/([^/]+)/([^/]+). For example if the site is github.com, then the default pattern will match github.com/<user>/<repo>."`
	Substitute            string `yaml:"substitute"                   help:"If specified, a regular expression substitution will be preformed to determine which URL is cloned by git. Values like $1, $2, ... will be replaced with matched subgroup data. If no substitute is given, a URL will be created based on the requested SSH authentication mode."`
	Suffix                string `yaml:"suffix"                       help:"The git repository suffix, like .git."`                                       // .git
	Auth                  string `yaml:"auth"                         help:"What authentication method do you use? Valid options are: http, https, ssh."` // http, https, ssh
	User                  string `yaml:"user"                         help:"The username to use when auth is set to git or https."`
	Port                  int    `yaml:"port"                         help:"The port to connect to when using git; has no effect for http(s)."`
	Prefix                string `yaml:"prefix"                       help:"This path is prefixed to the git clone url, e.g. ssh://user@host:port/prefix/project/repo.git"`
	Password              string `yaml:"password"                     help:"The https password to use when auth is set to https. This setting is ignored when auth is ssh."`
	ServerKey             string `yaml:"serverkey"                    help:"SSH fingerprints, like you would add in your known hosts file, or get from ssh-keyscan."`
	StrictHostKeyChecking *bool  `yaml:"strict_host_key_checking"     help:"Allow ssh access to hosts with unknown server keys (e.g. no entries in known_hosts), defaults to true."`
	SSHCommand            string `yaml:"ssh_command"                  help:"Set a value for the core.sshCommand git config option, which allows you to provide custom SSH configuration."`
}

// Config contains user's configuration values from ~/earthly/config.yml
type Config struct {
	Global GlobalConfig         `yaml:"global"    help:"Global configuration object. Requires YAML literal to set directly."`
	Git    map[string]GitConfig `yaml:"git"       help:"Git configuration object. Requires YAML literal to set directly."`
}

// PortOffset is the offset to use for dev ports.
func PortOffset(installationName string) int {
	if installationName == "earthly" {
		// No offset for the official release.
		return 0
	}
	return 10 + int(crc32.ChecksumIEEE([]byte(installationName)))%1000
}

// ParseYAML parse config data in yaml format.
func ParseYAML(yamlData []byte, installationName string) (Config, error) {
	defaultLocalRegistryPort := DefaultLocalRegistryPort + PortOffset(installationName)
	// prepopulate defaults
	config := Config{
		Global: GlobalConfig{
			BuildkitCacheSizeMb:     0,
			BuildkitCacheSizePct:    0,
			LocalRegistryHost:       fmt.Sprintf("tcp://127.0.0.1:%d", defaultLocalRegistryPort),
			DarwinProxyImage:        DefaultDarwinProxyImage,
			DarwinProxyWait:         DefaultDarwinProxyWait,
			BuildkitScheme:          DefaultBuildkitScheme,
			BuildkitRestartTimeoutS: 60,
			ConversionParallelism:   DefaultConversionParallelism,
			BuildkitMaxParallelism:  DefaultBuildkitMaxParallelism,
			BuildkitAdditionalArgs:  []string{},
			TLSEnabled:              true,
			TLSCAKey:                DefaultCAKey,
			TLSCACert:               DefaultCACert,
			ClientTLSCert:           DefaultClientTLSCert,
			ClientTLSKey:            DefaultClientTLSKey,
			ServerTLSCert:           DefaultServerTLSCert,
			ServerTLSKey:            DefaultServerTLSKey,
			ContainerFrontend:       DefaultContainerFrontend,
		},
	}

	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse YAML config")
	}

	if config.Git == nil {
		config.Git = make(map[string]GitConfig)
	}

	if err := parseRelPaths(installationName, &config); err != nil {
		return Config{}, errors.Wrap(err, "failed to parse relative path")
	}

	return config, nil
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

// Upsert adds or modifies the key to be the specified value.
// This is saved to disk in your earthly config file.
func Upsert(config []byte, path, value string) ([]byte, error) {
	base := &yaml.Node{}
	err := yaml.Unmarshal(config, base)
	if err != nil || base.IsZero() {
		// Possibly an empty file, or a simple comment results in a null document.
		// Not handled well, so manufacture somewhat acceptable document
		fullDoc := string(config) + "\n---"
		otherErr := yaml.Unmarshal([]byte(fullDoc), base)
		if otherErr != nil {
			// Give up.
			if err != nil {
				return []byte{}, errors.Wrapf(err, "failed to parse config file")
			}
			return []byte{}, errors.Wrapf(otherErr, "failed to parse config file")
		}
		base.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}

	pathParts := splitPath(path)

	t, _, err := validatePath(reflect.TypeOf(Config{}), pathParts)
	if err != nil {
		return []byte{}, errors.Wrap(err, "path is not valid")
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

// Delete removes the key and value at the specified path.
// If no key/value exists, the function will eventually return cleanly.
func Delete(config []byte, path string) ([]byte, error) {
	base := &yaml.Node{}
	err := yaml.Unmarshal(config, base)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to parse config file")
	}
	if base.IsZero() {
		return nil, errors.New("config is empty or missing")
	}

	pathParts := splitPath(path)

	_, _, err = validatePath(reflect.TypeOf(Config{}), pathParts)
	if err != nil {
		return []byte{}, errors.Wrap(err, "path is not valid")
	}

	deleteYamlValue(base, pathParts)

	newConfig, err := yaml.Marshal(base)
	if err != nil {
		return []byte{}, err
	}

	return newConfig, nil
}

// PrintHelp describes the provided config option by
// printing its type and help tags to the console.
func PrintHelp(path string) error {
	t, help, err := validatePath(reflect.TypeOf(Config{}), splitPath(path))
	if err != nil {
		return errors.Wrapf(err, "'%s' is not a valid config value", path)
	}
	fmt.Printf("(%s): %s\n", t.Kind(), help)
	return nil
}

func splitPath(path string) []string {
	// Allow quotes to group keys, since git repos are keys and have periods... this is why we don't just strings.Split
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
	err := yaml.Unmarshal([]byte(value), valueNode)
	if err != nil {
		return nil, errors.Errorf("%s is not a valid YAML value", value)
	}

	// Unfold all the yaml so it's not mixed inline and flow styles in the final document
	var fixStyling func(node *yaml.Node)
	fixStyling = func(node *yaml.Node) {
		node.Style = 0

		for _, n := range node.Content {
			fixStyling(n)
		}
	}
	fixStyling(valueNode)

	contentNode := &yaml.Node{}
	if len(valueNode.Content) > 0 {
		// ContentNode contains the user-provided value with its type etc
		contentNode = valueNode.Content[0]
	} else if value == "" {
		// Edge case where the yaml.Unmarshal above results in no nodes in valueNode.Content.
		// The code below ensures we can write an actual empty string to our yaml as requested.
		contentNode.SetString("")
	} else {
		// Very unlikely
		return nil, errors.New("failed setting value in yaml")
	}

	return contentNode, nil
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

func deleteYamlValue(node *yaml.Node, path []string) []string {

	switch node.Kind {
	case yaml.DocumentNode:
		for _, c := range node.Content {
			path = deleteYamlValue(c, path)
		}

	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			// Keys/Values are inline. Count by twos to get it right.
			key := node.Content[i]
			val := node.Content[i+1]

			if len(path) > 0 && key.Value == path[0] {
				path = path[1:]

				// We found the correct key/value pair.
				// Build new Content without those nodes.
				if len(path) == 0 {
					var newContent []*yaml.Node
					for j, n := range node.Content {
						if j != i && j != i+1 {
							newContent = append(newContent, n)
						}
					}
					node.Content = newContent
					return []string{}
				}

				path = deleteYamlValue(val, path)
			}
		}

	default: // Sequence, Scalar nodes get skipped
		return path
	}

	return []string{}
}

// ReadConfigFile reads in the config file from the disk, into a byte slice.
func ReadConfigFile(configPath string) ([]byte, error) {
	yamlData, err := os.ReadFile(configPath)
	if err != nil {
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

	return os.WriteFile(configPath, data, 0644)
}

func parseRelPaths(instName string, cfg *Config) error {
	if err := parseTLSPaths(instName, cfg); err != nil {
		return errors.Wrap(err, "could not parse relative TLS paths")
	}
	return nil
}

func parseTLSPaths(instName string, cfg *Config) error {
	if !cfg.Global.TLSEnabled {
		return nil
	}
	fields := map[string]*string{
		"ca key":      &cfg.Global.TLSCAKey,
		"ca cert":     &cfg.Global.TLSCACert,
		"client key":  &cfg.Global.ClientTLSKey,
		"client cert": &cfg.Global.ClientTLSCert,
		"server key":  &cfg.Global.ServerTLSKey,
		"server cert": &cfg.Global.ServerTLSCert,
	}
	for name, field := range fields {
		if err := parsePath(instName, field); err != nil {
			return errors.Wrapf(err, "could not parse %v path %q", name, *field)
		}
	}
	return nil
}

func parsePath(instName string, field *string) error {
	if field == nil {
		return errors.New("cannot parse nil field")
	}
	newPath, err := cfgPath(instName, *field)
	if err != nil {
		return err
	}
	*field = newPath
	return nil
}

func cfgPath(instName, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	cfgDir, err := cliutil.GetOrCreateEarthlyDir(instName)
	if err != nil {
		return "", err
	}
	return filepath.Join(cfgDir, path), nil
}
