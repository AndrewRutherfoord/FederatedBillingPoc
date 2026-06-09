package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"gorm.io/gorm"
)

// ProviderInfo holds the identity fields of this CSP instance, sourced from config.
type ProviderInfo struct {
	ID         string
	Name       string
	Currency   string
	RegionID   string
	RegionName string
}

// Repos is the single access point for all repositories.
// Pass it wherever data access is needed.
type Repos struct {
	Provider         ProviderInfo
	Customers        CustomerRepository
	Resources        ResourceRepository
	ResourceTypes    ResourceTypeRepository
	BillingProviders BillingProviderRepository
	BillingAccounts  BillingAccountRepository
	Focus            FocusRepository
	CostBatch        CostBatchRepository
	KeyValue         KeyValueRepository
}

func New(cfg *config.CspConfig, database *gorm.DB) *Repos {
	return &Repos{
		Provider: ProviderInfo{
			ID:         cfg.ProviderID,
			Name:       cfg.ProviderName,
			Currency:   cfg.Currency,
			RegionID:   cfg.RegionID,
			RegionName: cfg.RegionName,
		},
		Customers:        newCustomerRepo(database),
		Resources:        newResourceRepo(database),
		ResourceTypes:    newResourceTypeRepo(cfg.ResourceTypes),
		BillingProviders: newBillingProviderRepo(cfg.BillingProviders),
		BillingAccounts:  newBillingAccountRepo(database),
		Focus:            newFocusRepo(database),
		CostBatch:        newCostBatchRepo(database),
		KeyValue:         newKeyValueRepo(database),
	}
}
