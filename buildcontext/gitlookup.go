package buildcontext

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/exp/maps"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/fileutil"
	"github.com/earthly/earthly/util/stringutil"

	"github.com/jdxcode/netrc"
	"github.com/moby/buildkit/util/sshutil"
)

type gitMatcher struct {
	name                  string
	re                    *regexp.Regexp
	sub                   string
	user                  string
	suffix                string
	protocol              gitProtocol
	password              string
	strictHostKeyChecking bool
	port                  int
	prefix                string
	sshCommand            string
}

type gitProtocol string

const (
	autoProtocol  gitProtocol = "auto"
	sshProtocol   gitProtocol = "ssh"
	httpProtocol  gitProtocol = "http"
	httpsProtocol gitProtocol = "https"
)

// GitLookup looksup gits
type GitLookup struct {
	mu            sync.Mutex
	matchers      []*gitMatcher
	catchAll      *gitMatcher
	autoProtocols map[string]gitProtocol // host -> detected protocol type
	sshAuthSock   string
	keyScans      []string
	console       conslogging.ConsoleLogger
}

var defaultKeyScans = []string{
	// github.com
	"|1|+wkzm0y4RAEaLjnuB3lvMyNmqto=|A97DLdg1fwTjawL47CHJqEeE2lw= ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCj7ndNxQowgcQnjshcLrqPEiiphnt+VTTvDP6mHBL9j1aNUkY4Ue1gvwnGLVlOhGeYrnZaMgRK6+PKCUXaDbC7qtbW8gIkhL7aGCsOr/C56SJMy/BCZfxd1nWzAOxSDPgVsmerOBYfNqltV9/hWCqBywINIR+5dIg6JTJ72pcEpEjcYgXkE2YEFXV1JHnsKgbLWNlhScqb2UmyRkQyytRLtL+38TGxkxCflmO+5Z8CSSNY7GidjMIZ7Q4zMjA2n1nGrlTDkzwDCsw+wqFPGQA179cnfGWOWRVruj16z6XyvxvjJwbz0wQZ75XK5tKSb7FNyeIEs4TT4jk+S4dhPeAUC5y+bDYirYgM4GC7uEnztnZyaVWQ7B381AK4Qdrwt51ZqExKbQpTUNn+EjqoTwvqNj4kqx5QUCI0ThS/YkOxJCXmPUWZbhjpCg56i+2aB6CmK2JGhn57K5mj0MNdBXA4/WnwH6XoPWJzK5Nyu2zB3nAZp+S5hpQs+p1vN1/wsjk=",
	"|1|tNJE6wQBmC1c4lJm0wtToe8IHxY=|I7K0Cre2i8VXpAKTS2P6Y7bIdqg= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
	"|1|z0g6bpSrCXjh1vZdfzQP634n7SQ=|Xf+7/COPFwsdLXxWptK2/jRP2k0= ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",

	// gitlab.com
	"|1|an1urLLW36WT6FnJoB5BWqVwiEM=|RTcVDky6WhU+S+09yjALNiS4neo= ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsj2bNKTBSpIYDEGk9KxsGh3mySTRgMtXL583qmBpzeQ+jqCMRgBqB98u3z++J1sKlXHWfM9dyhSevkMwSbhoR8XIq/U0tCNyokEi/ueaBMCvbcTHhO7FcwzY92WK4Yt0aGROY5qX2UKSeOvuP4D6TPqKF1onrSzH9bx9XUf2lEdWT/ia1NEKjunUqu1xOB/StKDHMoX4/OKyIzuS0q/T1zOATthvasJFoPrAjkohTyaDUz2LN5JoH839hViyEG82yB+MjcFV5MU3N1l1QL3cVUCh93xSaua1N85qivl+siMkPGbO5xR/En4iEY6K2XPASUEMaieWVNTRCtJ4S8H+9",
	"|1|z2nqpVA8ymA2aLuV3ig57xKYDOw=|2JC7T/Oek2fpc/rw+YOfolDdDCI= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFSMqzJeV9rUzU4kWitGjeR4PWSa29SPqJ1fVkhtj3Hw9xjLVXVYrU9QlYWrOLXBpQ6KWjbjTDTdDkoohFzgbEY=",
	"|1|JAhjb/FmPaOSwPtfZlOYRmq7nlg=|MysQCX5GQaSfAKTn5R5AHdskAt4= ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAfuCHKVTjquxvt6CM6tdG4SLp1Btn/nOeHHE5UOzRdf",

	// bitbucket.com
	"|1|5myLBXBnkK609Pb0DTrYhK9hn3k=|7wQiytbsZpu1pDE7AOs7pfBw/4M= ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==",
}

