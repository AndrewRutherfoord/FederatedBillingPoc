package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type CloudServiceProvider struct {
	ID             string
	Name           string
	APIEndpointURL string
}

type CloudServiceProviderRepository interface {
	GetByID(id string) (CloudServiceProvider, error)
	Create(name string, APIEndpointURL string) (CloudServiceProvider, error)
}

type cloudServiceProvider struct {
	db *db.DB
}

func newCloudServiceProviderRepo(database *db.DB) CloudServiceProviderRepository {
	return &cloudServiceProvider{db: database}
}

func (r *cloudServiceProvider) GetByID(id string) (CloudServiceProvider, error) {
	var row db.CloudServiceProvider
	if err := r.db.Where("cloud_provider_id = ?", id).First(&row).Error; err != nil {
		return CloudServiceProvider{}, err
	}
	return CloudServiceProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.BaseURL}, nil
}

func (r *cloudServiceProvider) Create(name string, APIEndpointURL string) (CloudServiceProvider, error) {
	newCSP := db.CloudServiceProvider{Name: name, BaseURL: APIEndpointURL}
	if err := r.db.Create(&newCSP).Error; err != nil {
		return CloudServiceProvider{}, err
	}
	return CloudServiceProvider{ID: newCSP.ID, Name: newCSP.Name, APIEndpointURL: newCSP.BaseURL}, nil
}
