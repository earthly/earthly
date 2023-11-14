package billing

import (
	"github.com/earthly/cloud-api/billing"
	"time"
)

type tracker struct {
	plan          *billing.BillingPlan
	usedBuildTime time.Duration
}

var billingTracker = tracker{}

func (bt *tracker) AddPlanInfo(plan *billing.BillingPlan, usedBuildTime time.Duration) {
	bt.plan = plan
	bt.usedBuildTime = usedBuildTime
}
