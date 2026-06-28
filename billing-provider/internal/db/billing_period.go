package db

import (
	"time"
)

// BillingCycle is the cadence at which a billing account's billing periods
// recur, chosen by the account holder during onboarding.
type BillingCycle string

const (
	BillingCycleWeekly    BillingCycle = "weekly"
	BillingCycleMonthly   BillingCycle = "monthly"
	BillingCycleQuarterly BillingCycle = "quarterly"
	BillingCycleAnnual    BillingCycle = "annual"
)

func (c BillingCycle) Valid() bool {
	switch c {
	case BillingCycleWeekly, BillingCycleMonthly, BillingCycleQuarterly, BillingCycleAnnual:
		return true
	default:
		return false
	}
}

// Advance returns the end of the billing period that starts at periodStart
// for this cycle.
func (c BillingCycle) Advance(periodStart time.Time) time.Time {
	switch c {
	case BillingCycleWeekly:
		return periodStart.AddDate(0, 0, 7)
	case BillingCycleQuarterly:
		return periodStart.AddDate(0, 3, 0)
	case BillingCycleAnnual:
		return periodStart.AddDate(1, 0, 0)
	default:
		return periodStart.AddDate(0, 1, 0)
	}
}

// BillingPeriod is a single billing cycle for a billing account.
type BillingPeriod struct {
	ID               string `gorm:"column:id;primaryKey"`
	BillingAccountID string `gorm:"not null;index"`
	PeriodStart      time.Time
	PeriodEnd        time.Time
	Status           string
	OpenedAt         time.Time
	ClosedAt         time.Time
}
