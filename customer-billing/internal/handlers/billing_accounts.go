package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/services"
	"github.com/gin-gonic/gin"
)

type RegisterAcccountRequest struct {
	AccountAlias           string `json:"account_alias" binding:"required"`
	BillingProviderBaseURL string `json:"billing_provider_base_url" binding:"required"`
	ReturnURL              string `json:"return_url" binding:"required"` // TODO: I'm not sure I like the frontend setting the URL... Might need to change this...
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

	bp, err := services.SyncBillingProviderMetadata(c.Request.Context(), s.repos, req.BillingProviderBaseURL)
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
	if err := s.repos.BillingAccount.CreateBillingAccount(c.Request.Context(), bpResponse.ID, bp.ID, req.AccountAlias); err != nil {
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
	accounts, err := s.repos.BillingAccount.ListBillingAccounts(c.Request.Context())
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
	account, err := s.repos.BillingAccount.GetBillingAccountByID(c.Request.Context(), c.Param("id"))
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

type CloudProviderLink struct {
	ID                string  `json:"id"`
	CloudProviderID   string  `json:"cloud_provider_id"`
	CloudProviderName string  `json:"cloud_provider_name"`
	TotalCost         float64 `json:"total_cost"`
	BillingCurrency   string  `json:"billing_currency"`
}

// ListBillingProviderLinkedCloudProviders godoc
//
//	@Summary		List cloud providers linked to a billing account
//	@Tags			billing
//	@Produce		json
//	@Param			id	path		string	true	"Account ID"
//	@Success		200	{array}	[]CloudProviderLink
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id}/cloud-provider-accounts [get]
func (s *Server) ListBillingProviderLinkedCloudProviders(c *gin.Context) {
	links, err := s.repos.CloudServiceProviderAccount.ListByBillingAccount(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list linked cloud providers"})
		return
	}

	result := make([]CloudProviderLink, len(links))
	for i, l := range links {
		result[i] = CloudProviderLink{
			ID:                l.ID,
			CloudProviderID:   l.CloudProviderID,
			CloudProviderName: l.CloudProviderName,
			TotalCost:         l.TotalCost,
			BillingCurrency:   l.BillingCurrency,
		}
	}
	c.JSON(http.StatusOK, result)
}

type RegisterLinkedCloudProviderRequest struct {
	AccountID       string `json:"account_id" binding:"required"`        // The billing account ID to link the cloud provider to
	CloudProviderID string `json:"cloud_provider_id" binding:"required"` // The identifier from the billing provider metadata
	ReturnURL       string `json:"return_url" binding:"required"`
}

// RegisterCloudProviderAccount godoc
//
//	@Summary		Register a cloud provider account with a billing account
//	@Description	Registers a new cloud provider account with the billing provider and returns a redirect URL for the customer onboarding flow.
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string						true	"Account ID"
//	@Param			request	body	RegisterLinkedCloudProviderRequest	true	"Request body"
//	@Success		201	{object}	RegisterAccountResponse
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id}/cloud-provider-accounts/register [post]
func (s *Server) RegisterCloudProviderAccount(c *gin.Context) {
	var req RegisterLinkedCloudProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	billingAccount, err := s.repos.BillingAccount.GetBillingAccountByID(c.Request.Context(), req.AccountID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	var csp repository.CloudServiceProvider
	for _, provider := range billingAccount.BillingProvider.SupportedCloudProviders {
		if provider.ID == req.CloudProviderID {
			csp = provider
			break
		}
	}
	if csp.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cloud provider not supported by this billing provider"})
		return
	}
	log.Printf("Found matching cloud provider in billing provider metadata: %+v", csp)

	if _, err := services.SyncCloudServiceProviderMetadata(c.Request.Context(), s.repos, csp.APIEndpointURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to sync cloud provider metadata"})
		return
	}

	// Embed the billing account and CSP provider IDs into the return URL so the frontend
	// can call the complete endpoint after the CSP redirect.
	u, _ := url.Parse(req.ReturnURL)
	q := u.Query()
	q.Set("billing_account_id", billingAccount.ID)
	q.Set("csp_provider_id", csp.ID)
	u.RawQuery = q.Encode()
	cspReturnURL := u.String()

	cspClient := clients.NewCloudServiceProviderClient(csp.APIEndpointURL)
	cspResponse, err := cspClient.RegisterCloudProviderAccount(billingAccount.BillingProvider.ID, billingAccount.ID, cspReturnURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initiate cloud provider onboarding"})
		return
	}

	c.JSON(http.StatusCreated, RegisterAccountResponse{
		AccountID:   cspResponse.ID,
		RedirectURL: cspResponse.RedirectURL,
	})
}

type completeCloudProviderOnboardingRequest struct {
	CspProviderID string `json:"csp_provider_id" binding:"required"`
}

// CompleteCloudProviderAccountOnboarding godoc
//
//	@Summary		Complete cloud provider account onboarding
//	@Description	Called by the frontend after the CSP redirect. Stores the CSP account link in the billing account.
//	@Tags			billing
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string									true	"Billing account ID"
//	@Param			request	body	completeCloudProviderOnboardingRequest	true	"Request body"
//	@Success		201		{object}	CloudProviderLink
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/billing/accounts/{id}/cloud-provider-accounts/complete [post]
func (s *Server) CompleteCloudProviderAccountOnboarding(c *gin.Context) {
	billingAccountID := c.Param("id")

	var req completeCloudProviderOnboardingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	link, err := s.repos.CloudServiceProviderAccount.Create(c.Request.Context(), billingAccountID, req.CspProviderID)
	if err != nil {
		log.Printf("Failed to store CSP account link for billing account %s: %v", billingAccountID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store cloud provider account link"})
		return
	}

	c.JSON(http.StatusCreated, CloudProviderLink{
		ID:              link.ID,
		CloudProviderID: link.CloudProviderID,
	})
}
