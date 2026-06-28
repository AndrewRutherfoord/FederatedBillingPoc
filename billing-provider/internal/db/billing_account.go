package db

import (
	"time"
)

type BillingAccount struct {
	ID             string `gorm:"column:account_id;primaryKey"`
	Name           string
	Email          string
	BillingCycle   BillingCycle // Cadence of recurring billing periods, chosen during onboarding
	ApprovalStatus string       // e.g., "pending", "approved", "rejected"
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
