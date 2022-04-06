package secretprovider

import (
	"github.com/google/uuid"
)

type secretID struct {
	name          string
	localProject  string
	remoteProject string
}

var secretIDs map[string]secretID

// NewSecretID is used to store secretID metadata and returns a uuid
// which can be used to lookup this metadata
// This idirection is required due to earthly not being able to (easily) pass this data through buildkit
func NewSecretID(name, localProject, remoteProject string) string {
	if secretIDs == nil {
		secretIDs = map[string]secretID{}
	}
	id := uuid.New().String()
	secretIDs[id] = secretID{
		name:          name,
		localProject:  localProject,
		remoteProject: remoteProject,
	}
	return id
}
