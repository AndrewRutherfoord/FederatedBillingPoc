package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type BillingAccountWithProviderName struct {
	db.BillingAccount
	BillingProviderName string
}

type BillingAccountWithProvider struct {
	db.BillingAccount
	BillingProvider BillingProvider
}

type BillingAccountRepository interface {
	ListBillingAccounts() ([]*BillingAccountWithProviderName, error)
	GetBillingAccountByID(id string) (*BillingAccountWithProvider, error)
	CreateBillingAccount(id string, billingProviderID string, accountAlias string, token string) (*db.BillingAccount, error)
}

type billingAccountRepo struct {
	db *db.DB
}

func newBillingAccountRepo(database *db.DB) BillingAccountRepository {
	return &billingAccountRepo{db: database}
}

func (r *billingAccountRepo) ListBillingAccounts() ([]*BillingAccountWithProviderName, error) {
	var billingAccounts []*BillingAccountWithProviderName
	if err := r.db.
		Model(&db.BillingAccount{}).
		Select("billing_accounts.*, billing_providers.name AS billing_provider_name").
		Joins("LEFT JOIN billing_providers ON billing_providers.billing_provider_id = billing_accounts.billing_provider_id").
		Find(&billingAccounts).Error; err != nil {
		return nil, err
	}
	return billingAccounts, nil
}

func (r *billingAccountRepo) GetBillingAccountByID(id string) (*BillingAccountWithProvider, error) {
	var account db.BillingAccount
	if err := r.db.First(&account, "account_id = ?", id).Error; err != nil {
		return nil, err
	}

	var bpRow db.BillingProvider
	if err := r.db.First(&bpRow, "billing_provider_id = ?", account.BillingProviderID).Error; err != nil {
		return nil, err
	}

	var cspRows []db.SupportedCloudProvider
	if err := r.db.Where("billing_provider_id = ?", bpRow.ID).Find(&cspRows).Error; err != nil {
		return nil, err
	}
	csps := make([]SupportedCloudProvider, len(cspRows))
	for i, c := range cspRows {
		csps[i] = SupportedCloudProvider{ID: c.ID, Name: c.Name, APIEndpointURL: c.APIEndpointURL}
	}

	return &BillingAccountWithProvider{
		BillingAccount:  account,
		BillingProvider: toRepoBillingProvider(bpRow, csps),
	}, nil
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
