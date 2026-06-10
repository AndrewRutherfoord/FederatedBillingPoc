package db

import "time"

type Payment struct {
	ID          string `gorm:"column:id;primaryKey"`
	InvoiceID   string `gorm:"not null;index"`
	Amount      float64
	Currency    string
	GatewayRef  string // Reference from the payment gateway (e.g. transaction ID)
	InitiatedAt time.Time
	ExpiresAt   *time.Time // Nil if not yet settled
	PaidAt      *time.Time // Nil if not yet paid
}
