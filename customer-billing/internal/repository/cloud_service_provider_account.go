package repository

import (
	"context"
	"errors"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
)

type CloudServiceProviderAccountRepository interface {
	Create(ctx context.Context, accountID string, cloudServiceProviderID string, billingAccountID string) (CloudServiceProvider, error)
}

type cloudServiceProviderAccountRepository struct {
	db *db.DB
}

func newCloudServiceProviderAccountRepo(database *db.DB) CloudServiceProviderAccountRepository {
	return &cloudServiceProviderAccountRepository{db: database}
}

func (r *cloudServiceProviderAccountRepository) Create(ctx context.Context, accountID string, cloudServiceProviderID string, billingAccountID string) (CloudServiceProvider, error) {
	return CloudServiceProvider{}, errors.New("not implemented")
}
