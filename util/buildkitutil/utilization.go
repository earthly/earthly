package buildkitutil

import "fmt"

// FormatUtilization returns a string representing utilization
func FormatUtilization(numOtherSessions, load, maxParallelism int) string {
	otherSessions := "unknown"
	if numOtherSessions >= 0 {
		otherSessions = fmt.Sprintf("%d", numOtherSessions)
	}
	return fmt.Sprintf("Utilization: %s other builds, %d/%d op load", otherSessions, load, maxParallelism)
}
