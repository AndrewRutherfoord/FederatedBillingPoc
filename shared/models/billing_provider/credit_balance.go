package billingprovider

type CreditBalanceResponse struct {
	BillingAccountID   string  `json:"billing_account_id"`
	CreditAvailable    float64 `json:"credit_available"`     // Total credit available before subtracting usage
	CreditUsed         float64 `json:"credit_used"`          // Total credit used so far
	CreditCurrency     string  `json:"credit_currency"`      // Currency of the credit (e.g., EUR)
	PaymentModel       string  `json:"payment_model"`        // Payment model (e.g., prepaid, postpaid)
	BillingPeriodStart string  `json:"billing_period_start"` // Start of the billing period (ISO 8601 format)
	BillingPeriodEnd   string  `json:"billing_period_end"`   // End of the billing period (ISO 8601 format)
}