// NewGitLookup creates new lookuper
func NewGitLookup(console conslogging.ConsoleLogger, sshAuthSock string) *GitLookup {
	gl := &GitLookup{
		catchAll: &gitMatcher{
			name:     "",
			re:       regexp.MustCompile("[^/]+/[^/]+/[^/]+"),
			user:     "git",
			suffix:   ".git",
			protocol: autoProtocol,
		},
		autoProtocols: map[string]gitProtocol{},
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
		if m.protocol == sshProtocol || m.protocol == autoProtocol {
			gl.matchers[i].protocol = httpsProtocol
		}
	}
	if gl.catchAll.protocol == sshProtocol {
		gl.catchAll.protocol = httpsProtocol
	}
}

func knownHostsToKeyScans(knownHosts string) []string {
	knownHosts = strings.ReplaceAll(knownHosts, "\r\n", "\n")
	foundKeyScans := make(map[string]bool)
	for _, s := range strings.Split(knownHosts, "\n") {
		s = strings.TrimSpace(s)
		if s != "" && !strings.HasPrefix(s, "#") && !foundKeyScans[s] {
			foundKeyScans[s] = true
		}
	}
	return maps.Keys(foundKeyScans)
}

// AddMatcher adds a new matcher for looking up git repos
func (gl *GitLookup) AddMatcher(name, pattern, sub, user, password, prefix, suffix, protocol, knownHosts string, strictHostKeyChecking bool, port int, sshCommand string) error {
	gl.mu.Lock()
	defer gl.mu.Unlock()
	p := gitProtocol(protocol)
	if p == httpProtocol && password != "" {
		return errors.Errorf("using a password with http for %s is insecure", name)
	}

	if sub != "" && (port != 0 || prefix != "") {
		return errors.Errorf("unable to use substitution in combination with port or prefix values for %s git config", name)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return errors.Wrapf(err, "failed to compile regex %s", pattern)
	}
	switch p {
	case httpProtocol, httpsProtocol, sshProtocol, autoProtocol:
		break
	default:
		return errors.Errorf("unsupported git protocol %q", protocol)
	}

	gm := &gitMatcher{
		name:                  name,
		re:                    re,
		sub:                   sub,
		user:                  user,
		password:              password,
		port:                  port,
		prefix:                prefix,
		suffix:                suffix,
		protocol:              p,
		strictHostKeyChecking: strictHostKeyChecking,
		sshCommand:            sshCommand,
	}

	// update existing entry
	for i, m := range gl.matchers {
		if m.name == name {
			gl.matchers[i] = gm
			return nil
		}
	}

	gl.keyScans = append(gl.keyScans, knownHostsToKeyScans(knownHosts)...)

	// add new entry
	gl.matchers = append(gl.matchers, gm)
	return nil
}

// from crypto/ssh
// See https://android.googlesource.com/platform/external/openssh/+/ab28f5495c85297e7a597c1ba62e996416da7c7e/hostfile.c#120
func hashHost(hostname string, salt []byte) []byte {
	mac := hmac.New(sha1.New, salt)
	mac.Write([]byte(hostname))
	return mac.Sum(nil)
}

var errUnsupportedHash = errors.New("unsupported keyscan hash")
var errInvalidScan = errors.New("invalid keyscan")
var errKeyScanNoMatch = errors.New("keyscan does not match")

