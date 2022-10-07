package earthfile2llb

import (
	"testing"

	"github.com/earthly/earthly/features"
)

func Test_parseSecretFlag(t *testing.T) {

	tests := []struct {
		name              string
		val               string
		wantSecretID      string
		wantEnvVar        string
		wantErr           bool
		useProjectSecrets bool
	}{
		{
			name: "empty value",
		},
		{
			name:         "just the name",
			val:          "SECRET_ID",
			wantSecretID: "SECRET_ID",
			wantEnvVar:   "SECRET_ID",
		},
		{
			name:    "blank secret name",
			val:     "=BAR",
			wantErr: true,
		},
		{
			name:    "blank secret value",
			val:     "FOO=",
			wantErr: false,
		},
		{
			name:              "has flag but includes +secrets/",
			val:               "FOO=+secrets/BAR",
			wantErr:           false,
			wantSecretID:      "BAR",
			wantEnvVar:        "FOO",
			useProjectSecrets: true,
		},
		{
			name:              "has flag no prefix",
			val:               "FOO=BAR",
			wantErr:           false,
			wantSecretID:      "BAR",
			wantEnvVar:        "FOO",
			useProjectSecrets: true,
		},
		{
			name:         "no flag, has prefix",
			val:          "FOO=+secrets/BAR",
			wantErr:      false,
			wantSecretID: "BAR",
			wantEnvVar:   "FOO",
		},
		{
			name:    "no flag, no prefix",
			val:     "FOO=BAR",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := &Converter{ftrs: &features.Features{UseProjectSecrets: test.useProjectSecrets}}

			secretID, envVar, err := c.parseSecretFlag(test.val)

			if test.wantErr && err == nil {
				t.Error("expected error, got nil")
			} else if !test.wantErr && err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if test.wantSecretID != secretID {
				t.Errorf("expected secret ID %q, got %q", test.wantSecretID, secretID)
			}

			if test.wantEnvVar != envVar {
				t.Errorf("expected env var %q, got %q", test.wantEnvVar, envVar)
			}
		})
	}
}
