package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
	"github.com/google/uuid"
)

type CloudServiceProviderAccountLink struct {
	ID                string
	BillingAccountID  string
	CloudProviderID   string
	CloudProviderName string
}

type CloudServiceProviderAccountLinkWithTotalCost struct {
	CloudServiceProviderAccountLink
	TotalCost              float64
	BillingCurrency        string
	CustomerAPIEndpointURL string
}

type CloudServiceProviderAccountRepository interface {
	Create(ctx context.Context, billingAccountID string, cloudProviderID string) (CloudServiceProviderAccountLink, error)
	ListByBillingAccount(ctx context.Context, billingAccountID string) ([]CloudServiceProviderAccountLinkWithTotalCost, error)
}

type cloudServiceProviderAccountRepository struct {
	db *db.DB
}

func newCloudServiceProviderAccountRepo(database *db.DB) CloudServiceProviderAccountRepository {
	return &cloudServiceProviderAccountRepository{db: database}
}

func (r *cloudServiceProviderAccountRepository) Create(ctx context.Context, billingAccountID string, cloudProviderID string) (CloudServiceProviderAccountLink, error) {
	row, err := r.db.CreateCloudServiceProviderAccount(ctx, sqlcdb.CreateCloudServiceProviderAccountParams{
		ID:                     uuid.NewString(),
		BillingAccountID:       billingAccountID,
		CloudServiceProviderID: cloudProviderID,
	})
	if err != nil {
		return CloudServiceProviderAccountLink{}, err
	}
	return CloudServiceProviderAccountLink{
		ID:               row.ID,
		BillingAccountID: row.BillingAccountID,
		CloudProviderID:  row.CloudServiceProviderID,
	}, nil
}

func (r *cloudServiceProviderAccountRepository) ListByBillingAccount(ctx context.Context, billingAccountID string) ([]CloudServiceProviderAccountLinkWithTotalCost, error) {
	rows, err := r.db.ListCloudServiceProviderAccountsByBillingAccountWithTotalCost(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}
	result := make([]CloudServiceProviderAccountLinkWithTotalCost, len(rows))
	for i, row := range rows {
		result[i] = CloudServiceProviderAccountLinkWithTotalCost{
			CloudServiceProviderAccountLink: CloudServiceProviderAccountLink{
				ID:                row.ID,
				BillingAccountID:  row.BillingAccountID,
				CloudProviderID:   row.CloudServiceProviderID,
				CloudProviderName: row.CloudServiceProviderName,
			},
			TotalCost:              row.TotalCost.Float64,
			BillingCurrency:        row.BillingCurrency.String,
			CustomerAPIEndpointURL: row.CustomerApiEndpointUrl,
		}
	}
	return result, nil
}
