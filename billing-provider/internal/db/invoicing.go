package db

import (
	"time"
)

type Invoice struct {
	ID               string    `gorm:"column:id;primaryKey"`
	BillingAccountID string    `gorm:"not null;index"`
	BillingPeriodID  string    `gorm:"not null;index"`
	Amount           float64   `gorm:"not null"`
	Currency         string    `gorm:"not null"`
	Status           string    `gorm:"not null"`
	IssuedAt         time.Time `gorm:"not null"`
	DueAt            time.Time `gorm:"not null"` // If not paid by this date, invoice is considered overdue
}

type InvoiceBatch struct {
	ID                     string  `gorm:"column:id;primaryKey"`
	InvoiceID              string  `gorm:"not null;index"`
	Amount                 float64 `gorm:"not null"`
	CloudServiceProviderID string  `gorm:"not null;index"`
	CostBatchID            string  `gorm:"not null;index"` // Link back to the cost batch that generated this line item
}
