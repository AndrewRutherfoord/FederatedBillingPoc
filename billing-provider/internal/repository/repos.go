package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
)

type ProviderInfo struct {
	ID   string
	Name string
}

type Repos struct {
	Provider              ProviderInfo
	CloudServiceProviders CloudServiceProviderRepository
}

func New(cfg *config.Config, database *db.DB) *Repos {
	return &Repos{
		Provider: ProviderInfo{
			ID:   cfg.ProviderID,
			Name: cfg.ProviderName,
		},
		CloudServiceProviders: newCloudServiceProviderRepo(cfg.CloudServiceProviders),
	}
}
