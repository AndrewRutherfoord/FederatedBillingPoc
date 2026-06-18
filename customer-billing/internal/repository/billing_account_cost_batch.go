package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type BillingAccountCostBatchRepository interface {
	GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.BillingAccountCostBatch, error)
	Create(ctx context.Context, params CreateBillingAccountCostBatchParams) (*sqlcdb.BillingAccountCostBatch, error)
}

type billingAccountCostBatchRepo struct {
	db *db.DB
}

func newBillingAccountCostBatchRepo(database *db.DB) BillingAccountCostBatchRepository {
	return &billingAccountCostBatchRepo{db: database}
}

func (r *billingAccountCostBatchRepo) GetLatestBatchForBillingAccount(ctx context.Context, billingAccountID string) (*sqlcdb.BillingAccountCostBatch, error) {
	row, err := r.db.Queries.GetLatestBillingAccountCostBatchByAccount(ctx, billingAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

type CreateBillingAccountCostBatchParams struct {
	ID                     string
	BillingAccountID       string
	BillingPeriodID        string
	CloudServiceProviderID string
	TotalItems             int32
	TotalCost              float64
	BilledCurrency         string
	MerkelRoot             string
	BatchSignature         string
	CreatedAt              time.Time
	ReceivedAt             time.Time
}

func (r *billingAccountCostBatchRepo) Create(ctx context.Context, params CreateBillingAccountCostBatchParams) (*sqlcdb.BillingAccountCostBatch, error) {
	result, err := r.db.Queries.CreateBillingAccountCostBatch(ctx, sqlcdb.CreateBillingAccountCostBatchParams{
		ID:                     params.ID,
		BillingAccountID:       params.BillingAccountID,
		BillingPeriodID:        params.BillingPeriodID,
		CloudServiceProviderID: params.CloudServiceProviderID,
		TotalItems:             int64(params.TotalItems),
		TotalCost:              params.TotalCost,
		BilledCurrency:         params.BilledCurrency,
		MerkelRoot:             params.MerkelRoot,
		BatchSignature:         params.BatchSignature,
		CreatedAt:              params.CreatedAt,
		ReceivedAt:             params.ReceivedAt,
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}
