package fsutil

type VerboseProgressStatus int

const (
	StatusStat VerboseProgressStatus = iota
	StatusSkipped
	StatusSending
	StatusSent
	StatusReceiving
	StatusReceived
	StatusFailed
)

type VerboseProgressCB func(string, VerboseProgressStatus, int)
