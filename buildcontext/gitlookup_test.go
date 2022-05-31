package buildcontext

import (
	"fmt"
	"os"
	"testing"

	"github.com/earthly/earthly/conslogging"
	. "github.com/stretchr/testify/assert"
)

func Test_GetCloneURL(t *testing.T) {
	testcases := []struct {
		path    string
		url     string
		subPath string
		ok      bool
	}{
		{
			path:    "git.example.com/proj/repo",
			url:     "ssh://git.example.com:7777/proj/repo.git",
			subPath: "",
			ok:      true,
		},
		{
			path:    "git.example.com/proj/repo/inner/location",
			url:     "ssh://git.example.com:7777/proj/repo.git",
			subPath: "inner/location",
			ok:      true,
		},
	}

	logger := conslogging.Current(conslogging.NoColor, 0, conslogging.Info)
	gl := NewGitLookup(logger, "")
	err := gl.AddMatcher("git.example.com", "git.example.com/([^/]+)/([^/]+)", "ssh://git.example.com:7777/$1/$2.git",
		"", "", ".git", "ssh", "", false)
	Nil(t, err)

	for i, testcase := range testcases {
		t.Run(fmt.Sprintf("path test %d", i), func(t *testing.T) {
			url, subPath, _, err := gl.GetCloneURL(testcase.path)
			ok := err == nil
			Equal(t, ok, testcase.ok)
			Equal(t, url, testcase.url)
			Equal(t, subPath, testcase.subPath)
		})
	}
}

func Test_ConvertCloneURL(t *testing.T) {
	matcherName := "git.example.com"
	matcherPattern := "git.example.com/([^/]+)/([^/]+)"
	matcherSub := "ssh://git.example.com:7777/$1/$2.git"
	matcherSuffix := ".git"
	matcherProtocol := "ssh"

	err := os.Setenv("USER", "somebody")
	if err != nil {
		Error(t, err)
	}
	testcases := []struct {
		inURL   string
		matcher func(lookup *GitLookup) error
		outURL  string
		ok      bool
	}{
		{
			inURL: "ssh://git.example.com:7777/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"git", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://git@git.example.com:7777/proj/repo.git",
			ok:     true,
		},
		{
			inURL: "ssh://git.example.com:22/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"git", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://git@git.example.com:22/proj/repo.git",
			ok:     true,
		},
		{
			inURL: "ssh://git.example.com/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"git", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://git@git.example.com/proj/repo.git",
			ok:     true,
		},
		{
			inURL: "ssh://git.example.com:7777/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://somebody@git.example.com:7777/proj/repo.git",
			ok:     true,
		},
		{
			inURL: "ssh://git.example.com:22/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://somebody@git.example.com:22/proj/repo.git",
			ok:     true,
		},
		{
			inURL: "ssh://git.example.com/proj/repo.git",
			matcher: func(gl *GitLookup) error {
				return gl.AddMatcher(matcherName, matcherPattern, matcherSub,
					"", "", matcherSuffix, matcherProtocol, "", false)
			},
			outURL: "ssh://somebody@git.example.com/proj/repo.git",
			ok:     true,
		},
	}

	logger := conslogging.Current(conslogging.NoColor, 0, conslogging.Info)
	gl := NewGitLookup(logger, "")

	for i, testcase := range testcases {
		t.Run(fmt.Sprintf("inURL test %d", i), func(t *testing.T) {

			err := testcase.matcher(gl)
			Nil(t, err)
			url, _, err := gl.ConvertCloneURL(testcase.inURL)
			ok := err == nil
			Equal(t, testcase.ok, ok)
			Equal(t, testcase.outURL, url)
		})
	}
}

