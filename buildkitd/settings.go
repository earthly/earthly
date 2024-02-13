package buildkitd

import (
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/pkg/errors"
)

// Settings represents the buildkitd settings used to start up the daemon with.
type Settings struct {
	CacheSizeMb          int
	CacheSizePct         int
	CacheKeepDuration    int
	Debug                bool
	BuildkitAddress      string
	LocalRegistryAddress string
	AdditionalArgs       []string
	AdditionalConfig     string
	CniMtu               uint16
	Timeout              time.Duration `hash:"ignore"`
	TLSCA                string
	ClientTLSCert        string
	ClientTLSKey         string
	ServerTLSCert        string
	ServerTLSKey         string
	UseTCP               bool
	UseTLS               bool
	VolumeName           string
	IPTables             string
	MaxParallelism       int
	SatelliteName        string `hash:"ignore"`
	SatelliteDisplayName string `hash:"ignore"`
	SatelliteOrgID       string `hash:"ignore"`
	SatelliteToken       string `hash:"ignore"`
	SatelliteIsManaged   bool   `hash:"ignore"`
	EnableProfiler       bool
	NoUpdate             bool   `hash:"ignore"`
	StartUpLockPath      string `hash:"ignore"`
}

// Hash returns a secure hash of the settings.
func (s Settings) Hash() (string, error) {
	hash, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	if err != nil {
		return "", errors.Wrap(err, "hash settings")
	}

	return strconv.FormatUint(hash, 16), nil
}

// VerifyHash checks whether a given hash matches the settings.
func (s Settings) VerifyHash(hash string) (bool, error) {
	newHash, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	if err != nil {
		return false, errors.Wrap(err, "hash settings")
	}

	oldHash, err := strconv.ParseUint(strings.TrimSpace(hash), 16, 64)
	if err != nil {
		return false, errors.Wrap(err, "parse hash")
	}

	return oldHash == newHash, nil
}

// HasConfiguredCacheSize returns if the buildkitd cache size was configured
func (s Settings) HasConfiguredCacheSize() bool {
	return s.CacheSizeMb > 0 || s.CacheSizePct > 0
}
