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

type FetchBillingProviderRecordsJob struct {
	id     string
	repos  *repository.Repos
	config *config.Config
	clock  shared.Clock
}

func NewFetchBillingProviderRecordsJob(id string, repos *repository.Repos, config *config.Config, clock shared.Clock) *FetchBillingProviderRecordsJob {
	return &FetchBillingProviderRecordsJob{
		id:     id,
		repos:  repos,
		config: config,
		clock:  clock,
	}
}

func (j *FetchBillingProviderRecordsJob) ID() string {
	return j.id
}

func (j *FetchBillingProviderRecordsJob) fetchBillingAccountRecords(ctx context.Context, account *repository.BillingAccountWithProviderName) error {
	bpClient := clients.NewBillingProviderClient(account.BillingProviderBaseURL)

	latestBatch, err := j.repos.BillingAccountChargeBatch.GetLatestBatchForBillingAccount(ctx, account.ID)
	if err != nil {
		return err
	}

	fromTime := time.Time{}
	if latestBatch != nil {
		fromTime = latestBatch.CreatedAt
	}

	batches, err := bpClient.GetBillingAccountRecords(account.ID, fromTime, j.clock.Now())
	if err != nil {
		return err
	}

	storedBatches := 0
	for _, batch := range batches {
		if latestBatch != nil && batch.BatchID == latestBatch.ID {
			continue
		}
		_, err := j.repos.BillingAccountChargeBatch.Create(ctx, repository.CreateBillingAccountChargeBatchParams{
			ID:                     batch.BatchID,
			BillingAccountID:       account.ID,
			BillingPeriodID:        "", // TODO
			CloudServiceProviderID: batch.CloudServiceProviderID,
			TotalItems:             int32(batch.LineItemCount),
			TotalCost:              batch.TotalBilledCost,
			BilledCurrency:         batch.BilledCurrency,
			MerkleRoot:             batch.MerkleRoot,
			BatchSignature:         batch.BatchSignature,
			CreatedAt:              batch.CreatedAt,
			ReceivedAt:             j.clock.Now(),
		})
		if err != nil {
			log.Printf("Error saving billing account charge batch for account %s: %v", account.ID, err)
			continue
		}
		storedBatches++
	}

	log.Printf("Fetched %d charge batches for billing account %s, stored %d new batches", len(batches), account.ID, storedBatches)

	return nil
}

func (j *FetchBillingProviderRecordsJob) Execute(ctx context.Context, startTime time.Time) error {
	log.Printf("FetchBillingProviderRecordsJob: Executing job %s", j.id)

	billingAccounts, err := j.repos.BillingAccount.ListBillingAccounts(ctx)
	if err != nil {
		log.Printf("FetchBillingProviderRecordsJob: Error fetching billing accounts: %v", err)
		return err
	}

	for _, account := range billingAccounts {
		err := j.fetchBillingAccountRecords(ctx, account)
		if err != nil {
			log.Printf("FetchBillingProviderRecordsJob: Error fetching billing account records for account %s: %v", account.ID, err)
			return err
		}
	}

	return nil
}
