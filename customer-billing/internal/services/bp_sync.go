package services

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
)

// SyncBillingProviderMetadata fetches the .well-known metadata from a billing
// provider and upserts the provider record and its supported CSPs locally.
// Safe to call repeatedly — used both during account registration and by
// scheduled sync jobs.
func SyncBillingProviderMetadata(ctx context.Context, repos *repository.Repos, baseURL string) (repository.BillingProvider, error) {
	client := clients.NewBillingProviderClient(baseURL)

	metadata, err := client.GetMetadata()
	if err != nil {
		return repository.BillingProvider{}, err
	}

	csps := make([]repository.CloudServiceProvider, len(metadata.SupportedCloudProviders))
	for i, csp := range metadata.SupportedCloudProviders {
		csps[i] = repository.CloudServiceProvider{
			ID:             csp.ID,
			Name:           csp.Name,
			APIEndpointURL: csp.APIEndpoint,
		}
	}

	return repos.BillingProvider.UpsertBillingProvider(ctx, metadata.ID, metadata.Name, baseURL, csps)
}
