package repository

import (
	"context"
	"errors"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
)

type BillingProviderRepository interface {
	List(ctx context.Context) ([]config.BillingProvider, error)
	Get(ctx context.Context, id string) (*config.BillingProvider, error)
}

type billingProviderRepo struct {
	items []config.BillingProvider
}

func newBillingProviderRepo(items []config.BillingProvider) BillingProviderRepository {
	return &billingProviderRepo{items: items}
}

func (r *billingProviderRepo) List(ctx context.Context) ([]config.BillingProvider, error) {
	return r.items, nil
}

func (r *billingProviderRepo) Get(ctx context.Context, id string) (*config.BillingProvider, error) {
	for i := range r.items {
		if r.items[i].ID == id {
			return &r.items[i], nil
		}
	}
	return nil, errors.New("billing provider not found")
}
