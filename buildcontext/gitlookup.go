package buildcontext

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/fileutil"
)

type gitMatcher struct {
	name     string
	re       *regexp.Regexp
	sub      string
	user     string
	suffix   string
	protocol string
	password string
	keyScan  string
}

// GitLookup looksup gits
type GitLookup struct {
	mu            sync.Mutex
	matchers      []*gitMatcher
	catchAll      *gitMatcher
	autoProtocols map[string]string // host -> detected protocol type
	sshAuthSock   string
	console       conslogging.ConsoleLogger
}

// NewGitLookup creates new lookuper
func NewGitLookup(console conslogging.ConsoleLogger, sshAuthSock string) *GitLookup {
	matchers := []*gitMatcher{
		{
			name:     "github.com",
			re:       regexp.MustCompile("github.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "auto",
			keyScan:  "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			name:     "gitlab.com",
			re:       regexp.MustCompile("gitlab.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "auto",
			keyScan:  "gitlab.com ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9",
		},
		{
			name:     "bitbucket.com",
			re:       regexp.MustCompile("bitbucket.com/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "auto",
			keyScan:  "bitbucket.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==",
		},
	}

	gl := &GitLookup{
		matchers: matchers,
		catchAll: &gitMatcher{
			name:     "",
			re:       regexp.MustCompile("[^/]+/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: "auto",
		},
		autoProtocols: map[string]string{},
		sshAuthSock:   sshAuthSock,
		console:       console,
	}
	return gl
}

// ErrNoMatch occurs when no git matcher is found
var ErrNoMatch = errors.Errorf("no git match found")

// DisableSSH changes all git matchers from ssh to https
func (gl *GitLookup) DisableSSH() {
	gl.mu.Lock()
	defer gl.mu.Unlock()
	for i, m := range gl.matchers {
		if m.protocol == "ssh" || m.protocol == "auto" {
			gl.matchers[i].protocol = "https"
		}
	}
	if gl.catchAll.protocol == "ssh" {
		gl.catchAll.protocol = "https"
	}
}

// AddMatcher adds a new matcher for looking up git repos
func (gl *GitLookup) AddMatcher(name, pattern, sub, user, password, suffix, protocol, keyScan string) error {
	gl.mu.Lock()
	defer gl.mu.Unlock()
	if protocol == "http" && password != "" {
		return errors.Errorf("using a password with http for %s is insecure", name)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "failed to compile regex %s", pattern)
	}
	switch protocol {
	case "http", "https", "ssh", "auto":
		break
	default:
		return errors.Errorf("unsupported git protocol %q", protocol)
	}

	gm := &gitMatcher{
		name:     name,
		re:       re,
		sub:      sub,
		user:     user,
		password: password,
		suffix:   suffix,
		protocol: protocol,
		keyScan:  keyScan,
	}

	// update existing entry
	for i, m := range gl.matchers {
		if m.name == name {
			if gm.keyScan == "" {
				gm.keyScan = m.keyScan
			}
			gl.matchers[i] = gm
			return nil
		}
	}

	// add new entry
	gl.matchers = append(gl.matchers, gm)
	return nil
}

func (gl *GitLookup) hostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	for _, m := range gl.matchers {
		k, _, _, _, err := ssh.ParseAuthorizedKey([]byte(m.keyScan))
		if err != nil {
			gl.console.Warnf("failed to parse authorized key %q", m.keyScan)
			continue
		}
		if k.Type() == key.Type() && bytes.Equal(k.Marshal(), key.Marshal()) {
			return nil
		}
	}

	hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return err
	}
	return hostKeyCallback(hostname, remote, key)
}

func (gl *GitLookup) getGitMatcher(path string) (string, *gitMatcher, error) {
	if len(gl.matchers) == 0 {
		panic("no matchers")
	}
	for _, m := range gl.matchers {
		match := m.re.FindString(path)
		if match != "" {
			return match, m, nil
		}
	}

	match := gl.catchAll.re.FindString(path)
	if match != "" {
		return match, gl.catchAll, nil
	}

	return "", nil, ErrNoMatch
}

// detectProtocol will update the gitMatcher protocol if it is set to auto
func (gl *GitLookup) detectProtocol(host string) (protocol string, err error) {
	var ok bool
	protocol, ok = gl.autoProtocols[host]
	if ok {
		return
	}

	defer func() {
		if err == nil {
			gl.autoProtocols[host] = protocol
		}
	}()

	sshAgent, err := net.Dial("unix", gl.sshAuthSock)
	if err != nil {
		protocol = "https"
		err = nil
		return
	}
	config := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers),
		},
		HostKeyCallback: gl.hostKeyCallback,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		protocol = "https"
		err = nil
		return
	}
	defer client.Close()

	protocol = "ssh"
	err = nil
	return
}

// GetCloneURL returns the repo to clone, and a path relative to the repo
//   "github.com/earthly/earthly"             ---> ("git@github.com/earthly/earthly.git", "")
//   "github.com/earthly/earthly/examples"    ---> ("git@github.com/earthly/earthly.git", "examples")
//   "github.com/earthly/earthly/examples/go" ---> ("git@github.com/earthly/earthly.git", "examples/go")
// Additionally a ssh keyscan might be returned (or an empty string indicating none was configured)
func (gl *GitLookup) GetCloneURL(path string) (string, string, string, error) {
	gl.mu.Lock()
	defer gl.mu.Unlock()
	match, m, err := gl.getGitMatcher(path)
	if err != nil {
		return "", "", "", err
	}

	n := len(match)
	subPath := ""
	if len(path) > n {
		subPath = path[n+1:]
		path = path[:n]
	}
	host := path[:strings.IndexByte(path, '/')]

	protocol := m.protocol
	if protocol == "auto" {
		protocol, err = gl.detectProtocol(host)
		if err != nil {
			return "", "", "", err
		}
	}

	var gitURL, keyScan string
	switch protocol {
	case "ssh":
		gitURL = m.user + "@" + strings.Replace(match, "/", ":", 1) + m.suffix
		keyScan = m.keyScan
	case "http", "https":
		var userAndPass string
		if m.user != "" && m.password != "" {
			userAndPass = url.QueryEscape(m.user) + ":" + url.QueryEscape(m.password) + "@"
		}
		gitURL = protocol + "://" + userAndPass + match + m.suffix
	default:
		return "", "", "", errors.Errorf("unsupported protocol: %s", protocol)
	}

	if m.sub != "" {
		if !m.re.MatchString(path) {
			return "", "", "", errors.Errorf("failed to determine git path to clone for %q", path)
		}
		gitURL = m.re.ReplaceAllString(path, m.sub)
	}

	if keyScan == "" {
		keyScan, err = loadKnownHosts()
		if err != nil {
			return "", "", "", err
		}
	}

	return gitURL, subPath, keyScan, nil
}

func loadKnownHosts() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get user home dir")
	}

	knownHosts := filepath.Join(homeDir, ".ssh/known_hosts")

	if !fileutil.FileExists(knownHosts) {
		return "", nil
	}

	b, err := ioutil.ReadFile(knownHosts)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %s", knownHosts)
	}
	return string(b), nil
}