func hasPort(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '[' {
		return strings.Contains(s, "]:")
	}
	return strings.Contains(s, ":")
}

func isHashedHost(hashAndSalt, hostname string) (bool, error) {
	prefix := "|1|"
	if !strings.HasPrefix(hashAndSalt, prefix) {
		return false, errUnsupportedHash
	}

	splits := strings.Split(hashAndSalt[len(prefix):], "|")
	if len(splits) != 2 {
		return false, errInvalidScan
	}

	salt, err := base64.StdEncoding.DecodeString(splits[0])
	if err != nil {
		return false, errors.Wrap(err, "failed to decode known_hosts salt")
	}
	hash, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		return false, errors.Wrap(err, "failed to decode known_hosts hash")
	}

	hostnameHash := hashHost(hostname, salt)
	ok := bytes.Equal(hostnameHash, hash)

	if !ok {
		// try hashing the hostname only (ssh-keyscan -p 2222 -H hostname doesn't include the port in the hash)
		// ses bugfix: https://github.com/openssh/openssh-portable/commit/e9c71498a083a8b502aa831ea931ce294228eda0
		if hasPort(hostname) {
			host, _, err := net.SplitHostPort(hostname)
			if err != nil {
				return false, errors.Wrapf(err, "SplitHostPort on %q failed", hostname)
			}
			hostnameHash := hashHost(host, salt)
			ok = bytes.Equal(hostnameHash, hash)
		}
	}

	return ok, nil
}

func parseKeyScanIfHostMatches(keyScan, hostname string) (keyAlg, keyData string, err error) {
	hostname = knownhosts.Normalize(hostname)
	splits := strings.Fields(keyScan)
	if len(splits) < 3 {
		err = errInvalidScan
		return
	}
	scannedHostname := splits[0]

	if strings.HasPrefix(scannedHostname, "|") {
		var ok bool
		ok, err = isHashedHost(scannedHostname, hostname)
		if err != nil {
			return
		}
		if !ok {
			err = errKeyScanNoMatch
			return
		}
	} else {
		// entry isn't hashed
		// either the entry is of the form `[hostname]:port` or simply `hostname`
		if scannedHostname != hostname {
			// check for entry without a port
			// TODO: ACB is not sure if this part is needed ( https://github.com/openssh/openssh-portable/commit/e9c71498a083a8b502aa831ea931ce294228eda0 is a bugfix
			// that only affects hashed entries, however, it's not clear if old versions of ssh dropped the port in the non-hashed version).
			if !hasPort(hostname) {
				err = errKeyScanNoMatch
				return
			}
			var host string
			host, _, err = net.SplitHostPort(hostname)
			if err != nil {
				err = errors.Wrapf(err, "SplitHostPort on %q failed", hostname)
				return
			}
			if scannedHostname != host {
				err = errKeyScanNoMatch
				return
			}
		}
	}

	keyAlg = splits[1]
	keyData = splits[2]
	err = nil
	return
}

// This comes from crypto/ssh/common.go
// supportedHostKeyAlgos specifies the supported host-key algorithms (i.e. methods
// of authenticating servers) in preference order.
var supportedHostKeyAlgos = []string{
	ssh.CertAlgoRSAv01, ssh.CertAlgoDSAv01, ssh.CertAlgoECDSA256v01,
	ssh.CertAlgoECDSA384v01, ssh.CertAlgoECDSA521v01, ssh.CertAlgoED25519v01,

	ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521,
	ssh.KeyAlgoRSA, ssh.KeyAlgoDSA,

	ssh.KeyAlgoED25519,
}

