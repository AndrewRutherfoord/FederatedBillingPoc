package db

import "time"

type CostBatch struct {
	ID                               string  `gorm:"primaryKey"` // ID is set by the CSP when the batch is created, and is used as the reference for all items in the batch. UUID4
	BillingAccountID                 string  `gorm:"not null;index"`
	BillingPeriodID                  string  `gorm:"not null;index"`
	CloudServiceProviderSettlementID *string // Nullable foreign key to the settlement that this cost batch is part of (if it has been settled yet)
	CloudServiceProviderID           string  `gorm:"not null;index"`

	TotalItems int
	TotalCost  float64

	MerkelRoot string
	Signature  string // Cloud Service Provider Signature of MerkelRoot to ensure authenticity of the data

	CreatedAt time.Time `gorm:"not null"`
}
