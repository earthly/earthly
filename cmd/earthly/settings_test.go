package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/config"
	"github.com/earthly/earthly/conslogging"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

var noopArgs = []string{""}

type results struct {
	buildkit      string
	debugger      string
	localRegistry string
}

func TestBuildArgMatrix(t *testing.T) {
	var tests = []struct {
		testName string
		config   config.GlobalConfig
		args     []string
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
				buildkit:      buildkitd.DockerAddress,
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
				buildkit:      buildkitd.DockerAddress,
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
				buildkit:      buildkitd.DockerAddress,
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
			[]string{"", "--buildkit-host", "tcp://ok-bk:42", "--debugger-host", "tcp://ok-db:43"},
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
		var trash strings.Builder

		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, false)
		logger = logger.WithWriter(&logs)

		earthlyApp := newEarthlyApp(ctx, logger)
		earthlyApp.cfg = &config.Config{Global: tt.config}
		earthlyApp.cliApp.Writer = &trash    // Just chuck the help output
		earthlyApp.cliApp.ErrWriter = &trash // All of it, we dont care

		// Before is called at about the time that we would parse these, plus it a nice place to hook.
		earthlyApp.cliApp.Before = func(context *cli.Context) error {

			err := earthlyApp.setupAndValidateAddresses(context)
			assert.NoError(t, err, tt.testName)
			assert.Equal(t, tt.expected, results{
				buildkit:      earthlyApp.buildkitHost,
				debugger:      earthlyApp.debuggerHost,
				localRegistry: earthlyApp.localRegistryHost,
			}, tt.testName)

			return nil
		}
		earthlyApp.cliApp.RunContext(ctx, tt.args)
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
			"Buildkit and Local Registry URLs are pointed at different hosts",
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
			"Buildkit and Debugger URLs are pointed at different hosts",
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
			"Local registry with remote buildkit",
			config.GlobalConfig{
				BuildkitHost:      "tcp://cool-host:1234",
				DebuggerHost:      "",
				DebuggerPort:      config.DefaultDebuggerPort,
				LocalRegistryHost: "tcp://127.0.0.1:8373",
			},
			nil,
			"Local registry host is specified while using remote buildkit. Local registry will not be used.",
		},
	}

	for _, tt := range tests {
		ctx := context.Background()

		var logs strings.Builder
		var output strings.Builder

		logger := conslogging.Current(conslogging.NoColor, conslogging.DefaultPadding, false)
		logger = logger.WithWriter(&logs)

		earthlyApp := newEarthlyApp(ctx, logger)
		earthlyApp.cfg = &config.Config{Global: tt.config}
		earthlyApp.cliApp.Writer = &output
		earthlyApp.cliApp.ErrWriter = &output

		// Before is called at about the time that we would parse these, plus it a nice place to hook.
		earthlyApp.cliApp.Before = func(context *cli.Context) error {
			err := earthlyApp.setupAndValidateAddresses(context)

			assert.ErrorIs(t, err, tt.expected)
			assert.Contains(t, logs.String(), tt.log)

			return nil
		}
		earthlyApp.cliApp.RunContext(ctx, []string{""})
	}
}

func TestParseAndvalidateURLFailures(t *testing.T) {
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
		_, err := parseAndvalidateURL(tt.url)
		assert.ErrorIs(t, err, tt.expected)
	}
}

func TestParseAndvalidateURL(t *testing.T) {
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
		_, err := parseAndvalidateURL(tt.url)
		assert.NoError(t, err)
	}
}