func (gl *GitLookup) getHostKeyAlgorithms(hostname string) ([]string, []string, error) {
	foundAlgs := map[string]bool{}

	knownHostsKeyScans, err := loadKnownHosts()
	if err != nil {
		gl.console.Warnf("failed to load ~/.ssh/known_hosts: %s", err)
	}
	gl.console.VerbosePrintf("loaded %d key(s) from known_hosts and %d default key(s)", len(knownHostsKeyScans), len(defaultKeyScans))

	foundKeys := make(map[string]bool)
	for _, keyScans := range [][]string{
		knownHostsKeyScans,
		defaultKeyScans,
	} {
		for _, keyScan := range keyScans {
			keyAlg, keyData, err := parseKeyScanIfHostMatches(keyScan, hostname)
			switch err {
			case nil:
			case errKeyScanNoMatch:
				gl.console.VerbosePrintf("ignoring key scan %q: due to host mismatch", keyScan)
				continue
			default:
				gl.console.Warnf("failed to parse key scan %q: %s", keyScan, err)
				continue
			}
			foundAlgs[keyAlg] = true
			key := fmt.Sprintf("%s %s %s", knownhosts.Normalize(hostname), keyAlg, keyData)
			if !foundKeys[key] {
				gl.console.VerbosePrintf("found (normalized) key %s", key)
				foundKeys[key] = true
			}
		}
	}

	keys := maps.Keys(foundKeys)

	algs := []string{}
	for _, alg := range supportedHostKeyAlgos {
		if _, ok := foundAlgs[alg]; ok {
			algs = append(algs, alg)
		}
	}
	return algs, keys, nil
}

func (gl *GitLookup) newHostKeyCallback(keys []string) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		for _, keyScan := range keys {
			k, _, _, _, err := ssh.ParseAuthorizedKey([]byte(keyScan))
			if err != nil {
				gl.console.Warnf("failed to parse authorized key %q", keyScan)
				continue
			}
			if k.Type() == key.Type() && bytes.Equal(k.Marshal(), key.Marshal()) {
				return nil
			}
		}
		return fmt.Errorf("no known_host entry for %s", hostname)
	}
}

func (gl *GitLookup) getGitMatcherByPath(path string) (string, *gitMatcher, error) {
	for _, m := range gl.matchers {
		match := m.re.FindString(path)
		if match != "" {
			gl.console.VerbosePrintf("matched earthly reference %s with git config entry %s (regex %s)", path, m.name, m.re)
			return match, m, nil
		}
	}
	match := gl.catchAll.re.FindString(path)
	if match != "" {
		gl.console.VerbosePrintf("matched earthly reference %s with pre-configured catch-all (regex %s)", path, gl.catchAll.re)
		return match, gl.catchAll, nil
	}
	gl.console.VerbosePrintf("failed to match earthly reference %s with any git matchers", path)
	return "", nil, ErrNoMatch
}

func (gl *GitLookup) getGitMatcherByName(name string) *gitMatcher {
	for _, m := range gl.matchers {
		if m.name == name {
			gl.console.VerbosePrintf("found git config specific for %s", name)
			return m
		}
	}
	gl.console.VerbosePrintf("no host-specific git config found for %s, using global git settings", name)
	return gl.catchAll
}

// detectProtocol will update the gitMatcher protocol if it is set to auto
func (gl *GitLookup) detectProtocol(host string) (protocol gitProtocol, err error) {
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
		gl.console.VerbosePrintf("failed to connect to ssh-agent (using %s) due to %s; falling back to https", gl.sshAuthSock, err.Error())
		protocol = httpsProtocol
		err = nil
		return
	}

	algs, keys, err := gl.getHostKeyAlgorithms(host)
	if err != nil {
		gl.console.VerbosePrintf("failed to get accepted host key algorithms for %s: %s; falling back to https", host, err.Error())
		protocol = httpsProtocol
		err = nil
		return
	}
	if len(keys) == 0 {
		gl.console.VerbosePrintf("no known_hosts entries found for %s; falling back to https", host)
		protocol = httpsProtocol
		err = nil
		return
	}

	config := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers),
		},
		HostKeyAlgorithms: algs,
		HostKeyCallback:   gl.newHostKeyCallback(keys),
		Timeout:           time.Second * 3,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", host), config)
	if err != nil {
		gl.console.VerbosePrintf("failed to connect to %s over ssh due to %s; falling back to https", host, err.Error())
		protocol = httpsProtocol
		err = nil
		return
	}
	defer client.Close()

	gl.console.VerbosePrintf("defaulting to ssh protocol for %s", host)
	protocol = sshProtocol
	err = nil
	return
}

