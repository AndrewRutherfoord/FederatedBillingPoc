package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type CloudServiceProviderRepository interface {
	GetByID(ctx context.Context, id string) (CloudServiceProvider, error)
	Upsert(ctx context.Context, id string, name string, apiEndpointURL string) (CloudServiceProvider, error)
}

type cloudServiceProviderRepo struct {
	db *db.DB
}

func newCloudServiceProviderRepo(database *db.DB) CloudServiceProviderRepository {
	return &cloudServiceProviderRepo{db: database}
}

func (r *cloudServiceProviderRepo) GetByID(ctx context.Context, id string) (CloudServiceProvider, error) {
	row, err := r.db.Queries.GetBillingProviderSupportedCSP(ctx, id)
	if err != nil {
		return CloudServiceProvider{}, err
	}
	return CloudServiceProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.ApiEndpointUrl}, nil
}

func (r *cloudServiceProviderRepo) Upsert(ctx context.Context, id string, name string, apiEndpointURL string) (CloudServiceProvider, error) {
	row, err := r.db.Queries.UpdateBillingProviderSupportedCSP(ctx, sqlcdb.UpdateBillingProviderSupportedCSPParams{
		Name:           name,
		ApiEndpointUrl: apiEndpointURL,
		ID:             id,
	})
	if err != nil {
		return CloudServiceProvider{}, err
	}
	return CloudServiceProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.ApiEndpointUrl}, nil
}
