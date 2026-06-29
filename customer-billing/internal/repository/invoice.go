package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
)

type InvoiceProviderLineItem struct {
	CloudServiceProviderID string
	Amount                 float64
	MerkleRoot             string
	BatchCount             int64
}

type Invoice struct {
	ID                string
	BillingAccountID  string
	BillingPeriodID   string
	Amount            float64
	Currency          string
	Status            string
	IssuedAt          time.Time
	DueAt             time.Time
	ProviderLineItems []InvoiceProviderLineItem
}

type CreateInvoiceParams struct {
	ID               string
	BillingAccountID string
	BillingPeriodID  string
	Amount           float64
	Currency         string
	Status           string
	IssuedAt         time.Time
	DueAt            time.Time
}

type CreateInvoiceProviderLineItemParams struct {
	ID                     string
	InvoiceID              string
	CloudServiceProviderID string
	Amount                 float64
	MerkleRoot             string
	BatchCount             int64
}

type InvoiceRepository interface {
	// Exists reports whether an invoice has already been synced locally.
	Exists(ctx context.Context, id string) (bool, error)
	Create(ctx context.Context, params CreateInvoiceParams) error
	CreateProviderLineItem(ctx context.Context, params CreateInvoiceProviderLineItemParams) error
	ListByBillingAccount(ctx context.Context, billingAccountID string) ([]Invoice, error)
}

type invoiceRepo struct {
	db *db.DB
}

func newInvoiceRepo(database *db.DB) InvoiceRepository {
	return &invoiceRepo{db: database}
}

func (r *invoiceRepo) Exists(ctx context.Context, id string) (bool, error) {
	_, err := r.db.Queries.GetInvoice(ctx, id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *invoiceRepo) Create(ctx context.Context, params CreateInvoiceParams) error {
	_, err := r.db.Queries.CreateInvoice(ctx, sqlcdb.CreateInvoiceParams{
		ID:               params.ID,
		BillingAccountID: params.BillingAccountID,
		BillingPeriodID:  params.BillingPeriodID,
		Amount:           params.Amount,
		Currency:         params.Currency,
		Status:           params.Status,
		IssuedAt:         params.IssuedAt,
		DueAt:            params.DueAt,
	})
	return err
}

func (r *invoiceRepo) CreateProviderLineItem(ctx context.Context, params CreateInvoiceProviderLineItemParams) error {
	_, err := r.db.Queries.CreateInvoiceProviderLineItem(ctx, sqlcdb.CreateInvoiceProviderLineItemParams{
		ID:                     params.ID,
		InvoiceID:              params.InvoiceID,
		CloudServiceProviderID: params.CloudServiceProviderID,
		Amount:                 params.Amount,
		MerkleRoot:             params.MerkleRoot,
		BatchCount:             params.BatchCount,
	})
	return err
}

func (r *invoiceRepo) ListByBillingAccount(ctx context.Context, billingAccountID string) ([]Invoice, error) {
	rows, err := r.db.Queries.ListInvoicesByBillingAccount(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}

	invoiceIDs := make([]string, len(rows))
	for i, row := range rows {
		invoiceIDs[i] = row.ID
	}

	lineItemRows, err := r.db.Queries.ListInvoiceProviderLineItemsByInvoiceIDs(ctx, invoiceIDs)
	if err != nil {
		return nil, err
	}
	lineItemsByInvoice := map[string][]InvoiceProviderLineItem{}
	for _, li := range lineItemRows {
		lineItemsByInvoice[li.InvoiceID] = append(lineItemsByInvoice[li.InvoiceID], InvoiceProviderLineItem{
			CloudServiceProviderID: li.CloudServiceProviderID,
			Amount:                 li.Amount,
			MerkleRoot:             li.MerkleRoot,
			BatchCount:             li.BatchCount,
		})
	}

	result := make([]Invoice, len(rows))
	for i, row := range rows {
		result[i] = Invoice{
			ID:                row.ID,
			BillingAccountID:  row.BillingAccountID,
			BillingPeriodID:   row.BillingPeriodID,
			Amount:            row.Amount,
			Currency:          row.Currency,
			Status:            row.Status,
			IssuedAt:          row.IssuedAt,
			DueAt:             row.DueAt,
			ProviderLineItems: lineItemsByInvoice[row.ID],
		}
	}
	return result, nil
}
