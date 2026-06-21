package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type CloudServiceProvider struct {
	ID                     string
	Name                   string
	APIEndpointURL         string
	CustomerAPIEndpointURL string
}

type BillingProvider struct {
	ID                      string
	Name                    string
	BaseURL                 string
	SupportedCloudProviders []CloudServiceProvider
}

type BillingProviderRepository interface {
	ListBillingProviders(ctx context.Context) ([]*BillingProvider, error)
	GetBillingProviderByID(ctx context.Context, id string) (*BillingProvider, error)
	CreateBillingProvider(ctx context.Context, id string, name string, baseURL string) (BillingProvider, error)
	UpsertBillingProvider(ctx context.Context, id string, name string, baseURL string, csps []CloudServiceProvider) (BillingProvider, error)
}

type billingProviderRepo struct {
	db *db.DB
}

func newBillingProviderRepo(database *db.DB) BillingProviderRepository {
	return &billingProviderRepo{db: database}
}

func toBillingProvider(row sqlcdb.BillingProvider, cspRows []sqlcdb.CloudServiceProvider) BillingProvider {
	csps := make([]CloudServiceProvider, len(cspRows))
	for i, c := range cspRows {
		csps[i] = CloudServiceProvider{ID: c.ID, Name: c.Name, APIEndpointURL: c.ApiEndpointUrl, CustomerAPIEndpointURL: c.CustomerApiEndpointUrl}
	}
	return BillingProvider{ID: row.ID, Name: row.Name, BaseURL: row.ApiEndpointUrl, SupportedCloudProviders: csps}
}

func (r *billingProviderRepo) ListBillingProviders(ctx context.Context) ([]*BillingProvider, error) {
	rows, err := r.db.Queries.ListBillingProviders(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*BillingProvider, len(rows))
	for i, row := range rows {
		csps, err := r.db.Queries.ListCloudServiceProvidersByBillingProvider(ctx, row.ID)
		if err != nil {
			return nil, err
		}
		bp := toBillingProvider(row, csps)
		result[i] = &bp
	}
	return result, nil
}

func (r *billingProviderRepo) GetBillingProviderByID(ctx context.Context, id string) (*BillingProvider, error) {
	row, err := r.db.Queries.GetBillingProvider(ctx, id)
	if err != nil {
		return nil, err
	}
	csps, err := r.db.Queries.ListCloudServiceProvidersByBillingProvider(ctx, id)
	if err != nil {
		return nil, err
	}
	bp := toBillingProvider(row, csps)
	return &bp, nil
}

func (r *billingProviderRepo) CreateBillingProvider(ctx context.Context, id string, name string, baseURL string) (BillingProvider, error) {
	err := r.db.Queries.CreateBillingProvider(ctx, sqlcdb.CreateBillingProviderParams{
		ID:             id,
		Name:           name,
		ApiEndpointUrl: baseURL,
	})
	if err != nil {
		return BillingProvider{}, err
	}
	return BillingProvider{ID: id, Name: name, BaseURL: baseURL}, nil
}

func (r *billingProviderRepo) UpsertBillingProvider(ctx context.Context, id string, name string, baseURL string, csps []CloudServiceProvider) (BillingProvider, error) {
	err := r.db.WithTx(ctx, func(q *sqlcdb.Queries) error {
		if err := q.UpsertBillingProvider(ctx, sqlcdb.UpsertBillingProviderParams{
			ID:             id,
			Name:           name,
			ApiEndpointUrl: baseURL,
		}); err != nil {
			return err
		}
		if err := q.DeleteCloudServiceProvidersByBillingProvider(ctx, id); err != nil {
			return err
		}
		for _, csp := range csps {
			if err := q.CreateCloudServiceProvider(ctx, sqlcdb.CreateCloudServiceProviderParams{
				ID:                     csp.ID,
				BillingProviderID:      id,
				Name:                   csp.Name,
				ApiEndpointUrl:         csp.APIEndpointURL,
				CustomerApiEndpointUrl: csp.CustomerAPIEndpointURL,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return BillingProvider{}, err
	}
	return BillingProvider{ID: id, Name: name, BaseURL: baseURL, SupportedCloudProviders: csps}, nil
}
