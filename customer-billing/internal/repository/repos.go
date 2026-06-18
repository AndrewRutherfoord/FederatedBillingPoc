package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
)

type Repos struct {
	Customer                    CustomerRepository
	BillingProvider             BillingProviderRepository
	BillingAccount              BillingAccountRepository
	BillingAccountCostBatch     BillingAccountCostBatchRepository
	CloudServiceProvider        CloudServiceProviderRepository
	CloudServiceProviderAccount CloudServiceProviderAccountRepository
}

func New(config *config.Config, database *db.DB) *Repos {
	return &Repos{
		Customer:                    newCustomerRepo(config),
		BillingProvider:             newBillingProviderRepo(database),
		BillingAccount:              newBillingAccountRepo(database),
		BillingAccountCostBatch:     newBillingAccountCostBatchRepo(database),
		CloudServiceProvider:        newCloudServiceProviderRepo(database),
		CloudServiceProviderAccount: newCloudServiceProviderAccountRepo(database),
	}
}
