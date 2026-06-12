package db

import "time"

type CostBatch struct {
	ID                 string    `gorm:"primaryKey"`
	BillingAccountID   string    `gorm:"not null;index"`
	BillingProviderID  string
	ResourceProviderID string
	MerkelRoot         string
	TotalItems         int
	TotalCost          float64
	BilledCurrency     string
	BatchSignature     string

	CreatedAt time.Time `gorm:"not null"`
}
