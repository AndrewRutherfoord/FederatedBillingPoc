package services

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
)

// SyncCloudServiceProviderMetadata fetches the .well-known metadata from a cloud
// service provider and upserts it locally.
// Safe to call repeatedly — used both during account registration and by
// scheduled sync jobs.
func SyncCloudServiceProviderMetadata(ctx context.Context, repos *repository.Repos, baseURL string) (repository.CloudServiceProvider, error) {
	client := clients.NewCloudServiceProviderClient(baseURL)

	metadata, err := client.GetMetadata()
	if err != nil {
		return repository.CloudServiceProvider{}, err
	}

	return repos.CloudServiceProvider.Upsert(ctx, metadata.ID, metadata.Name, metadata.APIEndpointURL)
}
