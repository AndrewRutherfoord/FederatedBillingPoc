package db

import "time"

type BillingAccount struct {
	ID                    string `gorm:"column:account_id;primaryKey"`
	BillingProviderID     string
	Alias                 string
	Token                 string
	Created               *time.Time
	OnboardingRedirectURL string
	OnboardingComplete    *time.Time
}