var errNoRCHostEntry = fmt.Errorf("no netrc host entry")

func (gl *GitLookup) lookupNetRCCredential(host string) (login, password string, err error) {
	var n *netrc.Netrc
	if content := os.Getenv("NETRC_CONTENT"); content != "" {
		n, err = netrc.ParseString(content)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to parse NETRC_CONTENT data")
		}
	} else if path := os.Getenv("NETRC"); path != "" {
		n, err = netrc.Parse(path)
		if err != nil {
			return "", "", errors.Wrapf(err, "failed to parse netrc file: %s", path)
		}
	} else {
		homeDir, _ := fileutil.HomeDir()
		path = filepath.Join(homeDir, ".netrc")
		n, err = netrc.Parse(path)
		if err != nil {
			return "", "", errors.Wrap(err, "failed to parse default .netrc file")
		}
	}
	machine := n.Machine(host)
	if machine == nil {
		return "", "", errors.Wrapf(errNoRCHostEntry, "failed to lookup netrc entry for %s", host)
	}
	login = n.Machine(host).Get("login")
	password = n.Machine(host).Get("password")
	return login, password, nil
}

var errMakeCloneURLSubNotSupported = fmt.Errorf("makeCloneURL does not support gitMatcher substitution")

func (gl *GitLookup) makeCloneURL(m *gitMatcher, host, gitPath string) (gitURL string, keyScans []string, sshCommand string, err error) {
	if m.sub != "" {
		return "", nil, "", errMakeCloneURLSubNotSupported
	}

	configuredProtocol := m.protocol
	user := m.user
	password := m.password
	if configuredProtocol == autoProtocol {
		configuredProtocol, err = gl.detectProtocol(host)
		if err != nil {
			return "", nil, "", err
		}
		switch configuredProtocol {
		case sshProtocol:
			user = "git"
		case httpProtocol, httpsProtocol:
			user = ""
			password = ""
		}
	}

	switch configuredProtocol {
	case sshProtocol:
		if user == "" {
			var ok bool
			user, ok = os.LookupEnv("USER")
			if !ok {
				user = "git"
				gl.console.VerbosePrintf("ssh auth configured without a user; failed to get current user, defaulting to git")
			} else {
				gl.console.VerbosePrintf("ssh auth configured without a user; defaulting to current user")
			}
		}

		port := m.port
		if port == 0 {
			port = 22
		}

		// careful about changing all clone paths to the explicit ssh://user@host:port/user/repo.git form.
		// as the implicit form assumes the repo is relative to the user's home directory.
		// for example "git clone alex@coho:junk/test.git", might work as ssh://alex@coho/~/junk/test.git
		// however not all git servers support the use of `~` to denode the user's repo.
		// For instance, github fails to `git clone ssh://git@github.com/~/user/repo.git`.
		if strings.HasPrefix(gitPath, "/") || port != 22 {
			var portStr string
			if port != 22 {
				portStr = fmt.Sprintf(":%d", port)
			}
			if !strings.HasPrefix(gitPath, "/") {
				gitPath = "/" + gitPath
			}
			gitURL = "ssh://" + user + "@" + host + portStr + gitPath
		} else {
			gitURL = user + "@" + host + ":" + gitPath
		}
		_, keyScans, err = gl.getHostKeyAlgorithms(host)
		if err != nil {
			return "", nil, "", err
		}
		if len(keyScans) == 0 && m.strictHostKeyChecking {
			return "", nil, "", errors.Errorf("no known_hosts entries exist for %s", host)
		}
	case httpProtocol:
		if user != "" || password != "" {
			gl.console.Warnf("%s has been configured to use basic access authentication with http; this is insecure and will be ignored; use https or ssh authentication instead", host)
		}
		gitURL = "http://" + host + "/" + gitPath
	case httpsProtocol:
		var userAndPass string
		if user == "" && password == "" {
			user, password, _ = gl.lookupNetRCCredential(host) // best effort
		}
		if user != "" && password != "" {
			userAndPass = url.QueryEscape(user) + ":" + url.QueryEscape(password) + "@"
		}
		gitURL = "https://" + userAndPass + host + "/" + gitPath
	default:
		return "", nil, "", errors.Errorf("unsupported protocol: %s", configuredProtocol)
	}
	return gitURL, keyScans, m.sshCommand, nil
}

