package db

import "time"

type CostBatch struct {
	ID               string `gorm:"primaryKey"`
	BillingAccountID string `gorm:"not null;index"`
	MerkelRoot       string
	TotalItems       int
	TotalCost        float64

	CreatedAt time.Time `gorm:"not null"`
}
