package containerutil

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/stretchr/testify/assert"
)

var noopArgs = parsedCLIVals{}

type results struct {
	buildkit      string
	debugger      string
	localRegistry string
}

type parsedCLIVals struct {
	buildkit string
	debugger string
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
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      DockerAddress,
				debugger:      fmt.Sprintf("tcp://127.0.0.1:%v", config.DefaultDebuggerPort),
				localRegistry: "",
			},
		},
		{
			"Remote Local in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "tcp://127.0.0.1:8372",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "tcp://127.0.0.1:8372",
				debugger:      fmt.Sprintf("tcp://127.0.0.1:%v", config.DefaultDebuggerPort),
				localRegistry: "",
			},
		},
		{
			"Remote remote in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "tcp://my-cool-host:8372",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "tcp://my-cool-host:8372",
				debugger:      fmt.Sprintf("tcp://my-cool-host:%v", config.DefaultDebuggerPort),
				localRegistry: "",
			},
		},
		{
			"Nonstandard local in config, no CLI",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://my-container",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      "docker-container://my-container",
				debugger:      fmt.Sprintf("tcp://127.0.0.1:%v", config.DefaultDebuggerPort),
				localRegistry: "",
			},
		},
		{
			"Debugger port specified",
			config.GlobalConfig{
				BuildkitHost:      "",
				DebuggerHost:      "",
				DebuggerPort:      5678,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      DockerAddress,
				debugger:      "tcp://127.0.0.1:5678",
				localRegistry: "",
			},
		},
		{
			"Debugger host and port specified",
			config.GlobalConfig{
				BuildkitHost:      "",
				DebuggerHost:      "tcp://127.0.0.1:1234",
				DebuggerPort:      5678,
				LocalRegistryHost: "",
			},
			noopArgs,
			results{
				buildkit:      DockerAddress,
				debugger:      "tcp://127.0.0.1:1234",
				localRegistry: "",
			},
		},
		{
			"No Config, buildkit and debugger on CLI",
			config.GlobalConfig{
				BuildkitHost:      "tcp://not-this-bk:1234",
				DebuggerHost:      "tcp://not-this-db:1234",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			parsedCLIVals{
				buildkit: "tcp://ok-bk:42",
				debugger: "tcp://ok-db:43",
			},
			results{
				buildkit:      "tcp://ok-bk:42",
				debugger:      "tcp://ok-db:43",
				localRegistry: "",
			},
		},
		{
			"Remote Local in config, no CLI, validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "tcp://127.0.0.1:8372",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "tcp://127.0.0.1:8371",
			},
			noopArgs,
			results{
				buildkit:      "tcp://127.0.0.1:8372",
				debugger:      fmt.Sprintf("tcp://127.0.0.1:%v", config.DefaultDebuggerPort),
				localRegistry: "tcp://127.0.0.1:8371",
			},
		},
		{
			"Remote remote in config, no CLI, skip validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "tcp://my-cool-host:8372",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "this-is-not-a-url",
			},
			noopArgs,
			results{
				buildkit:      "tcp://my-cool-host:8372",
				debugger:      fmt.Sprintf("tcp://my-cool-host:%v", config.DefaultDebuggerPort),
				localRegistry: "",
			},
		},
		{
			"Local in config, no CLI, validate registry host",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://my-cool-container",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "tcp://127.0.0.1:8371",
			},
			noopArgs,
			results{
				buildkit:      "docker-container://my-cool-container",
				debugger:      fmt.Sprintf("tcp://127.0.0.1:%v", config.DefaultDebuggerPort),
				localRegistry: "tcp://127.0.0.1:8371",
			},
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		urls, err := stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostCLIValue:       tt.args.buildkit,
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			DebuggerHostCLIValue:       tt.args.debugger,
			DebuggerHostFileValue:      tt.config.DebuggerHost,
			DebuggerPortFileValue:      tt.config.DebuggerPort,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			Console:                    logger,
		})
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, results{
			buildkit:      urls.BuildkitHost.String(),
			debugger:      urls.DebuggerHost.String(),
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
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			errURLParseFailure,
			"",
		},
		{
			"Invalid debugger URL",
			config.GlobalConfig{
				BuildkitHost:      "",
				DebuggerHost:      "http\r://foo.com/",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			errURLParseFailure,
			"",
		},
		{
			"Invalid registry URL",
			config.GlobalConfig{
				BuildkitHost:      "",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "http\r://foo.com/",
			},
			errURLParseFailure,
			"",
		},
		{
			"Buildkit/Local Registry host mismatch",
			config.GlobalConfig{
				BuildkitHost:      "tcp://127.0.0.1:8372",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "tcp://localhost:8371",
			},
			nil,
			"Buildkit and local registry URLs are pointed at different hosts",
		},
		{
			"Buildkit/Debugger host mismatch",
			config.GlobalConfig{
				BuildkitHost:      "tcp://bk:1234",
				DebuggerHost:      "tcp://db:5678",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			nil,
			"Buildkit and debugger URLs are pointed at different hosts",
		},
		{
			"Buildkit/Debugger port clash",
			config.GlobalConfig{
				BuildkitHost:      "tcp://cool-host:1234",
				DebuggerHost:      "tcp://cool-host:1234",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			errURLValidationFailure,
			"",
		},
		{
			"Homebrew test",
			config.GlobalConfig{
				BuildkitHost:      "127.0.0.1",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			errURLValidationFailure,
			"",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		_, err = stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			DebuggerHostFileValue:      tt.config.DebuggerHost,
			DebuggerPortFileValue:      tt.config.DebuggerPort,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			Console:                    logger,
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
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "tcp://localhost:8371",
			},
			"Buildkit and Local Registry URLs are pointed at different hosts",
		},
		{
			"Buildkit/Debugger host mismatch, schemes differ",
			config.GlobalConfig{
				BuildkitHost:      "docker-container://bk:1234",
				DebuggerHost:      "tcp://db:5678",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "",
			},
			"Buildkit and Debugger URLs are pointed at different hosts",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, conslogging.Info)
		logger = logger.WithWriter(&logs)

		frontend, err := NewStubFrontend(ctx, &FrontendConfig{})
		assert.NoError(t, err)

		stub, ok := frontend.(*stubFrontend)
		assert.True(t, ok)

		_, err = stub.setupAndValidateAddresses(FrontendDockerShell, &FrontendConfig{
			BuildkitHostFileValue:      tt.config.BuildkitHost,
			DebuggerHostFileValue:      tt.config.DebuggerHost,
			DebuggerPortFileValue:      tt.config.DebuggerPort,
			LocalRegistryHostFileValue: tt.config.LocalRegistryHost,
			Console:                    logger,
		})
		assert.NoError(t, err)
		assert.NotContains(t, logs.String(), tt.log)
	}
}
