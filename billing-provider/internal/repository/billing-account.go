package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/google/uuid"
)

type BillingAccountRepository interface {
	List(ctx context.Context) ([]db.BillingAccount, error)
	Get(ctx context.Context, id string) (*db.BillingAccount, error)
	Create(ctx context.Context) (*db.BillingAccount, error)
	Update(ctx context.Context, id string, name string, email string, billingCycle db.BillingCycle) (*db.BillingAccount, error)
}

type billingAccountRepo struct {
	db *db.DB
}

func newBillingAccountRepo(database *db.DB) BillingAccountRepository {
	return &billingAccountRepo{db: database}
}

func (r *billingAccountRepo) List(ctx context.Context) ([]db.BillingAccount, error) {
	var billingAccounts []db.BillingAccount
	if err := r.db.Find(&billingAccounts).Error; err != nil {
		return nil, err
	}
	return billingAccounts, nil
}

func (r *billingAccountRepo) Get(ctx context.Context, id string) (*db.BillingAccount, error) {
	var account db.BillingAccount
	if err := r.db.Where("account_id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *billingAccountRepo) Create(ctx context.Context) (*db.BillingAccount, error) {
	account := db.BillingAccount{
		ID: uuid.New().String(),
	}
	if err := r.db.Create(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *billingAccountRepo) Update(ctx context.Context, id string, name string, email string, billingCycle db.BillingCycle) (*db.BillingAccount, error) {
	var account db.BillingAccount
	if err := r.db.Where("account_id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}
	account.Name = name
	account.Email = email
	account.BillingCycle = billingCycle
	if err := r.db.Save(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
