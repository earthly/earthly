package conslogging

import (
	"time"
)

// These types are also used by our log server, too.
const (
	// StatusWaiting is the status for a target that has yet to execute.
	StatusWaiting = "waiting"

	// StatusInProgress is the status for a target that is currently executing.
	StatusInProgress = "in_progress"

	// StatusComplete is the status for a target that has run to completion.
	StatusComplete = "complete"

	// StatusCancelled is the status for a target that did not run to completion, and was interrupted.
	StatusCancelled = "cancelled"

	// ResultSuccess is the result for a target that exits successfully.
	ResultSuccess = "success"

	// ResultFailure is the result for a target that exited with some kind of error code.
	ResultFailure = "failure"

	// ResultCancelled is the results for a target that did not run to completion.
	ResultCancelled = "cancelled"
)

// Manifest is the structure for the log bundle manifest, including all overarching data we need.
type Manifest struct {
	Version    int              `json:"version"`
	Duration   int              `json:"duration"`
	Status     string           `json:"status"`
	Result     string           `json:"result"`
	CreatedAt  time.Time        `json:"created_at"`
	Targets    []TargetManifest `json:"targets"`
	Entrypoint string           `json:"entrypoint"`
}

// TargetManifest is the structure for an individual target, indicating all relevant information.
type TargetManifest struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Result   string `json:"result"`
	Duration int    `json:"duration"`
	Size     int    `json:"size"`
	Command  string `json:"command,omitempty"`
	Summary  string `json:"summary,omitempty"`
}

// Permissions is the structure for the permissions manifest that can grant view rights to other Earthly users.
type Permissions struct {
	Version int      `json:"version"`
	Users   []string `json:"users"`
	Orgs    []string `json:"orgs"`
}
