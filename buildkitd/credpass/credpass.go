package credpass

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// DockerCredsGet is the format of the <docker-creds-helper> get command.
type DockerCredsGet struct {
	ServerURL string `json:"ServerURL"`
	Username  string `json:"Username"`
	Secret    string `json:"Secret"`
}

// DockerConfig represents the ~/.docker/config.json file format
// (minus several items that are not important to us).
type DockerConfig struct {
	Auths      map[string]DockerConfigAuth `json:"auths"`
	CredsStore string                      `json:"credsStore,omitempty"`
}

// DockerConfigAuth is one auth entry of a docker config json structure.
type DockerConfigAuth struct {
	Auth string `json:"auth"`
}

// Read reads docker credentials using configured docker credentials store, if any available
// and returns a newly generated docker config file contents with credentials stored as base64. This
// file is not stored on disk - its generated contents are returned as a string (first return arg).
// If credentials store is not set, assume credentials are plain text, in which case, return
// path to docker config dir - the file can be found there (second return arg).
func Read(ctx context.Context) (string, string, error) {
	credsHelper, err := credsHelperBinary()
	if err != nil {
		return "", "", err
	}
	if credsHelper == "" {
		// No creds store specified. Not encrypted - can simply mount docker config file within
		// container.
		configDir, err := dockerConfigDir()
		if err != nil {
			return "", "", err
		}
		if configDir == "" {
			// Docker config dir does not exist. Don't bother then.
			return "", "", nil
		}
		return "", configDir, nil
	}

	// Run list command to get all logged-in registries.
	listCmd := exec.CommandContext(ctx, credsHelper, "list")
	listBytes, err := listCmd.Output()
	if err != nil {
		return "", "", errors.Wrapf(err, "error running %s list", credsHelper)
	}
	credsMap := make(map[string]string)
	err = json.Unmarshal(listBytes, &credsMap)
	if err != nil {
		return "", "", errors.Wrapf(err, "error unmarshalling output of %s list command", credsHelper)
	}
	// For each registry, fetch username and password
	var credsList []DockerCredsGet
	for serverURL := range credsMap {
		getCmd := exec.CommandContext(ctx, credsHelper, "get")
		serverURLReader := strings.NewReader(serverURL)
		getCmd.Stdin = serverURLReader
		getBytes, err := getCmd.Output()
		if err != nil {
			return "", "", errors.Wrapf(err, "failed to get creds for registry %s", serverURL)
		}
		var creds DockerCredsGet
		err = json.Unmarshal(getBytes, &creds)
		if err != nil {
			return "", "", errors.Wrapf(err, "failed to read json output of command %s get", credsHelper)
		}
		credsList = append(credsList, creds)
	}
	// Build out a DockerConfig JSON with credentials base64 encoded.
	dc := DockerConfig{
		Auths: make(map[string]DockerConfigAuth),
	}
	for _, creds := range credsList {
		credsPlainText := fmt.Sprintf("%s:%s", creds.Username, creds.Secret)
		credsBase64 := base64.StdEncoding.EncodeToString([]byte(credsPlainText))
		dc.Auths[creds.ServerURL] = DockerConfigAuth{
			Auth: credsBase64,
		}
	}
	dcBytes, err := json.Marshal(dc)
	if err != nil {
		return "", "", errors.Wrap(err, "error marshalling new docker config")
	}
	return string(dcBytes), "", nil
}

func credsHelperBinary() (string, error) {
	cs, err := readCredsStore()
	if err != nil {
		return "", err
	}
	if cs == "" {
		return "", nil
	}
	credsHelper := fmt.Sprintf("docker-credential-%s", cs)
	_, err = exec.LookPath(credsHelper)
	if err != nil {
		return "", errors.Wrapf(err, "could not find docker credentials helper %s", credsHelper)
	}
	return credsHelper, nil
}

func readCredsStore() (string, error) {
	configDir, err := dockerConfigDir()
	if err != nil {
		return "", err
	}
	if configDir == "" {
		return "", nil
	}
	configFile := path.Join(configDir, "config.json")
	configFileBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read docker config file %s", configFile)
	}
	var dc DockerConfig
	err = json.Unmarshal(configFileBytes, &dc)
	if err != nil {
		return "", errors.Wrapf(err, "unable to unmarshal json file %s", configFile)
	}
	return dc.CredsStore, nil
}

func dockerConfigDir() (string, error) {
	dir := os.Getenv("DOCKER_CONFIG")
	if dir == "" {
		homeDir := os.Getenv("HOME")
		if homeDir == "" {
			homeDir = "/root"
		}
		dir = path.Join(homeDir, ".docker")
	}
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "unable to stat directory %s", dir)
	}
	return dir, nil
}
