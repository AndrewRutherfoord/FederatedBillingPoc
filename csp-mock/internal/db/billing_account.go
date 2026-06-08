package db

import (
	"time"

	"gorm.io/gorm"
)

type BillingAccount struct {
	AccountID         string         `gorm:"column:account_id;primaryKey"`
	CustomerID        string         `gorm:"column:customer_id;not null;index"`
	BillingProviderID string         `gorm:"column:billing_provider_id;not null;index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}
