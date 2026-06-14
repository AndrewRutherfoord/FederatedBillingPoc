package services

import (
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
)

// SyncCloudServiceProviderMetadata fetches the .well-known metadata from a cloud
// provider and upserts the provider record and its supported CSPs locally.
// Safe to call repeatedly — used both during account registration and by
// scheduled sync jobs.
func SyncCloudServiceProviderMetadata(repos *repository.Repos, baseURL string) (repository.BillingProvider, error) {
	client := clients.NewCloudServiceProviderClient(baseURL)

	metadata, err := client.GetMetadata()
	if err != nil {
		return repository.BillingProvider{}, err
	}

	return repos.BillingProvider.UpsertBillingProvider(metadata.ID, metadata.Name, baseURL, csps)
}
