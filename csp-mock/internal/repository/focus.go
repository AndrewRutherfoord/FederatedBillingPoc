package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/google/uuid"
)

type FocusFilter struct {
	BillingAccountID string
	BatchID          string
	From             time.Time
	To               time.Time
}

type FocusRepository interface {
	Insert(ctx context.Context, item sharedmodels.FocusLineItem) error
	InsertBatch(ctx context.Context, items []sharedmodels.FocusLineItem) error
	InsertWithBatchID(ctx context.Context, item sharedmodels.FocusLineItem, batchID string) error
	InsertBatchWithBatchID(ctx context.Context, items []sharedmodels.FocusLineItem, batchID string) error
	List(ctx context.Context, filter FocusFilter) ([]sharedmodels.FocusLineItem, error)
}

type focusRepo struct {
	db *db.DB
}

func newFocusRepo(database *db.DB) FocusRepository {
	return &focusRepo{db: database}
}

func (r *focusRepo) Insert(ctx context.Context, item sharedmodels.FocusLineItem) error {
	return r.db.WithContext(ctx).Create(&db.FocusRecord{ID: uuid.NewString(), FocusLineItem: item}).Error
}

func (r *focusRepo) InsertBatch(ctx context.Context, items []sharedmodels.FocusLineItem) error {
	records := make([]db.FocusRecord, len(items))
	for i, item := range items {
		records[i] = db.FocusRecord{ID: uuid.NewString(), FocusLineItem: item}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *focusRepo) InsertWithBatchID(ctx context.Context, item sharedmodels.FocusLineItem, batchID string) error {
	return r.db.WithContext(ctx).Create(&db.FocusRecord{ID: uuid.NewString(), FocusLineItem: item, BatchID: batchID}).Error
}

func (r *focusRepo) InsertBatchWithBatchID(ctx context.Context, items []sharedmodels.FocusLineItem, batchID string) error {
	records := make([]db.FocusRecord, len(items))
	for i, item := range items {
		records[i] = db.FocusRecord{ID: uuid.NewString(), FocusLineItem: item, BatchID: batchID}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *focusRepo) List(ctx context.Context, filter FocusFilter) ([]sharedmodels.FocusLineItem, error) {
	q := r.db.WithContext(ctx).Model(&db.FocusRecord{})
	if filter.BillingAccountID != "" {
		q = q.Where("billing_account_id = ?", filter.BillingAccountID)
	}
	if filter.BatchID != "" {
		q = q.Where("batch_id = ?", filter.BatchID)
	}
	if !filter.From.IsZero() {
		q = q.Where("charge_period_start >= ?", filter.From)
	}
	if !filter.To.IsZero() {
		q = q.Where("charge_period_end <= ?", filter.To)
	}

	var records []db.FocusRecord
	if err := q.Find(&records).Error; err != nil {
		return nil, err
	}

	items := make([]sharedmodels.FocusLineItem, len(records))
	for i, rec := range records {
		items[i] = rec.FocusLineItem
	}
	return items, nil
}
