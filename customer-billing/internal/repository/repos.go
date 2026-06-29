package repository

import (
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
)

type Repos struct {
	Customer                        CustomerRepository
	BillingProvider                 BillingProviderRepository
	BillingAccount                  BillingAccountRepository
	BillingAccountChargeBatch       BillingAccountChargeBatchRepository
	CloudServiceProvider            CloudServiceProviderRepository
	CloudServiceProviderAccount     CloudServiceProviderAccountRepository
	CloudServiceProviderChargeBatch CloudServiceProviderChargeBatchRepository
	CloudServiceProviderFocusRecord CloudServiceProviderFocusRecordRepository
	ChargeBatchReconciliation       ChargeBatchReconciliationRepository
	Invoice                         InvoiceRepository
}

func New(config *config.Config, database *db.DB) *Repos {
	return &Repos{
		Customer:                        newCustomerRepo(config),
		BillingProvider:                 newBillingProviderRepo(database),
		BillingAccount:                  newBillingAccountRepo(database),
		BillingAccountChargeBatch:       newBillingAccountChargeBatchRepo(database),
		CloudServiceProvider:            newCloudServiceProviderRepo(database),
		CloudServiceProviderAccount:     newCloudServiceProviderAccountRepo(database),
		CloudServiceProviderChargeBatch: newCloudServiceProviderChargeBatchRepo(database),
		CloudServiceProviderFocusRecord: newCloudServiceProviderFocusRecordRepo(database),
		ChargeBatchReconciliation:       newChargeBatchReconciliationRepo(database),
		Invoice:                         newInvoiceRepo(database),
	}
}
