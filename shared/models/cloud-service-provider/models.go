package cloudserviceprovider

type RegisterLinkedCloudProviderRequest struct {
	BillingProviderID string `json:"billing_provider_id" binding:"required"`
	BillingAccountID  string `json:"billing_account_id" binding:"required"`
	ReturnURL         string `json:"return_url" binding:"required"`
}

type RegisterLinkedCloudProviderResponse struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirect_url"`
}
