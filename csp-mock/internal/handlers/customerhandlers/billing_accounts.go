package customerhandlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/gin-gonic/gin"
)

var _ db.Resource // swag type anchor

type createBillingAccountRequest struct {
	AccountID       string `json:"account_id" binding:"required"`
	BillingProvider string `json:"billing_provider" binding:"required"`
}

// ListBillingAccounts godoc
//
//	@Summary	List all billing accounts for the authenticated customer
//	@Tags		billing-accounts
//	@Produce	json
//	@Security	BearerAuth
//	@Success	200	{array}		db.BillingAccount
//	@Failure	401	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/billing-accounts [get]
func (cs *CustomerServer) ListBillingAccounts(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)

	accounts, err := cs.repos.BillingAccounts.ListByCustomer(c.Request.Context(), customer.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, accounts)
}

// GetBillingAccount godoc
//
//	@Summary	Get a billing account by account ID
//	@Tags		billing-accounts
//	@Produce	json
//	@Security	BearerAuth
//	@Param		account_id	path		string	true	"Billing account ID"
//	@Success	200			{object}	db.BillingAccount
//	@Failure	401			{object}	map[string]string
//	@Failure	404			{object}	map[string]string
//	@Failure	500			{object}	map[string]string
//	@Router		/billing-accounts/{account_id} [get]
func (cs *CustomerServer) GetBillingAccount(c *gin.Context) {
	accountID := c.Param("account_id")

	account, err := cs.repos.BillingAccounts.GetByAccountID(c.Request.Context(), accountID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Billing account not found"})
		return
	}

	c.JSON(200, account)
}

// CreateBillingAccount godoc
//
//	@Summary	Create a new billing account for the authenticated customer
//	@Tags		billing-accounts
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		body    body		createBillingAccountRequest	true	"Billing account details"
//	@Success	201		{object}	db.BillingAccount
//	@Failure	400		{object}	map[string]string
//	@Failure	401		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/billing-accounts [post]
func (cs *CustomerServer) CreateBillingAccount(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)

	var req createBillingAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	billingProvider, err := cs.repos.BillingProviders.Get(c.Request.Context(), req.BillingProvider)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid billing provider"})
		return
	}

	// TODO: Validate account ID with Billing Provider

	account := db.BillingAccount{
		AccountID:         req.AccountID,
		BillingProviderID: billingProvider.ID,
		CustomerID:        customer.ID,
		CreatedAt:         cs.clock.Now(),
		UpdatedAt:         cs.clock.Now(),
	}

	newAccount, err := cs.repos.BillingAccounts.Create(c.Request.Context(), account)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, newAccount)
}
