package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	"gorm.io/gorm"
)

type SupportedCloudProvider struct {
	ID             string
	Name           string
	APIEndpointURL string
}

type BillingProvider struct {
	ID                      string
	Name                    string
	BaseURL                 string
	SupportedCloudProviders []SupportedCloudProvider
}

type BillingProviderRepository interface {
	ListBillingProviders() ([]*BillingProvider, error)
	GetBillingProviderByID(id string) (*BillingProvider, error)
	CreateBillingProvider(id string, name string, baseURL string) (BillingProvider, error)
	UpsertBillingProvider(id string, name string, baseURL string, csps []db.SupportedCloudProvider) (BillingProvider, error)
}

type billingProviderRepo struct {
	db *db.DB
}

func newBillingProviderRepo(database *db.DB) BillingProviderRepository {
	return &billingProviderRepo{db: database}
}

func (r *billingProviderRepo) loadCSPs(bpID string) ([]SupportedCloudProvider, error) {
	var rows []db.SupportedCloudProvider
	if err := r.db.Where("billing_provider_id = ?", bpID).Find(&rows).Error; err != nil {
		return nil, err
	}
	csps := make([]SupportedCloudProvider, len(rows))
	for i, row := range rows {
		csps[i] = SupportedCloudProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.APIEndpointURL}
	}
	return csps, nil
}

func toRepoBillingProvider(row db.BillingProvider, csps []SupportedCloudProvider) BillingProvider {
	return BillingProvider{ID: row.ID, Name: row.Name, BaseURL: row.BaseURL, SupportedCloudProviders: csps}
}

func (r *billingProviderRepo) ListBillingProviders() ([]*BillingProvider, error) {
	var rows []db.BillingProvider
	if err := r.db.Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*BillingProvider, len(rows))
	for i, row := range rows {
		csps, err := r.loadCSPs(row.ID)
		if err != nil {
			return nil, err
		}
		bp := toRepoBillingProvider(row, csps)
		result[i] = &bp
	}
	return result, nil
}

func (r *billingProviderRepo) GetBillingProviderByID(id string) (*BillingProvider, error) {
	var row db.BillingProvider
	if err := r.db.First(&row, "billing_provider_id = ?", id).Error; err != nil {
		return nil, err
	}
	csps, err := r.loadCSPs(id)
	if err != nil {
		return nil, err
	}
	bp := toRepoBillingProvider(row, csps)
	return &bp, nil
}

func (r *billingProviderRepo) CreateBillingProvider(id string, name string, baseURL string) (BillingProvider, error) {
	row := db.BillingProvider{ID: id, Name: name, BaseURL: baseURL}
	if err := r.db.Create(&row).Error; err != nil {
		return BillingProvider{}, err
	}
	return toRepoBillingProvider(row, nil), nil
}

func (r *billingProviderRepo) UpsertBillingProvider(id string, name string, baseURL string, csps []db.SupportedCloudProvider) (BillingProvider, error) {
	row := db.BillingProvider{ID: id, Name: name, BaseURL: baseURL}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&row).Error; err != nil {
			return err
		}
		if err := tx.Where("billing_provider_id = ?", id).Delete(&db.SupportedCloudProvider{}).Error; err != nil {
			return err
		}
		if len(csps) > 0 {
			return tx.Create(&csps).Error
		}
		return nil
	})
	if err != nil {
		return BillingProvider{}, err
	}

	repoCSPs := make([]SupportedCloudProvider, len(csps))
	for i, c := range csps {
		repoCSPs[i] = SupportedCloudProvider{ID: c.ID, Name: c.Name, APIEndpointURL: c.APIEndpointURL}
	}
	return toRepoBillingProvider(row, repoCSPs), nil
}
