package billing

import (
	"fmt"
	"time"

	"github.com/earthly/cloud-api/billing"
	"google.golang.org/protobuf/proto"
)

// AddPlanInfo sets the billing plan to make it available to retrieve later with Plan()
func AddPlanInfo(plan *billing.BillingPlan, usedBuildMinutes time.Duration) {
	billingTracker.AddPlanInfo(plan, usedBuildMinutes)
}

// Plan returns a copy of the billing plan which should have been set earlier by AddPlanInfo
func Plan() *billing.BillingPlan {
	if billingTracker.plan != nil {
		return proto.Clone(billingTracker.plan).(*billing.BillingPlan)
	}
	return nil
}

// UsedBuildTime returns the used build time(duration) of org referenced by the earthly command
func UsedBuildTime() time.Duration {
	return billingTracker.usedBuildTime
}

// GetBillingURL returns the billing url based on the given host and org names
func GetBillingURL(hostName, orgName string) string {
	return fmt.Sprintf("%s/%s/settings", hostName, orgName)
}

// GetUpgradeURL returns the billing url based on the given host and org names
func GetUpgradeURL(hostName, orgName string) string {
	return fmt.Sprintf("%s/%s/upgrade-now", hostName, orgName)
}
