package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ResourceRepository interface {
	ListByCustomerID(ctx context.Context, customerID string) ([]*db.Resource, error)
	ListByBillingAccountID(ctx context.Context, billingAccountID string) ([]*db.Resource, error)
	Create(ctx context.Context, customerID, billingAccountID, resourceType string, storageGB *decimal.Decimal) (*db.Resource, error)
	SetStorageGB(ctx context.Context, id string, gb decimal.Decimal) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*db.Resource, error)
}

type resourceRepo struct {
	db *gorm.DB
}

func newResourceRepo(database *gorm.DB) ResourceRepository {
	return &resourceRepo{db: database}
}

func (r *resourceRepo) ListByCustomerID(ctx context.Context, customerID string) ([]*db.Resource, error) {
	var resources []*db.Resource
	if err := r.db.WithContext(ctx).Where("customer_id = ? AND deleted_at IS NULL", customerID).Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *resourceRepo) ListByBillingAccountID(ctx context.Context, billingAccountID string) ([]*db.Resource, error) {
	var resources []*db.Resource
	if err := r.db.WithContext(ctx).Where("billing_account_id = ? AND deleted_at IS NULL", billingAccountID).Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

func (r *resourceRepo) GetByID(ctx context.Context, id string) (*db.Resource, error) {
	var resource db.Resource
	if err := r.db.WithContext(ctx).First(&resource, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepo) Create(ctx context.Context, customerID, billingAccountID, resourceType string, storageGB *decimal.Decimal) (*db.Resource, error) {
	resource := &db.Resource{
		ID:               uuid.NewString(),
		CustomerID:       customerID,
		BillingAccountID: billingAccountID,
		ResourceType:     resourceType,
		StorageGB:        storageGB,
		StartedAt:        time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Create(resource).Error; err != nil {
		return nil, err
	}
	return resource, nil
}

func (r *resourceRepo) SetStorageGB(ctx context.Context, id string, gb decimal.Decimal) error {
	return r.db.WithContext(ctx).Model(&db.Resource{}).Where("id = ? AND deleted_at IS NULL", id).Update("storage_gb", gb).Error
}

func (r *resourceRepo) Delete(ctx context.Context, id string) error {
	// Soft delete by setting DeletedAt timestamp
	return r.db.WithContext(ctx).Model(&db.Resource{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now().UTC()).Error
}
