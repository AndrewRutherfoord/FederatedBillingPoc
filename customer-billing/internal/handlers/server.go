package handlers

import (
	"errors"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type Server struct {
	repos *repository.Repos
}

func NewServer(repos *repository.Repos) *Server {
	return &Server{
		repos: repos,
	}
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

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

	client := clients.NewBillingProviderClient(req.BillingProviderBaseURL)

	// Get the metadata from the bp's .well-known endpoint
	metadata, err := client.GetMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch billing provider metadata"})
		return
	}

	// Create the BP in the db if it doesn't exist
	bp, err := s.repos.BillingProvider.GetBillingProviderByID(metadata.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to look up billing provider"})
			return
		}
		bp, err = s.repos.BillingProvider.CreateBillingProvider(metadata.ID, metadata.Name, req.BillingProviderBaseURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save billing provider"})
			return
		}
	}

	// Register the billing account with the billing provider.
	// Name and email are collected later via the BP's onboarding form.
	bpResponse, err := client.RegisterBillingAccount(req.ReturnURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Store the pending account locally
	// TODO: Handle token
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

// GetBillingAccounts godoc
//
//	@Summary		Get all billing accounts
//	@Tags			billing
//	@Produce		json
//	@Success		200	{array}		[]billingAccountResponse
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts [get]
func (s *Server) GetBillingAccounts(c *gin.Context) {
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

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Health)

	r.POST("/billing/accounts/register", s.RegisterBillingAccount)
	r.GET("/billing/accounts", s.GetBillingAccounts)
}
