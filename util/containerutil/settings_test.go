package containerutil

import (
	"context"
	"strings"
	"testing"

	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/stretchr/testify/assert"
)

var noopArgs = parsedCLIVals{}

type results struct {
	buildkit      string
	localRegistry string
}

type parsedCLIVals struct {
	buildkit string
}

func TestBuildArgMatrix(t *testing.T) {
	var tests = []struct {
		testName string
		config   config.GlobalConfig
		args     parsedCLIVals
		expected results
	}{
		{
			"No Config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "",
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "docker-container://test-buildkitd",
				localRegistry: "",
			},
		},
		{
			"Remote Local in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "tcp://127.0.0.1:8372",
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "tcp://127.0.0.1:8372",
				localRegistry: "",
			},
		},
		{
			"Remote remote in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "tcp://my-cool-host:8372",
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "tcp://my-cool-host:8372",
				localRegistry: "",
			},
		},
		{
			"Nonstandard local in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://my-container",
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "docker-container://my-container",
				localRegistry: "",
			},
		},
		{
			"Remote Local in config, no CLI, validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "tcp://127.0.0.1:8372",
				LocalRegistryHost: "tcp://127.0.0.1:8371",
			},
			noopArgs,
			results{
				buildkit:      "tcp://127.0.0.1:8372",
				localRegistry: "tcp://127.0.0.1:8371",
			},
		},
		{
			"Remote remote in config, no CLI, skip validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "tcp://my-cool-host:8372",
				LocalRegistryHost: "this-is-not-a-url",
			},
			noopArgs,
			results{
				buildkit:      "tcp://my-cool-host:8372",
				localRegistry: "",
			},
		},
		{
			"Local in config, no CLI, validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://my-cool-container",
				LocalRegistryHost: "tcp://127.0.0.1:8371",
			},
			noopArgs,
			results{
				buildkit:      "docker-container://my-cool-container",
				localRegistry: "tcp://127.0.0.1:8371",
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info, false)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{
			InstallationName: "test-stub",
		})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		urls, err := stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostCLIValue:       tt.args.buildkit,
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			InstallationName:           "test",
			DefaultPort:                8372,
			Console:                    logger,
		})
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, results{
			buildkit:      urls.BuildkitHost.String(),
			localRegistry: urls.LocalRegistryHost.String(),
		})
	}
}

func TestBuildArgMatrixValidationFailures(t *testing.T) {
	var tests = []struct {
		testName string
		config   config.GlobalConfig
		expected error
		log      string
	}{
		{
			"Invalid buildkit URL",
			config.GlobalConfig{
				BuildkitHost:      "http\r://foo.com/",
				LocalRegistryHost: "",
			},
			errURLParseFailure,
			"",
		},
		{
			"Invalid registry URL",
			config.GlobalConfig{
				BuildkitHost:      "",
				LocalRegistryHost: "http\r://foo.com/",
			},
			errURLParseFailure,
			"",
		},
		{
			"Homebrew test",
			config.GlobalConfig{
				BuildkitHost:      "127.0.0.1",
				LocalRegistryHost: "",
			},
			errURLValidationFailure,
			"",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info, false)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{
			InstallationName: "test-stub",
		})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		_, err = stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			Console:                    logger,
			InstallationName:           "test",
			DefaultPort:                8372,
		})
		assert.ErrorIs(t, err, tt.expected)
		assert.Contains(t, logs.String(), tt.log)
	}
}

func TestParseAndValidateURLFailures(t *testing.T) {
	var tests = []struct {
		testName string
		url      string
		expected error
	}{
		{
			"Invalid URL",
			"http\r://foo.com/",
			errURLParseFailure,
		},
		{
			"Invalid Scheme",
			"gopher://my-hole",
			errURLValidationFailure,
		},
		{
			"Missing Port",
			"tcp://my-server",
			errURLValidationFailure,
		},
	}

	for _, tt := range tests {
		_, err := parseAndValidateURL(tt.url)
		assert.ErrorIs(t, err, tt.expected)
	}
}

func TestParseAndValidateURL(t *testing.T) {
	var tests = []struct {
		testName string
		url      string
	}{
		{
			"docker-container URL",
			"docker-container://my-container",
		},
		{
			"tcp URL",
			"tcp://my-host:42",
		},
	}

	for _, tt := range tests {
		_, err := parseAndValidateURL(tt.url)
		assert.NoError(t, err)
	}
}

func TestBuildArgMatrixValidationNonIssues(t *testing.T) {
	var tests = []struct {
		testName string
		config   config.GlobalConfig
		log      string
	}{
		{
			"Buildkit/Local Registry host mismatch, schemes differ",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://127.0.0.1:8372",
				LocalRegistryHost: "tcp://localhost:8371",
			},
			"Buildkit and Local Registry URLs are pointed at different hosts",
		},
		{
			"Buildkit/Debugger host mismatch, schemes differ",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://bk:1234",
				LocalRegistryHost: "",
			},
			"Buildkit and Debugger URLs are pointed at different hosts",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info, false)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{
			InstallationName: "test-stub",
		})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		_, err = stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			Console:                    logger,
			InstallationName:           "test",
			DefaultPort:                8372,
		})
		assert.NoError(t, err)
		assert.NotContains(t, logs.String(), tt.log)
	}
}
