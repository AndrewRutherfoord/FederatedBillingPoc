package billingprovider

type RegisterBillingAccountRequest struct {
	ReturnURL string `json:"return_url"`
}

type RegisterBillingAccountResponse struct {
	ID          string `json:"id"`
	RedirectURL string `json:"redirect_url"`
}
