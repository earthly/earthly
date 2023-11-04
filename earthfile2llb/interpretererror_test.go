package earthfile2llb

import (
	"github.com/earthly/earthly/ast/spec"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromError(t *testing.T) {
	ieWithStack := Errorf(&spec.SourceLocation{
		File:        "path/To/Earthfile",
		StartLine:   90,
		StartColumn: 8,
	}, "", "some stack", "some error message")

	ieWithoutStack := Errorf(&spec.SourceLocation{
		File:        "path/To/Earthfile",
		StartLine:   90,
		StartColumn: 8,
	}, "", "", "some error message")

	tests := map[string]struct {
		providerErr    error
		expectedResult *InterpreterError
		success        bool
	}{
		"nil error": {},
		"no file path": {
			providerErr: errors.New("line 5:4 some error message"),
		},
		"no line": {
			providerErr: errors.New("path/to/Earthfile 5:4 some error message"),
		},
		"no column": {
			providerErr: errors.New("path/to/Earthfile line 5:"),
		},
		"no error message": {
			providerErr: errors.New("path/to/Earthfile line 5:4"),
		},
		"success without stack": {
			providerErr:    ieWithStack,
			expectedResult: ieWithStack,
			success:        true,
		},
		"success with stack": {
			providerErr:    ieWithoutStack,
			expectedResult: ieWithoutStack,
			success:        true,
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ie, ok := FromError(tc.providerErr)
			assert.Equal(t, tc.expectedResult, ie)
			assert.Equal(t, tc.success, ok)
		})
	}
}
