package billing

import (
	"fmt"
	"time"

	"github.com/earthly/cloud-api/billing"
)

func AddPlanInfo(plan *billing.BillingPlan, usedBuildMinutes time.Duration) {
	billingTracker.AddPlanInfo(plan, usedBuildMinutes)
}

func Plan() billing.BillingPlan {
	if billingTracker.plan != nil {
		return *billingTracker.plan
	}
	return billing.BillingPlan{}
}

func UsedBuildTime() time.Duration {
	return billingTracker.usedBuildTime
}

func GetBillingURL(hostName, orgName string) string {
	return fmt.Sprintf("%s/%s/settings", hostName, orgName)
}

func GetUpgradeURL(hostName, orgName string) string {
	return fmt.Sprintf("%s/%s/upgrade-now", hostName, orgName)
}
