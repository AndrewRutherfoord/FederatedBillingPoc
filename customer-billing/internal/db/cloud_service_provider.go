package db

type CloudServiceProvider struct {
	ID      string `gorm:"column:cloud_provider_id;primaryKey"`
	Name    string
	BaseURL string
}
