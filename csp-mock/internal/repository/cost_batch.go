package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CostBatchRepository interface {
	Create(ctx context.Context, billingAccountID string, merkelRoot string, totalItems int, totalCost float64, createdAt time.Time) (*db.CostBatch, error)
	GetByID(ctx context.Context, batchID string) (*db.CostBatch, error)
	ListByBillingProvider(ctx context.Context, billingProviderID string, startTime, endTime time.Time) ([]db.CostBatch, error)
}

type costBatchRepo struct {
	db *gorm.DB
}

func newCostBatchRepo(database *gorm.DB) CostBatchRepository {
	return &costBatchRepo{db: database}
}

func (r *costBatchRepo) Create(ctx context.Context, billingAccountID string, merkelRoot string, totalItems int, totalCost float64, createdAt time.Time) (*db.CostBatch, error) {
	batch := &db.CostBatch{
		ID:               uuid.NewString(),
		BillingAccountID: billingAccountID,
		CreatedAt:        createdAt,
		MerkelRoot:       merkelRoot,
		TotalItems:       totalItems,
		TotalCost:        totalCost,
	}
	return batch, r.db.WithContext(ctx).Create(batch).Error
}

func (r *costBatchRepo) GetByID(ctx context.Context, batchID string) (*db.CostBatch, error) {
	var batch db.CostBatch
	return &batch, r.db.WithContext(ctx).First(&batch, "id = ?", batchID).Error
}

func (r *costBatchRepo) ListByBillingProvider(ctx context.Context, billingProviderID string, startTime, endTime time.Time) ([]db.CostBatch, error) {
	var batches []db.CostBatch

	query := r.db.WithContext(ctx).
		Joins("JOIN billing_accounts ON cost_batches.billing_account_id = billing_accounts.id").
		Where("billing_accounts.billing_provider_id = ?", billingProviderID)

	// Add date range if provided
	if !startTime.IsZero() && !endTime.IsZero() {
		query = query.Where("cost_batches.created_at BETWEEN ? AND ?", startTime, endTime)
	}

	return batches, query.Find(&batches).Error
}
