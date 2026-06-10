package db

import (
	"time"
)

// Represents whether a billing account is linked to a customer a CSP.
type BillingAccountLink struct {
	AccountID                      string `gorm:"column:account_id;primaryKey"`
	CloudServiceProviderID         string `gorm:"column:cloud_service_provider_id;primaryKey"`
	CloudServiceProviderCustomerID string `gorm:"column:customer_id;not null;index"` // The CSPs customer ID for this billing account
	CreatedAt                      time.Time
}
