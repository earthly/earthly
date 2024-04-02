package cloud_test

import (
	"testing"
	"time"

	"github.com/earthly/earthly/cloud"
	"github.com/poy/onpar"
	"github.com/poy/onpar/expect"
)

// TestOffline tests that creating a NewCLient doesn't perform any network I/O
func TestOffline(t *testing.T) {
	o := onpar.New(t)
	defer o.Run()
	o.Spec("ClientConnectionIsLazy", func(t *testing.T) {
		expect := expect.New(t)
		_, err := cloud.NewClient("https://this", "shouldnt.matter:443", false, "", "", "", "", "", nil, nil, time.Second)
		expect(err).To(not(haveOccurred()))
	})
}
