package db

import "time"

type CloudServiceProviderSettlement struct {
	ID                     string     `gorm:"primaryKey"`
	CloudServiceProviderID string     `gorm:"not null;index"`
	Amount                 float64    `gorm:"not null"`
	Currency               string     `gorm:"not null"`
	InitiatedAt            time.Time  `gorm:"not null"`
	SettledAt              *time.Time // Nil if not yet settled
	GatewayRef             *string    // Nil if not yet processed
}
