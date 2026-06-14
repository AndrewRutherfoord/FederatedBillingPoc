package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type BillingProviderRepository interface {
	ListBillingProviders() ([]*db.BillingProvider, error)
	GetBillingProviderByID(id string) (*db.BillingProvider, error)
	CreateBillingProvider(id string, name string, baseURL string) (*db.BillingProvider, error)
}

type billingProviderRepo struct {
	db *db.DB
}

func newBillingProviderRepo(database *db.DB) BillingProviderRepository {
	return &billingProviderRepo{db: database}
}

func (r *billingProviderRepo) ListBillingProviders() ([]*db.BillingProvider, error) {
	var billingProviders []*db.BillingProvider
	if err := r.db.Find(&billingProviders).Error; err != nil {
		return nil, err
	}
	return billingProviders, nil
}

func (r *billingProviderRepo) GetBillingProviderByID(id string) (*db.BillingProvider, error) {
	var billingProvider db.BillingProvider
	if err := r.db.First(&billingProvider, "billing_provider_id = ?", id).Error; err != nil {
		return nil, err
	}
	return &billingProvider, nil
}

func (r *billingProviderRepo) CreateBillingProvider(id string, name string, baseURL string) (*db.BillingProvider, error) {
	billingProvider := &db.BillingProvider{
		ID:      id,
		Name:    name,
		BaseURL: baseURL,
	}
	if err := r.db.Create(billingProvider).Error; err != nil {
		return nil, err
	}
	return billingProvider, nil
}
