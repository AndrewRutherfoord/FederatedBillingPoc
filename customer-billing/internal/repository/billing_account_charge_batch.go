package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type BillingAccountChargeBatchRepository interface {
	GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.BillingAccountChargeBatch, error)
	Create(ctx context.Context, params CreateBillingAccountChargeBatchParams) (*sqlcdb.BillingAccountChargeBatch, error)
	// UpsertInvoiced stores a batch known to belong to an invoiced billing period, inserting
	// it if the regular sync job hasn't seen it yet, or marking it invoiced if it has.
	UpsertInvoiced(ctx context.Context, params CreateBillingAccountChargeBatchParams) (*sqlcdb.BillingAccountChargeBatch, error)
}

type billingAccountChargeBatchRepo struct {
	db *db.DB
}

func newBillingAccountChargeBatchRepo(database *db.DB) BillingAccountChargeBatchRepository {
	return &billingAccountChargeBatchRepo{db: database}
}

func (r *billingAccountChargeBatchRepo) GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.BillingAccountChargeBatch, error) {
	row, err := r.db.Queries.GetLatestBillingAccountChargeBatchByAccount(ctx, billingAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

type CreateBillingAccountChargeBatchParams struct {
	ID                     string
	BillingAccountID       string
	BillingPeriodID        string
	CloudServiceProviderID string
	TotalItems             int32
	TotalCost              float64
	BilledCurrency         string
	MerkleRoot             string
	BatchSignature         string
	CreatedAt              time.Time
	ReceivedAt             time.Time
}

func (r *billingAccountChargeBatchRepo) Create(ctx context.Context, params CreateBillingAccountChargeBatchParams) (*sqlcdb.BillingAccountChargeBatch, error) {
	result, err := r.db.Queries.CreateBillingAccountChargeBatch(ctx, sqlcdb.CreateBillingAccountChargeBatchParams{
		ID:                     params.ID,
		BillingAccountID:       params.BillingAccountID,
		BillingPeriodID:        params.BillingPeriodID,
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

func (r *billingAccountChargeBatchRepo) UpsertInvoiced(ctx context.Context, params CreateBillingAccountChargeBatchParams) (*sqlcdb.BillingAccountChargeBatch, error) {
	result, err := r.db.Queries.UpsertInvoicedBillingAccountChargeBatch(ctx, sqlcdb.UpsertInvoicedBillingAccountChargeBatchParams{
		ID:                     params.ID,
		BillingAccountID:       params.BillingAccountID,
		BillingPeriodID:        params.BillingPeriodID,
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
