package repository

import (
	"context"
	"sort"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const invoiceDueInDays = 14

type InvoiceRepository interface {
	// CreateForPeriod issues an invoice covering the given batches and links each batch to it.
	CreateForPeriod(ctx context.Context, billingAccountID, billingPeriodID, currency string, batches []db.ChargeBatch, issuedAt time.Time) (*db.Invoice, error)
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
			sort.Strings(roots)

			lineItem := &db.InvoiceProviderLineItem{
				ID:                     uuid.New().String(),
				InvoiceID:              invoice.ID,
				CloudServiceProviderID: providerID,
				Amount:                 providerAmount,
				MerkleRoot:             shared.NewMerkleTree(roots).Root(),
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
