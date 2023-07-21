package envutil

import (
	"os"
)

// IsTrue returns true if the env variable `k` is set
// to something bash would interpret as true.
func IsTrue(k string) bool {
	switch os.Getenv(k) {
	case "", "0", "false", "FALSE":
		return false
	default:
		return true
	}
}
