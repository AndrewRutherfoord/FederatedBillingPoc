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

type CloudServiceProviderAccount struct {
	ID                 string `gorm:"column:id;primaryKey"`
	BillingAccountID   string
	CloudProviderID    string
	InternalAccountID  string // Account ID used within the CSP
	OnboardingComplete *time.Time
	Created            *time.Time
}
