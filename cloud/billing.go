package cloud

import (
	"context"
	"github.com/earthly/cloud-api/billing"

	"github.com/pkg/errors"
)

func (c *Client) GetBillingPlan(ctx context.Context, org string) (*billing.GetBillingPlanResponse, error) {
	response, err := c.billing.GetBillingPlan(c.withAuth(ctx), &billing.GetBillingPlanRequest{
		OrgName: org,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed get billing plan")
	}
	return response, nil
}
