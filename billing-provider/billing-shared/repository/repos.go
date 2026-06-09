package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/config"
	"gorm.io/gorm"
)

type ProviderInfo struct {
	ID   string
	Name string
}

type Repos struct {
	Provider              ProviderInfo
	CloudServiceProviders CloudServiceProviderRepository
}

func New(cfg *config.BillingProviderConfig, database *gorm.DB) *Repos {
	return &Repos{
		Provider: ProviderInfo{
			ID:   cfg.ProviderID,
			Name: cfg.ProviderName,
		},
		CloudServiceProviders: newCloudServiceProviderRepo(cfg.CloudServiceProviders),
	}
}
