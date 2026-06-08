package db

import (
	"time"

	"github.com/shopspring/decimal"
)

type Resource struct {
	ID               string `gorm:"primaryKey"`
	CustomerID       string
	BillingAccountID string
	ResourceType     string
	StartedAt        time.Time
	DeletedAt        *time.Time       // Soft delete timestamp
	StorageGB        *decimal.Decimal `gorm:"type:text"`
}
