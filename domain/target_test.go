package domain

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	target1 := Target{
		GitURL:  "github.com/earthly/earthly",
		GitPath: "examples",
	}
	target2 := Target{
		LocalPath: "./go",
		Target:    "+docker",
	}

	target3, err := JoinTargets(target1, target2)
	Equal(t, nil, err)
	Equal(t, "github.com/earthly/earthly", target3.GitURL)
	Equal(t, "examples/go", target3.GitPath)
	Equal(t, "+docker", target3.Target)
}
