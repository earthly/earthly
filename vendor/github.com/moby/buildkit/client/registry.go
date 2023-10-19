package client

import registry "github.com/moby/buildkit/api/services/registry"

// RegistryClient creates a new gRPC client for the embedded Docker
// registry. The client & server use a gRPC data stream to proxy image pull
// requests to the embedded registry. Earthly-specific.
func (c *Client) RegistryClient() registry.RegistryClient {
	return registry.NewRegistryClient(c.conn)
}
