package db

import (
	"time"
)

type BillingAccount struct {
	ID             string `gorm:"column:account_id;primaryKey"`
	Name           string
	Email          string
	ApprovalStatus string // e.g., "pending", "approved", "rejected"
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
