package db

import "time"

type ChargeBatch struct {
	ID                               string  `gorm:"primaryKey"` // ID is set by the CSP when the batch is created, and is used as the reference for all items in the batch. UUID4
	BillingAccountID                 string  `gorm:"not null;index"`
	BillingPeriodID                  *string `gorm:"index"` // Set once a billing job associates this batch with a billing period/invoice; null means not yet invoiced
	CloudServiceProviderSettlementID *string // Nullable foreign key to the settlement that this charge batch is part of (if it has been settled yet)
	CloudServiceProviderID           string  `gorm:"not null;index"`

	TotalItems int
	TotalCost  float64

	MerkleRoot     string
	BatchSignature string // Cloud Service Provider signature of MerkleRoot to ensure authenticity of the data

	CreatedAt  time.Time `gorm:"not null"`
	ReceivedAt time.Time
}