func Test_parseKeyScanIfHostMatches(t *testing.T) {
	testcases := []struct {
		key      string
		hostname string
		ok       bool
		keyAlg   string
		keyData  string
	}{
		{
			key:      "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ssh-rsa",
			keyData:  "AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			key:      "github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ecdsa-sha2-nistp256",
			keyData:  "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
		},
		{
			key:      "github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ssh-ed25519",
			keyData:  "AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
		},
		// test that port will be ignored
		{
			key:      "github.com ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
			hostname: "github.com:1234",
			ok:       true,
			keyAlg:   "ssh-rsa",
			keyData:  "AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			key:      "github.com ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
			hostname: "github.com:1234",
			ok:       true,
			keyAlg:   "ecdsa-sha2-nistp256",
			keyData:  "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
		},
		{
			key:      "github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
			hostname: "github.com:1234",
			ok:       true,
			keyAlg:   "ssh-ed25519",
			keyData:  "AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
		},
		// test hashing
		{
			key:      "|1|bDbXcpQoMMPAtFztKvPqwjqYTNw=|I4eHPek2rZ+DQKtZN5VOad+Zccg= ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ssh-rsa",
			keyData:  "AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			key:      "|1|YOJyk3Bq4EmEH2MStLodIEcrwmU=|auDKvWzTkfBNRBBuPuQt8JhrU1w= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ecdsa-sha2-nistp256",
			keyData:  "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
		},
		{
			key:      "|1|EQsTUtPg//BcsyEO7v11tCiTxvs=|UvtKlPuh0OJzNWECXxhUDkxEnFM= ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
			hostname: "github.com",
			ok:       true,
			keyAlg:   "ssh-ed25519",
			keyData:  "AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
		},
		// test hashing and ignored port
		{
			key:      "|1|bDbXcpQoMMPAtFztKvPqwjqYTNw=|I4eHPek2rZ+DQKtZN5VOad+Zccg= ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
			hostname: "github.com:2222",
			ok:       true,
			keyAlg:   "ssh-rsa",
			keyData:  "AAAAB3NzaC1yc2EAAAABIwAAAQEAq2A7hRGmdnm9tUDbO9IDSwBK6TbQa+PXYPCPy6rbTrTtw7PHkccKrpp0yVhp5HdEIcKr6pLlVDBfOLX9QUsyCOV0wzfjIJNlGEYsdlLJizHhbn2mUjvSAHQqZETYP81eFzLQNnPHt4EVVUh7VfDESU84KezmD5QlWpXLmvU31/yMf+Se8xhHTvKSCZIFImWwoG6mbUoWf9nzpIoaSjB+weqqUUmpaaasXVal72J+UX2B+2RPW3RcT0eOzQgqlJL3RKrTJvdsjE3JEAvGq3lGHSZXy28G3skua2SmVi/w4yCE6gbODqnTWlg7+wC604ydGXA8VJiS5ap43JXiUFFAaQ==",
		},
		{
			key:      "|1|YOJyk3Bq4EmEH2MStLodIEcrwmU=|auDKvWzTkfBNRBBuPuQt8JhrU1w= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
			hostname: "github.com:2222",
			ok:       true,
			keyAlg:   "ecdsa-sha2-nistp256",
			keyData:  "AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBEmKSENjQEezOmxkZMy7opKgwFB9nkt5YRrYMjNuG5N87uRgg6CLrbo5wAdT/y6v0mKV0U2w0WZ2YB/++Tpockg=",
		},
		{
			key:      "|1|EQsTUtPg//BcsyEO7v11tCiTxvs=|UvtKlPuh0OJzNWECXxhUDkxEnFM= ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
			hostname: "github.com:2222",
			ok:       true,
			keyAlg:   "ssh-ed25519",
			keyData:  "AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl",
		},
	}

	for i, testcase := range testcases {
		t.Run(fmt.Sprintf("key test %d", i), func(t *testing.T) {
			keyAlg, keyData, err := parseKeyScanIfHostMatches(testcase.key, testcase.hostname)
			ok := err == nil
			Equal(t, ok, testcase.ok)
			Equal(t, keyAlg, testcase.keyAlg)
			Equal(t, keyData, testcase.keyData)

			// test other hostnames don't work
			_, _, err = parseKeyScanIfHostMatches(testcase.key, "this-hostname-should-not-exist.org")
			Error(t, err)
			_, _, err = parseKeyScanIfHostMatches(testcase.key, "this-hostname-should-not-exist.org:2222")
			Error(t, err)
		})
	}
}
