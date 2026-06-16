package db

import "time"

type OnboardingSession struct {
	ID                string    `gorm:"primaryKey"`
	BillingProviderID string    `gorm:"not null"`
	BillingAccountID  string    `gorm:"not null"`
	ReturnURL         string    `gorm:"not null"`
	CreatedAt         time.Time
}
