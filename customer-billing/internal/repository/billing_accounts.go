package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type BillingAccountWithProviderName struct {
	ID                  string
	BillingProviderID   string
	Alias               string
	BillingProviderName string
}

type BillingAccountWithProvider struct {
	ID              string
	BillingProvider BillingProvider
	Alias           string
}

type BillingAccountRepository interface {
	ListBillingAccounts(ctx context.Context) ([]*BillingAccountWithProviderName, error)
	GetBillingAccountByID(ctx context.Context, id string) (*BillingAccountWithProvider, error)
	CreateBillingAccount(ctx context.Context, id string, billingProviderID string, accountAlias string) error
}

type billingAccountRepo struct {
	db *db.DB
}

func newBillingAccountRepo(database *db.DB) BillingAccountRepository {
	return &billingAccountRepo{db: database}
}

func (r *billingAccountRepo) ListBillingAccounts(ctx context.Context) ([]*BillingAccountWithProviderName, error) {
	rows, err := r.db.Queries.ListBillingAccounts(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*BillingAccountWithProviderName, len(rows))
	for i, row := range rows {
		result[i] = &BillingAccountWithProviderName{
			ID:                  row.ID,
			BillingProviderID:   row.BillingProviderID,
			Alias:               row.Alias,
			BillingProviderName: row.BillingProviderName.String,
		}
	}
	return result, nil
}

func (r *billingAccountRepo) GetBillingAccountByID(ctx context.Context, id string) (*BillingAccountWithProvider, error) {
	account, err := r.db.Queries.GetBillingAccount(ctx, id)
	if err != nil {
		return nil, err
	}

	bpRow, err := r.db.Queries.GetBillingProvider(ctx, account.BillingProviderID)
	if err != nil {
		return nil, err
	}

	cspRows, err := r.db.Queries.ListBillingProviderSupportedCSPs(ctx, bpRow.ID)
	if err != nil {
		return nil, err
	}

	return &BillingAccountWithProvider{
		ID:              account.ID,
		Alias:           account.Alias,
		BillingProvider: toBillingProvider(bpRow, cspRows),
	}, nil
}

func (r *billingAccountRepo) CreateBillingAccount(ctx context.Context, id string, billingProviderID string, accountAlias string) error {
	return r.db.Queries.CreateBillingAccount(ctx, sqlcdb.CreateBillingAccountParams{
		ID:                id,
		BillingProviderID: billingProviderID,
		Alias:             accountAlias,
		CreatedAt:         time.Now(),
	})
}
