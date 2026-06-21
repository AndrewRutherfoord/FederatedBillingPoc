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

// Fetches CSP's report of a charge batches (including raw FOCUS line items)
// Allows comparison to those provided by the billing provider
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

func (j *FetchCloudServiceProviderMeteringRecordsJob) fetchCloudServiceProviderRecords(ctx context.Context, billingAccountID string, link repository.CloudServiceProviderAccountLinkWithTotalCost) error {
	if link.CustomerAPIEndpointURL == "" {
		log.Printf("Skipping CSP %s for billing account %s: no customer API endpoint on file", link.CloudProviderID, billingAccountID)
		return nil
	}

	cspClient := clients.NewCloudServiceProviderClient(link.CustomerAPIEndpointURL)

	latestBatch, err := j.repos.CloudServiceProviderChargeBatch.GetLatestBatchForBillingAccount(ctx, billingAccountID)
	if err != nil {
		return err
	}

	fromTime := time.Time{}
	if latestBatch != nil {
		fromTime = latestBatch.CreatedAt
	}

	batches, err := cspClient.GetBillingAccountRecords(billingAccountID, fromTime, j.clock.Now())
	if err != nil {
		return err
	}

	storedBatches := 0
	for _, batch := range batches {
		if latestBatch != nil && batch.BatchID == latestBatch.ID {
			continue
		}

		_, err := j.repos.CloudServiceProviderChargeBatch.Create(ctx, repository.CreateCloudServiceProviderChargeBatchParams{
			ID:                     batch.BatchID,
			BillingAccountID:       billingAccountID,
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
			log.Printf("Error saving CSP charge batch %s for billing account %s: %v", batch.BatchID, billingAccountID, err)
			continue
		}

		if len(batch.LineItems) > 0 {
			if err := j.repos.CloudServiceProviderFocusRecord.CreateBatch(ctx, batch.BatchID, billingAccountID, batch.LineItems); err != nil {
				log.Printf("Error saving line items for CSP charge batch %s: %v", batch.BatchID, err)
			}
		}

		storedBatches++
	}

	log.Printf("Fetched %d charge batches from CSP %s for billing account %s, stored %d new batches", len(batches), link.CloudProviderID, billingAccountID, storedBatches)

	return nil
}

func (j *FetchCloudServiceProviderMeteringRecordsJob) fetchBillingAccountRecords(ctx context.Context, account *repository.BillingAccountWithProviderName) error {
	links, err := j.repos.CloudServiceProviderAccount.ListByBillingAccount(ctx, account.ID)
	if err != nil {
		return err
	}

	for _, link := range links {
		if err := j.fetchCloudServiceProviderRecords(ctx, account.ID, link); err != nil {
			log.Printf("Error fetching CSP records for billing account %s, CSP %s: %v", account.ID, link.CloudProviderID, err)
			continue
		}
	}

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
