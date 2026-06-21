package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

// Stores what the CCSP sent for a charge batch retrived from CSP rather than via the BP for comparions
type CloudServiceProviderChargeBatchRepository interface {
	GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.CloudServiceProviderChargeBatch, error)
	Create(ctx context.Context, params CreateCloudServiceProviderChargeBatchParams) (*sqlcdb.CloudServiceProviderChargeBatch, error)
}

type cloudServiceProviderChargeBatchRepo struct {
	db *db.DB
}

func newCloudServiceProviderChargeBatchRepo(database *db.DB) CloudServiceProviderChargeBatchRepository {
	return &cloudServiceProviderChargeBatchRepo{db: database}
}

func (r *cloudServiceProviderChargeBatchRepo) GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.CloudServiceProviderChargeBatch, error) {
	row, err := r.db.Queries.GetLatestCloudServiceProviderChargeBatchByAccount(ctx, billingAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

type CreateCloudServiceProviderChargeBatchParams struct {
	ID                     string
	BillingAccountID       string
	CloudServiceProviderID string
	TotalItems             int32
	TotalCost              float64
	BilledCurrency         string
	MerkleRoot             string
	BatchSignature         string
	CreatedAt              time.Time
	ReceivedAt             time.Time
}

func (r *cloudServiceProviderChargeBatchRepo) Create(ctx context.Context, params CreateCloudServiceProviderChargeBatchParams) (*sqlcdb.CloudServiceProviderChargeBatch, error) {
	result, err := r.db.Queries.CreateCloudServiceProviderChargeBatch(ctx, sqlcdb.CreateCloudServiceProviderChargeBatchParams{
		ID:                     params.ID,
		BillingAccountID:       params.BillingAccountID,
		CloudServiceProviderID: params.CloudServiceProviderID,
		TotalItems:             int64(params.TotalItems),
		TotalCost:              params.TotalCost,
		BilledCurrency:         params.BilledCurrency,
		MerkleRoot:             params.MerkleRoot,
		BatchSignature:         params.BatchSignature,
		CreatedAt:              params.CreatedAt,
		ReceivedAt:             params.ReceivedAt,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
