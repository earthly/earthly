package common

// Only functions that do NOT touch the app CLI should go here!

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/hint"
	"github.com/earthly/earthly/variables"
	gsysinfo "github.com/elastic/go-sysinfo"
	"github.com/pkg/errors"
)

func Wrap(s ...string) string {
	return strings.Join(s, "\n\t")
}

func CombineVariables(dotEnvMap map[string]string, flagArgs []string, buildFlagArgs []string) (*variables.Scope, error) {
	dotEnvVars := variables.NewScope()
	for k, v := range dotEnvMap {
		dotEnvVars.Add(k, variables.NewStringVariable(v))
	}
	buildArgs := append([]string{}, buildFlagArgs...)
	buildArgs = append(buildArgs, flagArgs...)
	overridingVars, err := variables.ParseCommandLineArgs(buildArgs)
	if err != nil {
		return nil, errors.Wrap(err, "parse build args")
	}
	return variables.CombineScopes(overridingVars, dotEnvVars), nil
}

func ProcessSecrets(secrets, secretFiles []string, dotEnvMap map[string]string, secretsFilePath string) (map[string][]byte, error) {
	finalSecrets := make(map[string][]byte)
	for k, v := range dotEnvMap {
		finalSecrets[k] = []byte(v)
	}
	for _, secret := range secrets {
		parts := strings.SplitN(secret, "=", 2)
		key := parts[0]
		var data []byte
		if len(parts) == 2 {
			// secret value passed as argument
			data = []byte(parts[1])
		} else {
			// Not set. Use environment to fetch it.
			value, found := os.LookupEnv(secret)
			if !found {
				return nil, hint.Wrapf(errors.Errorf("failed to set secret %q via --secret flag without a value", secret), "Try to set an env var by the name %q with the secret value or pass the value as part of the --secret flag", secret)
			}
			data = []byte(value)
		}
		if _, ok := finalSecrets[key]; ok {
			return nil, hint.Wrapf(errors.Errorf("failed to set secret %q via --secret flag", key), "Check the secret %q has not already been set in the file %q or passed more than once to the command", key, secretsFilePath)
		}
		finalSecrets[key] = data
	}
	for _, secret := range secretFiles {
		parts := strings.SplitN(secret, "=", 2)
		if len(parts) != 2 {
			return nil, errors.Errorf("unable to parse --secret-file argument: %q", secret)
		}
		k := parts[0]
		path := fileutil.ExpandPath(parts[1])
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open %q", path)
		}
		if _, ok := finalSecrets[k]; ok {
			return nil, hint.Wrapf(errors.Errorf("failed to set secret %q via --secret-file flag", k), "Check the secret %q has not already been set in the file %q, or passed via --secret flag", k, secretsFilePath)
		}
		finalSecrets[k] = data
	}
	return finalSecrets, nil
}

func GetPlatform() string {
	h, err := gsysinfo.Host()
	if err != nil {
		return "unknown"
	}
	info := h.Info()
	return fmt.Sprintf("%s/%s; %s %s", runtime.GOOS, runtime.GOARCH, info.OS.Name, info.OS.Version)
}

func GetBinaryName() string {
	if len(os.Args) == 0 {
		return "earthly"
	}
	binPath := os.Args[0] // can't use os.Executable() here; because it will give us earthly if executed via the earth symlink
	baseName := path.Base(binPath)
	return baseName
}

func PromptInput(ctx context.Context, question string) (string, error) {
	fmt.Printf("%s", question)
	var line string
	var readErr error
	ch := make(chan struct{})
	go func() {
		rbuf := bufio.NewReader(os.Stdin)
		line, readErr = rbuf.ReadString('\n')
		close(ch)
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-ch:
		if readErr != nil {
			return "", readErr
		}
		return strings.TrimRight(line, "\n"), nil
	}
}

func IfNilBoolDefault(ptr *bool, defaultValue bool) bool {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}

func IsEarthlyBinary(path string) bool {
	// apply heuristics to see if binary is a version of earthly
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	if !bytes.Contains(data, []byte("docs.earthly.dev")) {
		return false
	}
	if !bytes.Contains(data, []byte("api.earthly.dev")) {
		return false
	}
	if !bytes.Contains(data, []byte("Earthfile")) {
		return false
	}
	return true
}
