//go:build offline

package cloud_test

import (
	"testing"
	"time"

	"github.com/earthly/earthly/cloud"
	"github.com/poy/onpar/expect"
)

// TestClientConnectionIsLazy tests that creating a NewCLient doesn't perform any network I/O
func TestClientConnectionIsLazy(t *testing.T) {
	expect := expect.New(t)
	c, err := cloud.NewClient("https://this", "shouldnt.matter:443", false, "", "", "", "", "", nil, time.Second)
	expect(err).To(not(haveOccurred()))
}
