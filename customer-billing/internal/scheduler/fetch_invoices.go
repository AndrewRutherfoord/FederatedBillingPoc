package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	billingprovidermodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
)

const invoiceBatchPageSize = 100

type FetchInvoicesJob struct {
	id    string
	repos *repository.Repos
	clock shared.Clock
}

func NewFetchInvoicesJob(id string, repos *repository.Repos, clock shared.Clock) *FetchInvoicesJob {
	return &FetchInvoicesJob{id: id, repos: repos, clock: clock}
}

func (j *FetchInvoicesJob) ID() string {
	return j.id
}

func (j *FetchInvoicesJob) syncInvoiceBatches(ctx context.Context, bpClient clients.BillingProviderClient, account *repository.BillingAccountWithProviderName, invoice billingprovidermodels.Invoice) (map[string][]sharedmodels.ChargeBatch, error) {
	batchesByProvider := map[string][]sharedmodels.ChargeBatch{}

	for offset := 0; ; offset += invoiceBatchPageSize {
		batches, err := bpClient.GetBillingPeriodBatches(account.ID, invoice.BillingPeriodID, invoiceBatchPageSize, offset)
		if err != nil {
			return nil, err
		}
		if len(batches) == 0 {
			return batchesByProvider, nil
		}

		for _, batch := range batches {
			batchesByProvider[batch.CloudServiceProviderID] = append(batchesByProvider[batch.CloudServiceProviderID], batch)

			_, err := j.repos.BillingAccountChargeBatch.UpsertInvoiced(ctx, repository.CreateBillingAccountChargeBatchParams{
				ID:                     batch.BatchID,
				BillingAccountID:       account.ID,
				BillingPeriodID:        invoice.BillingPeriodID,
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
				log.Printf("FetchInvoicesJob: error upserting invoiced batch %s: %v", batch.BatchID, err)
			}
		}
	}
}

// verifyProviderLineItemMerkleRoot rebuilds the merkle tree proving totals are backed by batches
func verifyProviderLineItemMerkleRoot(batches []sharedmodels.ChargeBatch, claimedRoot string) bool {
	roots := make([]string, len(batches))
	for i, b := range batches {
		roots[i] = b.MerkleRoot
	}
	return shared.MerkleRootFromHashes(roots) == claimedRoot
}

func (j *FetchInvoicesJob) syncAccountInvoices(ctx context.Context, account *repository.BillingAccountWithProviderName) error {
	bpClient := clients.NewBillingProviderClient(account.BillingProviderBaseURL)

	invoices, err := bpClient.GetInvoices(account.ID)
	if err != nil {
		return err
	}

	for _, inv := range invoices {
		exists, err := j.repos.Invoice.Exists(ctx, inv.ID)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		batchesByProvider, err := j.syncInvoiceBatches(ctx, bpClient, account, inv)
		if err != nil {
			log.Printf("FetchInvoicesJob: error syncing batches for invoice %s, will retry next run: %v", inv.ID, err)
			continue
		}

		if err := j.repos.Invoice.Create(ctx, repository.CreateInvoiceParams{
			ID:               inv.ID,
			BillingAccountID: account.ID,
			BillingPeriodID:  inv.BillingPeriodID,
			Amount:           inv.Amount,
			Currency:         inv.Currency,
			Status:           inv.Status,
			IssuedAt:         inv.IssuedAt,
			DueAt:            inv.DueAt,
		}); err != nil {
			log.Printf("FetchInvoicesJob: error storing invoice %s: %v", inv.ID, err)
			continue
		}

		for _, li := range inv.ProviderLineItems {
			valid := verifyProviderLineItemMerkleRoot(batchesByProvider[li.CloudServiceProviderID], li.MerkleRoot)
			if !valid {
				log.Printf("FetchInvoicesJob: MERKLE ROOT MISMATCH for invoice %s, provider %s: line item claims root %s over %d batch(es)", inv.ID, li.CloudServiceProviderID, li.MerkleRoot, li.BatchCount)
			}

			if err := j.repos.Invoice.CreateProviderLineItem(ctx, repository.CreateInvoiceProviderLineItemParams{
				ID:                     inv.ID + ":" + li.CloudServiceProviderID,
				InvoiceID:              inv.ID,
				CloudServiceProviderID: li.CloudServiceProviderID,
				Amount:                 li.Amount,
				MerkleRoot:             li.MerkleRoot,
				BatchCount:             int64(li.BatchCount),
				MerkleValid:            valid,
			}); err != nil {
				log.Printf("FetchInvoicesJob: error storing line item for invoice %s: %v", inv.ID, err)
			}
		}

		log.Printf("FetchInvoicesJob: synced invoice %s for billing account %s", inv.ID, account.ID)
	}

	return nil
}

func (j *FetchInvoicesJob) Execute(ctx context.Context, startTime time.Time, lastExecution time.Time) error {
	log.Printf("FetchInvoicesJob: Executing job %s", j.id)

	billingAccounts, err := j.repos.BillingAccount.ListBillingAccounts(ctx)
	if err != nil {
		log.Printf("FetchInvoicesJob: Error fetching billing accounts: %v", err)
		return err
	}

	for _, account := range billingAccounts {
		if err := j.syncAccountInvoices(ctx, account); err != nil {
			log.Printf("FetchInvoicesJob: Error syncing invoices for account %s: %v", account.ID, err)
		}
	}

	return nil
}