// TODO eventually we should use gitutil.parseURL directly; but for now we want to avoid this change to keep this commit smaller
const (
	HTTPProtocol = iota + 1
	HTTPSProtocol
	SSHProtocol
	GitProtocol
	UnknownProtocol
)

// parseGitProtocol comes from buildkit (which was named ParseProtocol); it was since deleted and replaced with ParseURL)
func parseGitProtocol(remote string) (string, int) {
	prefixes := map[string]int{
		"http://":  HTTPProtocol,
		"https://": HTTPSProtocol,
		"git://":   GitProtocol,
		"ssh://":   SSHProtocol,
	}
	protocolType := UnknownProtocol
	for prefix, potentialType := range prefixes {
		if strings.HasPrefix(remote, prefix) {
			remote = strings.TrimPrefix(remote, prefix)
			protocolType = potentialType
		}
	}

	if protocolType == UnknownProtocol && sshutil.IsImplicitSSHTransport(remote) {
		protocolType = SSHProtocol
	}

	// remove name from ssh
	if protocolType == SSHProtocol {
		parts := strings.SplitN(remote, "@", 2)
		if len(parts) == 2 {
			remote = parts[1]
		}
	}

	return remote, protocolType
}

// GetCloneURL returns the repo to clone, and a path relative to the repo
//
//	"github.com/earthly/earthly"             ---> ("git@github.com/earthly/earthly.git", "")
//	"github.com/earthly/earthly/examples"    ---> ("git@github.com/earthly/earthly.git", "examples")
//	"github.com/earthly/earthly/examples/tutorial/go/part3" ---> ("git@github.com/earthly/earthly.git", "examples/go")
//
// Additionally a ssh keyscan might be returned (or an empty string indicating none was configured)
// Also, a custom "git ssh command" may be returned. This is part of this function since the user may
// specify a command necessary to clone their repository successfully.
func (gl *GitLookup) GetCloneURL(path string) (gitURL string, subPath string, keyScans []string, sshCommand string, err error) {
	gl.mu.Lock()
	defer gl.mu.Unlock()
	match, m, err := gl.getGitMatcherByPath(path)
	if err != nil {
		return "", "", nil, "", err
	}

	n := len(match)
	subPath = ""
	if len(path) > n {
		subPath = path[n+1:]
		path = path[:n]
	}
	host := path[:strings.IndexByte(path, '/')]
	gitPath := m.prefix + match[(strings.IndexByte(match, '/')+1):] + m.suffix

	sshCommand = m.sshCommand

	if m.sub != "" {
		if !m.re.MatchString(path) {
			return "", "", nil, "", errors.Errorf("failed to determine git path to clone for %q", path)
		}
		gitURL := m.re.ReplaceAllString(path, m.sub)
		gl.console.VerbosePrintf("converted earthly reference %s to git url %s (using regex substitution %s)", path, stringutil.ScrubCredentials(gitURL), stringutil.ScrubCredentials(m.sub))
		var keyScans []string
		remote, protocol := parseGitProtocol(gitURL)
		if protocol == SSHProtocol {
			subHost := remote[:strings.IndexByte(remote, '/')]
			_, keyScans, err = gl.getHostKeyAlgorithms(subHost)
			if err != nil {
				return "", "", nil, "", err
			}
			if len(keyScans) == 0 && m.strictHostKeyChecking {
				return "", "", nil, "", errors.Errorf("no known_hosts entries exist for substituted host %s", subHost)
			}
		}
		return gitURL, subPath, keyScans, sshCommand, nil
	}

	gitURL, keyScans, sshCommand, err = gl.makeCloneURL(m, host, gitPath)
	if err != nil {
		return "", "", nil, "", err
	}
	gl.console.VerbosePrintf("converted earthly reference %s to git url %s", path, stringutil.ScrubCredentials(gitURL))
	return gitURL, subPath, keyScans, sshCommand, nil
}

