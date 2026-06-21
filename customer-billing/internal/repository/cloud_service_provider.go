package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type CloudServiceProviderRepository interface {
	GetByID(ctx context.Context, id string) (CloudServiceProvider, error)
	Upsert(ctx context.Context, id string, name string, customerAPIEndpointURL string) (CloudServiceProvider, error)
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
	return CloudServiceProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.ApiEndpointUrl, CustomerAPIEndpointURL: row.CustomerApiEndpointUrl}, nil
}

// Upsert refreshes a CSP's customer-facing endpoint with the value the CSP itself
// reports via its own .well-known metadata, which is authoritative over whatever
// the billing provider initially advertised.
func (r *cloudServiceProviderRepo) Upsert(ctx context.Context, id string, name string, customerAPIEndpointURL string) (CloudServiceProvider, error) {
	row, err := r.db.Queries.UpdateBillingProviderSupportedCSPCustomerEndpoint(ctx, sqlcdb.UpdateBillingProviderSupportedCSPCustomerEndpointParams{
		Name:                   name,
		CustomerApiEndpointUrl: customerAPIEndpointURL,
		ID:                     id,
	})
	if err != nil {
		return CloudServiceProvider{}, err
	}
	return CloudServiceProvider{ID: row.ID, Name: row.Name, APIEndpointURL: row.ApiEndpointUrl, CustomerAPIEndpointURL: row.CustomerApiEndpointUrl}, nil
}
