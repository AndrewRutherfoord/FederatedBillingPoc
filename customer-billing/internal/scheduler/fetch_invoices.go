package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
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

func (j *FetchInvoicesJob) syncInvoiceBatches(ctx context.Context, bpClient clients.BillingProviderClient, account *repository.BillingAccountWithProviderName, invoice billingprovidermodels.Invoice) error {
	for offset := 0; ; offset += invoiceBatchPageSize {
		batches, err := bpClient.GetBillingPeriodBatches(account.ID, invoice.BillingPeriodID, invoiceBatchPageSize, offset)
		if err != nil {
			return err
		}
		if len(batches) == 0 {
			return nil
		}

		for _, batch := range batches {
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
			if err := j.repos.Invoice.CreateProviderLineItem(ctx, repository.CreateInvoiceProviderLineItemParams{
				ID:                     inv.ID + ":" + li.CloudServiceProviderID,
				InvoiceID:              inv.ID,
				CloudServiceProviderID: li.CloudServiceProviderID,
				Amount:                 li.Amount,
				MerkleRoot:             li.MerkleRoot,
				BatchCount:             int64(li.BatchCount),
			}); err != nil {
				log.Printf("FetchInvoicesJob: error storing line item for invoice %s: %v", inv.ID, err)
			}
		}

		if err := j.syncInvoiceBatches(ctx, bpClient, account, inv); err != nil {
			log.Printf("FetchInvoicesJob: error syncing batches for invoice %s: %v", inv.ID, err)
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
