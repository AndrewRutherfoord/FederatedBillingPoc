package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
)

type ChargeBatchRepository interface {
	Create(ctx context.Context, batch *db.ChargeBatch) error
	Get(ctx context.Context, id string) (*db.ChargeBatch, error)
	GetByBillingAccountAndTimeRange(ctx context.Context, billingAccountID string, from, to time.Time) ([]db.ChargeBatch, error)
}

type chargeBatchRepo struct {
	db *db.DB
}

func newChargeBatchRepo(database *db.DB) ChargeBatchRepository {
	return &chargeBatchRepo{db: database}
}

func (r *chargeBatchRepo) Create(ctx context.Context, batch *db.ChargeBatch) error {
	return r.db.WithContext(ctx).Create(batch).Error
}

func (r *chargeBatchRepo) Get(ctx context.Context, id string) (*db.ChargeBatch, error) {
	var batch db.ChargeBatch
	if err := r.db.WithContext(ctx).First(&batch, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *chargeBatchRepo) GetByBillingAccountAndTimeRange(ctx context.Context, billingAccountID string, from, to time.Time) ([]db.ChargeBatch, error) {
	var batches []db.ChargeBatch
	err := r.db.WithContext(ctx).Where("billing_account_id = ? AND created_at >= ? AND created_at <= ?", billingAccountID, from, to).Find(&batches).Error
	return batches, err
}
