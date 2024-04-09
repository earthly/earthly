package cloud

import "context"

const (
	CloudStatusConnected = "Connected"
	CloudStatusActive    = "Active"
	CloudStatusProblem   = "Problem"
)

type Installation struct {
	Name          string
	Org           string
	Status        string
	NumSatellites int
	IsDefault     bool
}

func (c *Client) ConfigureCloud(ctx context.Context, orgName, cloudName string, setDefault bool) (*Installation, error) {
	// TODO
	return nil, nil
}

func (c *Client) ListClouds(ctx context.Context, orgName string) ([]Installation, error) {
	// TODO
	return nil, nil
}

func (c *Client) DeleteCloud(ctx context.Context, orgName, cloudName string) error {
	// TODO
	return nil
}
