package builder

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestParseVertexName(t *testing.T) {

	for _, tt := range []struct {
		//input
		name string

		//expect output
		targetStr      string
		targetBrackets string
		salt           string
		operation      string
	}{
		{
			name:           "docker-image://docker.io/alpine/git:v2.24.1",
			targetStr:      "internal",
			targetBrackets: "",
			salt:           "internal",
			operation:      "docker-image://docker.io/alpine/git:v2.24.1",
		},
		{
			name:           "[internal] GET GIT META github.com/earthly/buildkit:earthly-main",
			targetStr:      "internal",
			targetBrackets: "",
			salt:           "internal",
			operation:      "GET GIT META github.com/earthly/buildkit:earthly-main",
		},
		{
			name:           "[./earthfile2llb/parser+base(cGxhdGZvcm09bGludXgvYW1kNjQ=) 7504504064263669287]",
			targetStr:      "internal",
			targetBrackets: "",
			salt:           "internal",
			operation:      "[./earthfile2llb/parser+base(cGxhdGZvcm09bGludXgvYW1kNjQ=) 7504504064263669287]",
		},
		{
			name:           "[internal] load metadata for docker.io/tonistiigi/xx:golang@sha256:6f7d999551dd471b58f70716754290495690efa8421e0a1fcf18eb11d0c0a537",
			targetStr:      "internal",
			targetBrackets: "",
			salt:           "internal",
			operation:      "load metadata for docker.io/tonistiigi/xx:golang@sha256:6f7d999551dd471b58f70716754290495690efa8421e0a1fcf18eb11d0c0a537",
		},
		{
			name:           "[busybox:1.32.1 2423175906] Load metadata linux/amd64",
			targetStr:      "busybox:1.32.1",
			targetBrackets: "",
			salt:           "2423175906",
			operation:      "Load metadata linux/amd64",
		},
		{
			name:           "[./examples/tests/local+test-local *local* 5577006791947779410] RUN whoami",
			targetStr:      "./examples/tests/local+test-local *local*",
			targetBrackets: "",
			salt:           "5577006791947779410",
			operation:      "RUN whoami",
		},
	} {
		targetStr, targetBrackets, salt, operation := parseVertexName(tt.name)
		Equal(t, tt.targetStr, targetStr)
		Equal(t, tt.targetBrackets, targetBrackets)
		Equal(t, tt.salt, salt)
		Equal(t, tt.operation, operation)

	}
}
