package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/services"
	"github.com/gin-gonic/gin"
)

type RegisterAcccountRequest struct {
	AccountAlias           string `json:"account_alias" binding:"required"`
	BillingProviderBaseURL string `json:"billing_provider_base_url" binding:"required"`
	ReturnURL              string `json:"return_url" binding:"required"`
}

type RegisterAccountResponse struct {
	AccountID   string `json:"account_id"`
	RedirectURL string `json:"redirect_url"`
}

// RegisterBillingAccount godoc
//
//	@Summary		Register billing account with billing provider
//	@Description	Registers a new billing account with a billing provider and returns a redirect URL for the customer onboarding flow.
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Param			request	body		RegisterAcccountRequest		true	"Request body"
//	@Success		201		{object}	RegisterAccountResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/billing/accounts/register [post]
func (s *Server) RegisterBillingAccount(c *gin.Context) {
	var req RegisterAcccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bp, err := services.SyncBillingProviderMetadata(s.repos, req.BillingProviderBaseURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync billing provider metadata"})
		return
	}

	client := clients.NewBillingProviderClient(req.BillingProviderBaseURL)

	// Register the billing account with the billing provider.
	// Name and email are collected later via the BP's onboarding form.
	bpResponse, err := client.RegisterBillingAccount(req.ReturnURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the pending account locally
	if _, err := s.repos.BillingAccount.CreateBillingAccount(bpResponse.ID, bp.ID, req.AccountAlias, ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save billing account"})
		return
	}

	c.JSON(http.StatusCreated, RegisterAccountResponse{
		AccountID:   bpResponse.ID,
		RedirectURL: bpResponse.RedirectURL,
	})
}

type billingAccountResponse struct {
	ID                  string `json:"id"`
	Alias               string `json:"alias"`
	BillingProviderID   string `json:"billing_provider_id"`
	BillingProviderName string `json:"billing_provider_name"`
}

// ListBillingAccounts godoc
//
//	@Summary		Get all billing accounts
//	@Tags			billing
//	@Produce		json
//	@Success		200	{array}		[]billingAccountResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts [get]
func (s *Server) ListBillingAccounts(c *gin.Context) {
	accounts, err := s.repos.BillingAccount.ListBillingAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get billing accounts"})
		return
	}

	accountsResponse := make([]*billingAccountResponse, len(accounts))
	for i := range accounts {
		accountsResponse[i] = &billingAccountResponse{
			ID:                  accounts[i].ID,
			Alias:               accounts[i].Alias,
			BillingProviderID:   accounts[i].BillingProviderID,
			BillingProviderName: accounts[i].BillingProviderName,
		}
	}
	c.JSON(http.StatusOK, accountsResponse)
}

type billingAccountDetailResponse struct {
	ID                      string                        `json:"id"`
	Alias                   string                        `json:"alias"`
	BillingProviderID       string                        `json:"billing_provider_id"`
	BillingProviderName     string                        `json:"billing_provider_name"`
	SupportedCloudProviders []supportedCloudProviderEntry `json:"supported_cloud_providers"`
}

type supportedCloudProviderEntry struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	APIEndpointURL string `json:"api_endpoint_url"`
}

// GetBillingAccount godoc
//
//	@Summary		Get a billing account by ID
//	@Tags			billing
//	@Produce		json
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{object}	billingAccountDetailResponse
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id} [get]
func (s *Server) GetBillingAccount(c *gin.Context) {
	account, err := s.repos.BillingAccount.GetBillingAccountByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	csps := make([]supportedCloudProviderEntry, len(account.BillingProvider.SupportedCloudProviders))
	for i, csp := range account.BillingProvider.SupportedCloudProviders {
		csps[i] = supportedCloudProviderEntry{ID: csp.ID, Name: csp.Name, APIEndpointURL: csp.APIEndpointURL}
	}

	c.JSON(http.StatusOK, billingAccountDetailResponse{
		ID:                      account.ID,
		Alias:                   account.Alias,
		BillingProviderID:       account.BillingProvider.ID,
		BillingProviderName:     account.BillingProvider.Name,
		SupportedCloudProviders: csps,
	})
}
