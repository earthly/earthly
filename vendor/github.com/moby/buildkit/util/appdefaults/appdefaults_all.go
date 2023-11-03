package appdefaults

import "time"

const (
	HealthAllowedFailures = 1
	HealthFrequency       = 1 * time.Second
	HealthTimeout         = 10 * time.Second
)
