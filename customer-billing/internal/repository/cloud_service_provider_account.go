package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type CloudServiceProviderAccountRepository interface {
	Create(accountID string, cloudServiceProviderID string, billingAccountID string) (CloudServiceProvider, error)
}

type cloudServiceProviderAccountRepository struct {
	db *db.DB
}

func newCloudServiceProviderAccountRepo(database *db.DB) CloudServiceProviderAccountRepository {
	return &cloudServiceProviderAccountRepository{db: database}
}
