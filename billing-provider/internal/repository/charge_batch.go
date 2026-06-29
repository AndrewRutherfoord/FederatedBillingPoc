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
	// GetByBillingAccountAndPeriod paginates a billing period's batches, oldest first.
	GetByBillingAccountAndPeriod(ctx context.Context, billingAccountID, billingPeriodID string, limit, offset int) ([]db.ChargeBatch, error)
	// GetUnassociatedByBillingAccount returns batches not yet attributed to a billing period, oldest first.
	GetUnassociatedByBillingAccount(ctx context.Context, billingAccountID string) ([]db.ChargeBatch, error)
	AssociateWithPeriod(ctx context.Context, batchIDs []string, billingPeriodID string) error
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

func (r *chargeBatchRepo) GetByBillingAccountAndPeriod(ctx context.Context, billingAccountID, billingPeriodID string, limit, offset int) ([]db.ChargeBatch, error) {
	var batches []db.ChargeBatch
	err := r.db.WithContext(ctx).
		Where("billing_account_id = ? AND billing_period_id = ?", billingAccountID, billingPeriodID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&batches).Error
	return batches, err
}

func (r *chargeBatchRepo) GetUnassociatedByBillingAccount(ctx context.Context, billingAccountID string) ([]db.ChargeBatch, error) {
	var batches []db.ChargeBatch
	err := r.db.WithContext(ctx).
		Where("billing_account_id = ? AND billing_period_id IS NULL", billingAccountID).
		Order("created_at ASC").
		Find(&batches).Error
	return batches, err
}

func (r *chargeBatchRepo) AssociateWithPeriod(ctx context.Context, batchIDs []string, billingPeriodID string) error {
	if len(batchIDs) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Model(&db.ChargeBatch{}).
		Where("id IN ?", batchIDs).
		Update("billing_period_id", billingPeriodID).Error
}
