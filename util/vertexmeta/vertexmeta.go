package vertexmeta

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"regexp"
	"sort"
	"strings"

	"github.com/earthly/earthly/ast/spec"
)

// VertexMeta is metadata associated with the vertex. This is passed from the
// converter to the solver monitor via BuildKit.
type VertexMeta struct {
	SourceLocation      *spec.SourceLocation `json:"sl,omitempty"`
	RepoGitURL          string               `json:"rgu,omitempty"`
	RepoGitHash         string               `json:"rgh,omitempty"`
	RepoFileRelToRepo   string               `json:"rfr,omitempty"`
	CommandID           string               `json:"cid,omitempty"`
	TargetID            string               `json:"tid,omitempty"`
	TargetName          string               `json:"tnm,omitempty"`
	CanonicalTargetName string               `json:"ctnm,omitempty"`
	Platform            string               `json:"plt,omitempty"`
	NonDefaultPlatform  bool                 `json:"defplt,omitempty"`
	Local               bool                 `json:"lcl,omitempty"`
	Interactive         bool                 `json:"itrctv,omitempty"`
	OverridingArgs      map[string]string    `json:"args,omitempty"`
	Secrets             []string             `json:"secrets,omitempty"`
	Internal            bool                 `json:"itrnl,omitempty"`
	Runner              string               `json:"runner,omitempty"`
}

var vertexRegexp = regexp.MustCompile(`(?s)^\[([^\]]*)\] (.*)$`)

// ParseFromVertexPrefix parses the vertex prefix from the given string.
func ParseFromVertexPrefix(in string) (*VertexMeta, string) {
	vm := &VertexMeta{}
	tail := in
	if strings.HasPrefix(in, "importing cache manifest") ||
		strings.HasPrefix(in, "exporting cache") {
		vm.TargetName = "cache"
		return vm, tail
	}
	match := vertexRegexp.FindStringSubmatch(in)
	if len(match) < 2 {
		vm.TargetName = "internal"
		vm.Internal = true
		return vm, tail
	}
	vmDt64 := match[1]
	tail = match[2]
	dt, err := base64.StdEncoding.DecodeString(vmDt64)
	if err != nil {
		// Either "context <context-name>"
		// or "internal"
		// or coming from Dockerfile: "<target> <step>/<total-steps>".
		splits := strings.SplitN(vmDt64, " ", 2)
		if len(splits) > 0 {
			vm.TargetName = splits[0]
		}
		if vm.TargetName == "internal" {
			vm.Internal = true
		}
		return vm, tail
	}
	err = json.Unmarshal(dt, vm)
	if err != nil {
		vm.TargetName = vmDt64
		return vm, tail
	}
	return vm, tail
}

// ToVertexPrefix returns the vertex prefix for the given VertexMeta.
func (vm *VertexMeta) ToVertexPrefix() string {
	dt, err := json.Marshal(vm)
	if err != nil {
		panic(err) // should never happen
	}
	b64Str := base64.StdEncoding.EncodeToString(dt)
	return fmt.Sprintf("[%s] ", b64Str)
}

// OverridingArgsString returns the string representation of the overriding
// arguments.
func (vm *VertexMeta) OverridingArgsString() string {
	if vm.OverridingArgs == nil {
		return ""
	}
	var args []string
	for k, v := range vm.OverridingArgs {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(args)
	return strings.Join(args, " ")
}

// SecretsString returns the string representation of the secrets. It will only
// include the names of secrets, no values.
func (vm *VertexMeta) SecretsString() string {
	if len(vm.Secrets) == 0 {
		return ""
	}
	return "secrets: " + strings.Join(vm.Secrets, " ")
}

// Salt returns a string identifying the target as uniquely as possible.
func (vm *VertexMeta) Salt() string {
	if vm.TargetID != "" {
		return vm.TargetID
	}
	var name string
	switch {
	case vm.TargetName != "":
		name = vm.TargetName
	case vm.Internal:
		name = "internal"
	default:
		name = "unknown"
	}

	overridingArgsString := vm.OverridingArgsString()
	if vm.Platform == "" && overridingArgsString == "" {
		return name // don't hash when unmarshalling VertexMeta failed
	}

	h := fnv.New32a()
	h.Write([]byte(vm.Platform))
	h.Write([]byte(overridingArgsString))
	return fmt.Sprintf("%s-%d", name, h.Sum32())
}
