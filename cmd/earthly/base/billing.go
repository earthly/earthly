package base

import (
	"context"
	"fmt"
	"github.com/earthly/earthly/billing"
	"github.com/earthly/earthly/cloud"
	"time"
)

// CollectBillingInfo will collect billing plan info from billing service and make it available for other commands
// to use later
func (cli *CLI) CollectBillingInfo(ctx context.Context, cloudClient *cloud.Client, orgName string) error {
	if !cloudClient.IsLoggedIn(ctx) {
		return nil
	}
	resp, err := cloudClient.GetBillingPlan(ctx, orgName)
	if err != nil {
		return fmt.Errorf("failed to get billing plan: %w", err)
	}
	billing.AddPlanInfo(resp.GetPlan(), time.Second*time.Duration(resp.GetBillingCycleUsedBuildSeconds()))
	return nil
}
