package conslogging

import (
	"time"
)

// These types are also used by our log server, too.
const (
	StatusWaiting    = "waiting"
	StatusInProgress = "in_progress"
	StatusComplete   = "complete"
	StatusCancelled  = "cancelled"

	ResultSuccess   = "success"
	ResultFailure   = "failure"
	ResultCancelled = "cancelled"
)

type Manifest struct {
	Version    int              `json:"version"`
	Duration   int              `json:"duration"`
	Status     string           `json:"status"`
	Result     string           `json:"result"`
	CreatedAt  time.Time        `json:"created_at"`
	Targets    []TargetManifest `json:"targets"`
	Entrypoint string           `json:"entrypoint"`
}

type TargetManifest struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Result   string `json:"result"`
	Duration int    `json:"duration"`
	Size     int    `json:"size"`
	Command  string `json:"command,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

type Permissions struct {
	Version int      `json:"version"`
	Users   []uint64 `json:"users"`
	Orgs    []uint64 `json:"orgs"`
}
