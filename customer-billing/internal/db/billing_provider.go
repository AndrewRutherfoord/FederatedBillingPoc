package db

type BillingProvider struct {
	ID      string `gorm:"column:billing_provider_id;primaryKey"`
	Name    string
	BaseURL string
}
