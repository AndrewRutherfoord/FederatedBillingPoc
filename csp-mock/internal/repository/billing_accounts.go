package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"gorm.io/gorm"
)

type BillingAccountRepository interface {
	Create(ctx context.Context, account db.BillingAccount) (*db.BillingAccount, error)
	GetByAccountID(ctx context.Context, accountID string) (*db.BillingAccount, error)
	ListByCustomer(ctx context.Context, customerID string) ([]db.BillingAccount, error)
	ListByBillingProvider(ctx context.Context, billingProviderID string) ([]db.BillingAccount, error)
}

type billingAccountRepo struct {
	db *gorm.DB
}

func newBillingAccountRepo(database *gorm.DB) BillingAccountRepository {
	return &billingAccountRepo{db: database}
}

func (r *billingAccountRepo) GetByAccountID(ctx context.Context, accountID string) (*db.BillingAccount, error) {
	var account db.BillingAccount
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *billingAccountRepo) ListByCustomer(ctx context.Context, customerID string) ([]db.BillingAccount, error) {
	var accounts []db.BillingAccount
	if err := r.db.WithContext(ctx).Where("customer_id = ?", customerID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *billingAccountRepo) ListByBillingProvider(ctx context.Context, billingProviderID string) ([]db.BillingAccount, error) {
	var accounts []db.BillingAccount
	if err := r.db.WithContext(ctx).Where("billing_provider_id = ?", billingProviderID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
func (r *billingAccountRepo) Create(ctx context.Context, account db.BillingAccount) (*db.BillingAccount, error) {
	if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *billingAccountRepo) Delete(ctx context.Context, account *db.BillingAccount) (*db.BillingAccount, error) {
	if err := r.db.WithContext(ctx).Delete(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}
