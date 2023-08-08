package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/util/cliutil"
	"github.com/earthly/earthly/util/fileutil"
)

func (app *earthlyApp) newCloudClient(opts ...cloud.ClientOpt) (*cloud.Client, error) {
	cloudClient, err := cloud.NewClient(app.cloudHTTPAddr, app.cloudGRPCAddr, app.cloudGRPCInsecure, app.sshAuthSock,
		app.authToken, app.authJWT, app.installationName, app.requestID, app.console.Warnf, app.serverConnTimeout, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cloud client")
	}
	return cloudClient, nil
}

// getCIHost returns protocol://hostname
func (app *earthlyApp) getCIHost() string {
	switch {
	case strings.Contains(app.cloudGRPCAddr, "staging"):
		return "https://cloud.staging.earthly.dev"
	case strings.Contains(app.cloudGRPCAddr, "earthly.local"):
		return "http://earthly.local:3000"
	}
	return "https://cloud.earthly.dev"
}

func wrap(s ...string) string {
	return strings.Join(s, "\n\t")
}

func promptInput(ctx context.Context, question string) (string, error) {
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

func processSecrets(secrets, secretFiles []string, dotEnvMap map[string]string) (map[string][]byte, error) {
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
				return nil, errors.Errorf("env var %s not set", secret)
			}
			data = []byte(value)
		}
		if _, ok := finalSecrets[key]; ok {
			return nil, errors.Errorf("secret %q already contains a value", key)
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
			return nil, errors.Errorf("secret %q already contains a value", k)
		}
		finalSecrets[k] = []byte(data)
	}
	return finalSecrets, nil
}

func defaultConfigPath(installName string) string {
	earthlyDir := cliutil.GetEarthlyDir(installName)
	oldConfig := filepath.Join(earthlyDir, "config.yaml")
	newConfig := filepath.Join(earthlyDir, "config.yml")
	oldConfigExists, _ := fileutil.FileExists(oldConfig)
	newConfigExists, _ := fileutil.FileExists(newConfig)
	if oldConfigExists && !newConfigExists {
		return oldConfig
	}
	return newConfig
}

func isEarthlyBinary(path string) bool {
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

func profhandler() {
	addr := "127.0.0.1:6060"
	fmt.Printf("listening for pprof on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("error listening for pprof: %v", err)
	}
}

type byteSizeValue uint64

func (b *byteSizeValue) Set(s string) error {
	v, err := humanize.ParseBytes(s)
	if err != nil {
		return err
	}
	*b = byteSizeValue(v)
	return nil
}

func (b *byteSizeValue) String() string { return humanize.Bytes(uint64(*b)) }
