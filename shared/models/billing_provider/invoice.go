package billingprovider

import "time"

type InvoiceProviderLineItem struct {
	CloudServiceProviderID string  `json:"cloud_service_provider_id"`
	Amount                 float64 `json:"amount"`
	MerkleRoot             string  `json:"merkle_root"`
	BatchCount             int     `json:"batch_count"`
}

type Invoice struct {
	ID                string                    `json:"id"`
	BillingAccountID  string                    `json:"billing_account_id"`
	BillingPeriodID   string                    `json:"billing_period_id"`
	Amount            float64                   `json:"amount"`
	Currency          string                    `json:"currency"`
	Status            string                    `json:"status"`
	IssuedAt          time.Time                 `json:"issued_at"`
	DueAt             time.Time                 `json:"due_at"`
	ProviderLineItems []InvoiceProviderLineItem `json:"provider_line_items"`
}

type GetInvoicesResponse struct {
	Invoices []Invoice `json:"invoices"`
	Count    int       `json:"count"`
}
