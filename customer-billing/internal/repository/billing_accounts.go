package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type BillingAccountWithProvider struct {
	db.BillingAccount
	BillingProviderName string
}

type BillingAccountRepository interface {
	ListBillingAccounts() ([]*BillingAccountWithProvider, error)
	CreateBillingAccount(id string, billingProviderID string, accountAlias string, token string) (*db.BillingAccount, error)
}

type billingAccountRepo struct {
	db *db.DB
}

func newBillingAccountRepo(database *db.DB) BillingAccountRepository {
	return &billingAccountRepo{db: database}
}

func (r *billingAccountRepo) ListBillingAccounts() ([]*BillingAccountWithProvider, error) {
	var billingAccounts []*BillingAccountWithProvider
	if err := r.db.
		Model(&db.BillingAccount{}).
		Select("billing_accounts.*, billing_providers.name AS billing_provider_name").
		Joins("LEFT JOIN billing_providers ON billing_providers.billing_provider_id = billing_accounts.billing_provider_id").
		Find(&billingAccounts).Error; err != nil {
		return nil, err
	}
	return billingAccounts, nil
}

func (r *billingAccountRepo) CreateBillingAccount(id string, billingProviderID string, accountAlias string, token string) (*db.BillingAccount, error) {
	billingAccount := &db.BillingAccount{
		ID:                id,
		BillingProviderID: billingProviderID,
		Alias:             accountAlias,
		Token:             token,
	}
	if err := r.db.Create(billingAccount).Error; err != nil {
		return nil, err
	}
	return billingAccount, nil
}
