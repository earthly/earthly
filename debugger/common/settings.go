package common

// DebuggerSettings is used to pass settings to the debugger
type DebuggerSettings struct {
	DebugLevelLogging bool   `json:"debug_level"`
	RemoteConsoleAddr string `json:"remote_console_addr"`
	Enabled           bool   `json:"enabled"`
}
