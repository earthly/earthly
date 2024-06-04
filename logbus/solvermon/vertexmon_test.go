package solvermon

import (
	"math"
	"testing"

	"github.com/earthly/cloud-api/logstream"
)

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name         string
		errString    string
		expectedCode uint64
	}{
		{
			name:         "no match",
			errString:    "random error message",
			expectedCode: 0,
		},
		{
			name:         "match with exit code",
			errString:    "process \"foo\" did not complete successfully: exit code: 123",
			expectedCode: 123,
		},
		{
			name:         "match with max uint32",
			errString:    "process \"foo\" did not complete successfully: exit code: 4294967295",
			expectedCode: math.MaxUint32,
		},
		{
			name:         "match with max uint32",
			errString:    "some wrap message: process \"foo\" did not complete successfully: exit code: 8",
			expectedCode: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := getExitCode(tt.errString)
			if code != tt.expectedCode {
				t.Errorf("getExitCode(%q) = %d, want %d", tt.errString, code, tt.expectedCode)
			}
		})
	}
}

func TestDetermineFatalErrorType(t *testing.T) {
	tests := []struct {
		name          string
		errString     string
		exitCode      uint64
		expectedType  logstream.FailureType
		expectedFatal bool
	}{
		{
			name:          "context canceled",
			errString:     "context canceled",
			exitCode:      0,
			expectedType:  logstream.FailureType_FAILURE_TYPE_UNKNOWN,
			expectedFatal: false,
		},
		{
			name:          "exit code 123",
			errString:     "process \"foo\" did not complete successfully: exit code: 123",
			exitCode:      123,
			expectedType:  logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT,
			expectedFatal: true,
		},
		{
			name:          "exit code max uint32",
			errString:     "process \"foo\" did not complete successfully: exit code: 4294967295",
			exitCode:      math.MaxUint32,
			expectedType:  logstream.FailureType_FAILURE_TYPE_OOM_KILLED,
			expectedFatal: true,
		},
		{
			name:          "file not found",
			errString:     "failed to calculate checksum of ref foo: bar",
			exitCode:      0,
			expectedType:  logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND,
			expectedFatal: true,
		},
		{
			name:          "git error",
			errString:     "EARTHLY_GIT_STDERR: Z2l0IC1jI...",
			exitCode:      0,
			expectedType:  logstream.FailureType_FAILURE_TYPE_GIT,
			expectedFatal: true,
		},
		{
			name:          "unknown error",
			errString:     "unknown error",
			exitCode:      0,
			expectedType:  logstream.FailureType_FAILURE_TYPE_UNKNOWN,
			expectedFatal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fatalType, fatal := determineFatalErrorType(tt.errString, tt.exitCode)
			if fatalType != tt.expectedType {
				t.Errorf("determineFatalErrorType(%q, %d) = %v, want %v", tt.errString, tt.exitCode, fatalType, tt.expectedType)
			}
			if fatal != tt.expectedFatal {
				t.Errorf("determineFatalErrorType(%q, %d) = %v, want %v", tt.errString, tt.exitCode, fatal, tt.expectedFatal)
			}
		})
	}
}
