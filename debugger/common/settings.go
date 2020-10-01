package common

// DebuggerSettingsSecretsKey stores the secrets key name
const DebuggerSettingsSecretsKey = "earthly_debugger_settings"

// DebuggerSettings is used to pass settings to the debugger
type DebuggerSettings struct {
	DebugLevelLogging bool   `json:"debugLevel"`
	Enabled           bool   `json:"enabled"`
	SockPath          string `json:"sockPath"`
	Term              string `json:"term"`
}
