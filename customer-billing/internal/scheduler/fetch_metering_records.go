package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

type FetchCloudServiceProviderMeteringRecordsJob struct {
	id     string
	repos  *repository.Repos
	config *config.Config
	clock  shared.Clock
}

func NewFetchCloudServiceProviderMeteringRecordsJob(id string, repos *repository.Repos, config *config.Config, clock shared.Clock) *FetchCloudServiceProviderMeteringRecordsJob {
	return &FetchCloudServiceProviderMeteringRecordsJob{
		id:     id,
		repos:  repos,
		config: config,
		clock:  clock,
	}
}

func (j *FetchCloudServiceProviderMeteringRecordsJob) ID() string {
	return j.id
}

func (j *FetchCloudServiceProviderMeteringRecordsJob) fetchBillingAccountRecords(ctx context.Context, account *repository.BillingAccountWithProviderName) error {
	cspClient := clients.NewCloudServiceProviderClient(account.BillingProviderBaseURL)

	cspClient.

	return nil
}

func (j *FetchCloudServiceProviderMeteringRecordsJob) Execute(ctx context.Context, startTime time.Time) error {
	log.Printf("FetchCloudServiceProviderMeteringRecordsJob: Executing job %s", j.id)

	billingAccounts, err := j.repos.BillingAccount.ListBillingAccounts(ctx)
	if err != nil {
		log.Printf("FetchCloudServiceProviderMeteringRecordsJob: Error fetching billing accounts: %v", err)
		return err
	}

	for _, account := range billingAccounts {
		err := j.fetchBillingAccountRecords(ctx, account)
		if err != nil {
			log.Printf("FetchCloudServiceProviderMeteringRecordsJob: Error fetching billing account records for account %s: %v", account.ID, err)
			return err
		}
	}

	return nil
}
