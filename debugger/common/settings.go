package common

// DebuggerSettingsSecretsKey stores the secrets key name
const DebuggerSettingsSecretsKey = "earthly_debugger_settings"

// DebuggerDefaultSocketPath is the default socket to connect to (path is inside the container)
const DebuggerDefaultSocketPath = "/var/run/earthly_interactive"

// DebuggerSettings is used to pass settings to the debugger
type DebuggerSettings struct {
	DebugLevelLogging bool   `json:"debugLevel"`
	Enabled           bool   `json:"enabled"`
	SocketPath        string `json:"socketPath"`
	Term              string `json:"term"`
}
