package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
)

type ChargeBatchRepository interface {
	Create(ctx context.Context, billingAccountID string, billingProviderID string, cloudServiceProviderID string, merkleRoot string, totalItems int, totalCost float64, createdAt time.Time) (*db.ChargeBatch, error)
	GetByID(ctx context.Context, batchID string) (*db.ChargeBatch, error)
	ListByBillingProvider(ctx context.Context, billingProviderID string, startTime, endTime time.Time) ([]db.ChargeBatch, error)
}

type chargeBatchRepo struct {
	db *db.DB
}

func newChargeBatchRepo(database *db.DB) ChargeBatchRepository {
	return &chargeBatchRepo{db: database}
}

func (r *chargeBatchRepo) Create(ctx context.Context, billingAccountID string, billingProviderID string, cloudServiceProviderID string, merkleRoot string, totalItems int, totalCost float64, createdAt time.Time) (*db.ChargeBatch, error) {
	batch := &db.ChargeBatch{
		ID:                     uuid.NewString(),
		BillingAccountID:       billingAccountID,
		BillingProviderID:      billingProviderID,
		CloudServiceProviderID: cloudServiceProviderID,
		CreatedAt:              createdAt,
		MerkleRoot:             merkleRoot,
		TotalItems:             totalItems,
		TotalCost:              totalCost,
	}
	return batch, r.db.WithContext(ctx).Create(batch).Error
}

func (r *chargeBatchRepo) GetByID(ctx context.Context, batchID string) (*db.ChargeBatch, error) {
	var batch db.ChargeBatch
	return &batch, r.db.WithContext(ctx).First(&batch, "id = ?", batchID).Error
}

func (r *chargeBatchRepo) ListByBillingProvider(ctx context.Context, billingProviderID string, startTime, endTime time.Time) ([]db.ChargeBatch, error) {
	var batches []db.ChargeBatch

	query := r.db.WithContext(ctx).
		Joins("JOIN billing_accounts ON charge_batches.billing_account_id = billing_accounts.id").
		Where("billing_accounts.billing_provider_id = ?", billingProviderID)

	// Add date range if provided
	if !startTime.IsZero() && !endTime.IsZero() {
		query = query.Where("charge_batches.created_at BETWEEN ? AND ?", startTime, endTime)
	}

	return batches, query.Find(&batches).Error
}
