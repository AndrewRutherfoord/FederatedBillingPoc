package models

type BillingRecord struct {
	BillingProviderID  string `json:"billing_provider_id"`
	ResourceProviderID string `json:"resource_provider_id"`
	BillingAccountID   string `json:"billing_account_id"`
}

type AggregatedChargeRecord struct {
	BillingRecord // embed common billing record fields

	BatchID         string  `json:"batch_id"`
	TotalBilledCost float64 `json:"total_billed_cost"`
	BilledCurrency  string  `json:"billed_currency"`
	LineItemCount   int     `json:"line_item_count"`
	BatchHash       string  `json:"batch_hash"`
	BatchSignature  string  `json:"batch_signature"`
}

type AggregatedChargeMeteringRecords struct {
	AggregatedChargeRecord // embed common aggregated charge record fields

	LineItems []FocusLineItem `json:"line_items"`
}
