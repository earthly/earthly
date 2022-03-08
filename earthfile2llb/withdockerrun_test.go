package earthfile2llb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"testing"
)

//goland:noinspection HttpUrlsUsage
func Test_processRegistryValue(t *testing.T) {
	type args struct {
		registry        string
		defaultProtocol string
	}
	tests := []struct {
		name                string
		args                args
		wantedString        string
		wantedProcessedList []string
	}{
		// Happy path
		{"Unset/empty", args{registry: "", defaultProtocol: "http"},
			"", []string{}},
		{"Single unqualified", args{registry: "example.com:999", defaultProtocol: "http"},
			"\\\"http://example.com:999\\\"", []string{"http://example.com:999"}},
		{"Single qualified http with http default untouched", args{registry: "http://example.com:999", defaultProtocol: "http"},
			"\\\"http://example.com:999\\\"", []string{}},
		{"Single qualified https with https default untouched", args{registry: "https://example.com:999", defaultProtocol: "https"},
			"\\\"https://example.com:999\\\"", []string{}},

		// Happy, but uncommon
		{"Multiple unqualified", args{registry: "example.com:998,example.com:999", defaultProtocol: "http"},
			"\\\"http://example.com:998\\\",\\\"http://example.com:999\\\"", []string{"http://example.com:998", "http://example.com:999"}},
		{"Multiple qualified", args{registry: "https://example.com:998,https://example.com:999", defaultProtocol: "https"},
			"\\\"https://example.com:998\\\",\\\"https://example.com:999\\\"", []string{}},
		{"Multiple with whitespace", args{registry: " example.com:997, example.com:998 ,example.com:999 ", defaultProtocol: "https"},
			"\\\"https://example.com:997\\\",\\\"https://example.com:998\\\",\\\"https://example.com:999\\\"", []string{"https://example.com:997", "https://example.com:998", "https://example.com:999"}},

		// Uncommon
		{"Single qualified http with https default untouched", args{registry: "http://example.com:999", defaultProtocol: "https"},
			"\\\"http://example.com:999\\\"", []string{}},
		{"Single qualified https with http default untouched", args{registry: "https://example.com:999", defaultProtocol: "http"},
			"\\\"https://example.com:999\\\"", []string{}},
		{"Protocol insignificant", args{registry: "example.com:999", defaultProtocol: "gopher"},
			"\\\"gopher://example.com:999\\\"", []string{"gopher://example.com:999"}},
		{"Invalid values processed as valid", args{registry: "junk 2 the hundredth power!", defaultProtocol: "https"},
			"\\\"https://junk 2 the hundredth power!\\\"", []string{"https://junk 2 the hundredth power!"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registryStr, processed := processRegistryValue(tt.args.registry, tt.args.defaultProtocol)
			assert.Equalf(t, tt.wantedString, registryStr, "processRegistryValue(%v, %v)", tt.args.registry, tt.args.defaultProtocol)
			assert.Equalf(t, tt.wantedProcessedList, processed, "processRegistryValue(%v, %v)", tt.args.registry, tt.args.defaultProtocol)
		})
	}
	t.Run("Sanity check of intended use", func(t *testing.T) {
		result, _ := processRegistryValue("example.com,127.0.0.1", "https")
		// Mimic the code passing the value to a shelled process and its usage in a "templated" JSON heredoc
		cmd := exec.Command("sh", "-c", "echo -n \"[${MIRROR}]\"")
		cmd.Env = []string{fmt.Sprintf("MIRROR=\"%s\"", result)}
		var out bytes.Buffer
		cmd.Stdout = &out
		assert.NoError(t, cmd.Run())
		jsonString := out.String()

		// Validate the result is valid JSON, compliments of https://stackoverflow.com/a/36922225/15001552
		var js json.RawMessage
		assert.NoError(t, json.Unmarshal([]byte(jsonString), &js))
	})
}
