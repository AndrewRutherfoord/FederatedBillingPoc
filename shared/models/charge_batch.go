package models

import "time"

// BillingContext identifies the parties and billing account a charge batch belongs to.
type BillingContext struct {
	BillingProviderID      string `json:"billing_provider_id"`
	CloudServiceProviderID string `json:"cloud_service_provider_id"`
	BillingAccountID       string `json:"billing_account_id"`
}

// ChargeBatch is the metadata for a batch of FocusLineItems that a CSP hands off to a BP.
type ChargeBatch struct {
	BillingContext

	BatchID         string    `json:"batch_id"`
	TotalBilledCost float64   `json:"total_billed_cost"`
	BilledCurrency  string    `json:"billed_currency"`
	LineItemCount   int       `json:"line_item_count"`
	MerkleRoot      string    `json:"merkle_root"`
	BatchSignature  string    `json:"batch_signature"`
	CreatedAt       time.Time `json:"created_at"`
	BillingPeriodID *string   `json:"billing_period_id"` // Nil until a billing period job assigns this batch to a period
}

// ChargeBatchDetail is a ChargeBatch plus the underlying FocusLineItems it aggregates.
type ChargeBatchDetail struct {
	ChargeBatch

	LineItems []FocusLineItem `json:"line_items"`
}
