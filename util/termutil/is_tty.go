package termutil

import "os"

// IsTTY returns true if a terminal is detected
func IsTTY() bool {
	return isFileDescriptorTTY(os.Stdin) && isFileDescriptorTTY(os.Stdout)
}

func isFileDescriptorTTY(fd *os.File) bool {
	if fileInfo, _ := fd.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
