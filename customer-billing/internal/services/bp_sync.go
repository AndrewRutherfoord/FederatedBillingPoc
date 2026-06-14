package services

import (
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
)

// SyncBillingProviderMetadata fetches the .well-known metadata from a billing
// provider and upserts the provider record and its supported CSPs locally.
// Safe to call repeatedly — used both during account registration and by
// scheduled sync jobs.
func SyncBillingProviderMetadata(repos *repository.Repos, baseURL string) (repository.BillingProvider, error) {
	client := clients.NewBillingProviderClient(baseURL)

	metadata, err := client.GetMetadata()
	if err != nil {
		return repository.BillingProvider{}, err
	}

	csps := make([]db.SupportedCloudProvider, len(metadata.SupportedCloudProviders))
	for i, csp := range metadata.SupportedCloudProviders {
		csps[i] = db.SupportedCloudProvider{
			ID:                csp.ID,
			BillingProviderID: metadata.ID,
			Name:              csp.Name,
			APIEndpointURL:    csp.APIEndpoint,
		}
	}

	return repos.BillingProvider.UpsertBillingProvider(metadata.ID, metadata.Name, baseURL, csps)
}
