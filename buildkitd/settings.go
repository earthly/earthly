package buildkitd

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// Settings represents the buildkitd settings used to start up the daemon with.
type Settings struct {
	SSHAuthSock     string   `json:"sshAuthSock"`
	GitURLInsteadOf string   `json:"gitUrlInsteadOf"`
	CacheSizeMb     int      `json:"cacheSizeMb"`
	GitConfig       string   `json:"gitConfig"`
	GitCredentials  []string `json:"gitCredentials"`
	RunDir          string   `json:"runDir"`
	Debug           bool     `json:"debug"`
	DebuggerPort    int      `json:"debuggerPort"`
}

// Hash returns a secure hash of the settings.
func (s Settings) Hash() (string, error) {
	dt, err := json.Marshal(s)
	if err != nil {
		return "", errors.Wrap(err, "json marshal settings")
	}
	// Extra sha256 wrap is needed due to bcrypt password length limit.
	dtSha256 := sha256.Sum256(dt)
	hash, err := bcrypt.GenerateFromPassword(dtSha256[:], bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "generate from password")
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

// VerifyHash checks whether a given hash matches the settings.
func (s Settings) VerifyHash(hash string) (bool, error) {
	hashBytes, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, errors.Wrap(err, "base64 decode hash")
	}
	dt, err := json.Marshal(s)
	if err != nil {
		return false, errors.Wrap(err, "json marshal settings")
	}
	dtSha256 := sha256.Sum256(dt)
	err = bcrypt.CompareHashAndPassword(hashBytes, dtSha256[:])
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, errors.Wrap(err, "compare hash and password")
	}
	return true, nil
}
