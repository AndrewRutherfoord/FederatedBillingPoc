package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const invoiceDueInDays = 14

type InvoiceWithLineItems struct {
	Invoice   db.Invoice
	LineItems []db.InvoiceProviderLineItem
}

type InvoiceRepository interface {
	// CreateForPeriod issues an invoice covering the given batches and links each batch to it.
	CreateForPeriod(ctx context.Context, billingAccountID, billingPeriodID, currency string, batches []db.ChargeBatch, issuedAt time.Time) (*db.Invoice, error)
	ListByBillingAccount(ctx context.Context, billingAccountID string) ([]InvoiceWithLineItems, error)
}

type invoiceRepo struct {
	db *db.DB
}

func newInvoiceRepo(database *db.DB) InvoiceRepository {
	return &invoiceRepo{db: database}
}

func (r *invoiceRepo) CreateForPeriod(ctx context.Context, billingAccountID, billingPeriodID, currency string, batches []db.ChargeBatch, issuedAt time.Time) (*db.Invoice, error) {
	var total float64
	batchesByProvider := map[string][]db.ChargeBatch{}
	for _, b := range batches {
		total += b.TotalCost
		batchesByProvider[b.CloudServiceProviderID] = append(batchesByProvider[b.CloudServiceProviderID], b)
	}

	invoice := &db.Invoice{
		ID:               uuid.New().String(),
		BillingAccountID: billingAccountID,
		BillingPeriodID:  billingPeriodID,
		Amount:           total,
		Currency:         currency,
		Status:           "issued",
		IssuedAt:         issuedAt,
		DueAt:            issuedAt.AddDate(0, 0, invoiceDueInDays),
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(invoice).Error; err != nil {
			return err
		}
		for _, b := range batches {
			invoiceBatch := &db.InvoiceBatch{
				ID:                     uuid.New().String(),
				InvoiceID:              invoice.ID,
				Amount:                 b.TotalCost,
				CloudServiceProviderID: b.CloudServiceProviderID,
				ChargeBatchID:          b.ID,
			}
			if err := tx.Create(invoiceBatch).Error; err != nil {
				return err
			}
		}

		for providerID, providerBatches := range batchesByProvider {
			roots := make([]string, len(providerBatches))
			var providerAmount float64
			for i, b := range providerBatches {
				roots[i] = b.MerkleRoot
				providerAmount += b.TotalCost
			}

			lineItem := &db.InvoiceProviderLineItem{
				ID:                     uuid.New().String(),
				InvoiceID:              invoice.ID,
				CloudServiceProviderID: providerID,
				Amount:                 providerAmount,
				MerkleRoot:             shared.MerkleRootFromHashes(roots),
				BatchCount:             len(providerBatches),
			}
			if err := tx.Create(lineItem).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *invoiceRepo) ListByBillingAccount(ctx context.Context, billingAccountID string) ([]InvoiceWithLineItems, error) {
	var invoices []db.Invoice
	if err := r.db.WithContext(ctx).
		Where("billing_account_id = ?", billingAccountID).
		Order("issued_at DESC").
		Find(&invoices).Error; err != nil {
		return nil, err
	}
	if len(invoices) == 0 {
		return nil, nil
	}

	invoiceIDs := make([]string, len(invoices))
	for i, inv := range invoices {
		invoiceIDs[i] = inv.ID
	}

	var lineItems []db.InvoiceProviderLineItem
	if err := r.db.WithContext(ctx).Where("invoice_id IN ?", invoiceIDs).Find(&lineItems).Error; err != nil {
		return nil, err
	}

	lineItemsByInvoice := map[string][]db.InvoiceProviderLineItem{}
	for _, li := range lineItems {
		lineItemsByInvoice[li.InvoiceID] = append(lineItemsByInvoice[li.InvoiceID], li)
	}

	result := make([]InvoiceWithLineItems, len(invoices))
	for i, inv := range invoices {
		result[i] = InvoiceWithLineItems{Invoice: inv, LineItems: lineItemsByInvoice[inv.ID]}
	}
	return result, nil
}
