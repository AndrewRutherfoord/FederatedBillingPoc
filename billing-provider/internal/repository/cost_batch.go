package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
)

type CostBatchRepository interface {
	Create(ctx context.Context, batch *db.CostBatch) error
	Get(ctx context.Context, id string) (*db.CostBatch, error)
	GetByBillingAccountAndTimeRange(ctx context.Context, billingAccountID string, from, to time.Time) ([]db.CostBatch, error)
}

type costBatchRepo struct {
	db *db.DB
}

func newCostBatchRepo(database *db.DB) CostBatchRepository {
	return &costBatchRepo{db: database}
}

func (r *costBatchRepo) Create(ctx context.Context, batch *db.CostBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

func (r *costBatchRepo) Get(ctx context.Context, id string) (*db.CostBatch, error) {
	var batch db.CostBatch
	if err := r.db.WithContext(ctx).First(&batch, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *costBatchRepo) GetByBillingAccountAndTimeRange(ctx context.Context, billingAccountID string, from, to time.Time) ([]db.CostBatch, error) {
	var batches []db.CostBatch
	err := r.db.WithContext(ctx).Where("billing_account_id = ? AND created_at >= ? AND created_at <= ?", billingAccountID, from, to).Find(&batches).Error
	return batches, err
}
