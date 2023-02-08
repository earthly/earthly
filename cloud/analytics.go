package cloud

import (
	"context"

	"github.com/earthly/cloud-api/analytics"
	"github.com/pkg/errors"
)

// SendAnalytics send an analytics event to the Cloud server.
func (c *Client) SendAnalytics(ctx context.Context, data *analytics.SendAnalyticsRequest) error {
	if _, err := c.analytics.SendAnalytics(ctx, data); err != nil {
		return errors.Wrap(err, "failed sending analytics")
	}
	return nil
}
