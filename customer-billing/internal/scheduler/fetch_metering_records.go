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

// Fetches the billing provider's report of charge batches for a billing account, then
// independently fetches each batch (including raw FOCUS line items) directly from the CSP
// that produced it, allowing comparison between the two.
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
	links, err := j.repos.CloudServiceProviderAccount.ListByBillingAccount(ctx, account.ID)
	if err != nil {
		return err
	}
	linksByCSP := make(map[string]repository.CloudServiceProviderAccountLinkWithTotalCost, len(links))
	for _, link := range links {
		linksByCSP[link.CloudProviderID] = link
	}

	latestBatch, err := j.repos.CloudServiceProviderChargeBatch.GetLatestBatchForBillingAccount(ctx, account.ID)
	if err != nil {
		return err
	}

	fromTime := time.Time{}
	if latestBatch != nil {
		fromTime = latestBatch.CreatedAt
	}

	bpClient := clients.NewBillingProviderClient(account.BillingProviderBaseURL)
	batches, err := bpClient.GetBillingAccountRecords(account.ID, fromTime, j.clock.Now())
	if err != nil {
		return err
	}

	storedBatches := 0
	for _, batch := range batches {
		if latestBatch != nil && batch.BatchID == latestBatch.ID {
			continue
		}

		link, ok := linksByCSP[batch.CloudServiceProviderID]
		if !ok || link.CustomerAPIEndpointURL == "" {
			log.Printf("Skipping CSP verification for batch %s: no customer API endpoint on file for CSP %s", batch.BatchID, batch.CloudServiceProviderID)
			continue
		}

		cspClient := clients.NewCloudServiceProviderClient(link.CustomerAPIEndpointURL)
		detail, err := cspClient.GetChargeBatch(batch.BatchID, account.ID)
		if err != nil {
			log.Printf("Error fetching charge batch %s from CSP %s: %v", batch.BatchID, batch.CloudServiceProviderID, err)
			continue
		}

		_, err = j.repos.CloudServiceProviderChargeBatch.Create(ctx, repository.CreateCloudServiceProviderChargeBatchParams{
			ID:                     detail.BatchID,
			BillingAccountID:       account.ID,
			CloudServiceProviderID: detail.CloudServiceProviderID,
			TotalItems:             int32(detail.LineItemCount),
			TotalCost:              detail.TotalBilledCost,
			BilledCurrency:         detail.BilledCurrency,
			MerkleRoot:             detail.MerkleRoot,
			BatchSignature:         detail.BatchSignature,
			CreatedAt:              detail.CreatedAt,
			ReceivedAt:             j.clock.Now(),
		})
		if err != nil {
			log.Printf("Error saving CSP charge batch %s for billing account %s: %v", detail.BatchID, account.ID, err)
			continue
		}

		if len(detail.LineItems) > 0 {
			if err := j.repos.CloudServiceProviderFocusRecord.CreateBatch(ctx, detail.BatchID, account.ID, detail.LineItems); err != nil {
				log.Printf("Error saving line items for CSP charge batch %s: %v", detail.BatchID, err)
			}
		}

		storedBatches++
	}

	log.Printf("Fetched %d charge batches from billing provider for billing account %s, verified and stored %d directly with their CSPs", len(batches), account.ID, storedBatches)

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