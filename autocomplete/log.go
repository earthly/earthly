package autocomplete

import (
	"fmt"
	"os"
)

var logPath string

func SetupLog(path string) {
	logPath = path
}

func Logf(format string, args ...interface{}) {
	Log(fmt.Sprintf(format, args...))
}

func Log(s string) {
	if logPath == "" {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.WriteString(s + "\n")
}
