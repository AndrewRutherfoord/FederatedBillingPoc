package db

import "time"

type ChargeBatch struct {
	ID                     string `gorm:"primaryKey"`
	BillingAccountID       string `gorm:"not null;index"`
	BillingProviderID      string
	CloudServiceProviderID string
	MerkleRoot             string
	TotalItems             int
	TotalCost              float64
	BilledCurrency         string
	BatchSignature         string

	CreatedAt time.Time `gorm:"not null"`
}
