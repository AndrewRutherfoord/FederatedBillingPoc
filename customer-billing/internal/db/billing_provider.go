package db

type BillingProvider struct {
	ID      string `gorm:"column:billing_provider_id;primaryKey"`
	Name    string
	BaseURL string
}

type SupportedCloudProvider struct {
	ID                string `gorm:"column:cloud_provider_id;primaryKey"`
	BillingProviderID string
	Name              string
	APIEndpointURL    string
}