// ConvertCloneURL takes a url such as https://github.com/user/repo.git or git@github.com:user/repo.git
// and makes use of configured git credentials and protocol preferences to convert it into the appropriate
// https or ssh protocol.
// it also returns a keyScan and sshCommand
func (gl *GitLookup) ConvertCloneURL(inURL string) (gitURL string, keyScans []string, sshCommand string, err error) {
	var host string
	var gitPath string

	remote, protocol := parseGitProtocol(inURL)
	switch protocol {
	case HTTPProtocol, HTTPSProtocol:
		splits := strings.SplitN(remote, "/", 2)
		if len(splits) != 2 {
			return "", nil, "", errors.Errorf("failed to split path from host in %s", remote)
		}
		host = splits[0]
		gitPath = splits[1]
	case SSHProtocol:
		if sshutil.IsImplicitSSHTransport(inURL) {
			splits := strings.SplitN(remote, ":", 2)
			if len(splits) != 2 {
				return "", nil, "", errors.Errorf("failed to split path from host in %s", remote)
			}
			host = splits[0]
			gitPath = splits[1]
		} else {
			u, err := url.Parse(inURL)
			if err != nil {
				return "", nil, "", errors.Wrapf(err, "failed to parse %s", inURL)
			}
			if u.Scheme != "ssh" {
				panic(fmt.Sprintf("expected scheme of ssh; got %s", u.Scheme)) // shouldn't happen
			}
			host = strings.TrimSuffix(u.Host, ":22")
			gitPath = u.Path
		}
	default:
		return "", nil, "", errors.Errorf("unsupported git protocol %v", protocol)
	}

	m := gl.getGitMatcherByName(host)
	if m.sub != "" {
		path := host + strings.TrimSuffix(gitPath, ".git")
		if !m.re.MatchString(path) {
			return "", nil, "", errors.Errorf("failed to determine git path to clone for %q", path)
		}
		gitURL = m.re.ReplaceAllString(path, m.sub)
		remote, protocol := parseGitProtocol(gitURL)
		if protocol == SSHProtocol {
			subHost := remote[:strings.IndexByte(remote, '/')]
			_, keyScans, err = gl.getHostKeyAlgorithms(subHost)
			if err != nil {
				return "", nil, "", err
			}
			if len(keyScans) == 0 && m.strictHostKeyChecking {
				return "", nil, "", errors.Errorf("no known_hosts entries exist for substituted host %s", subHost)
			}
		}
		return gitURL, keyScans, m.sshCommand, nil
	}

	return gl.makeCloneURL(m, host,
		m.prefix+gitPath, // Note that inURL already contains the suffix
	)
}

func loadKnownHostsFromPath(path string) ([]string, error) {
	knownHostsExists, err := fileutil.FileExists(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if %s exists", path)
	}
	if !knownHostsExists {
		return nil, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %s", path)
	}
	return knownHostsToKeyScans(string(b)), nil
}

func loadKnownHosts() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user home dir")
	}
	knownHosts, err := loadKnownHostsFromPath(filepath.Join(homeDir, ".ssh/known_hosts"))
	if err != nil {
		return nil, err
	}
	etcKnownHosts, err := loadKnownHostsFromPath("/etc/ssh/ssh_known_hosts")
	if err != nil {
		return nil, err
	}
	knownHosts = append(knownHosts, etcKnownHosts...)
	return knownHosts, nil
}
