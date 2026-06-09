package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/config"
)

type CloudServiceProviderRepository interface {
	List(ctx context.Context) ([]config.CloudServiceProvider, error)
	Get(ctx context.Context, id string) (*config.CloudServiceProvider, error)
}

type cloudServiceProviderRepo struct {
	items []config.CloudServiceProvider
}

func newCloudServiceProviderRepo(items []config.CloudServiceProvider) CloudServiceProviderRepository {
	return &cloudServiceProviderRepo{items: items}
}

func (r *cloudServiceProviderRepo) List(ctx context.Context) ([]config.CloudServiceProvider, error) {
	return r.items, nil
}

func (r *cloudServiceProviderRepo) Get(ctx context.Context, id string) (*config.CloudServiceProvider, error) {
	for i := range r.items {
		if r.items[i].ID == id {
			return &r.items[i], nil
		}
	}
	return nil, nil // Not found; return nil without an error
}
