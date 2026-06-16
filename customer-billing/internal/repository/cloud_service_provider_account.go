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

type CloudServiceProviderAccountRepository interface {
	Create(ctx context.Context, billingAccountID string, cloudProviderID string) (CloudServiceProviderAccountLink, error)
	ListByBillingAccount(ctx context.Context, billingAccountID string) ([]CloudServiceProviderAccountLink, error)
}

type cloudServiceProviderAccountRepository struct {
	db *db.DB
}

func newCloudServiceProviderAccountRepo(database *db.DB) CloudServiceProviderAccountRepository {
	return &cloudServiceProviderAccountRepository{db: database}
}

func (r *cloudServiceProviderAccountRepository) Create(ctx context.Context, billingAccountID string, cloudProviderID string) (CloudServiceProviderAccountLink, error) {
	row, err := r.db.CreateCloudServiceProviderAccount(ctx, sqlcdb.CreateCloudServiceProviderAccountParams{
		ID:               uuid.NewString(),
		BillingAccountID: billingAccountID,
		CloudProviderID:  cloudProviderID,
	})
	if err != nil {
		return CloudServiceProviderAccountLink{}, err
	}
	return CloudServiceProviderAccountLink{
		ID:               row.ID,
		BillingAccountID: row.BillingAccountID,
		CloudProviderID:  row.CloudProviderID,
	}, nil
}

func (r *cloudServiceProviderAccountRepository) ListByBillingAccount(ctx context.Context, billingAccountID string) ([]CloudServiceProviderAccountLink, error) {
	rows, err := r.db.ListCloudServiceProviderAccountsByBillingAccount(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}
	result := make([]CloudServiceProviderAccountLink, len(rows))
	for i, row := range rows {
		result[i] = CloudServiceProviderAccountLink{
			ID:                row.ID,
			BillingAccountID:  row.BillingAccountID,
			CloudProviderID:   row.CloudProviderID,
			CloudProviderName: row.CloudProviderName,
		}
	}
	return result, nil
}
