package termutil

import "os"

// IsTTY returns true if a terminal is detected
func IsTTY() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
