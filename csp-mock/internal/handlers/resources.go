package handlers

import (
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var _ db.Resource // swag type anchor

type createResourceRequest struct {
	BillingAccountID string `json:"billing_account_id" binding:"required"`
	ResourceType     string `json:"resource_type" binding:"required"`
}

type setStorageRequest struct {
	StorageGB float64 `json:"storage_gb" binding:"required"`
}

// ListResources godoc
//
//	@Summary	List all resources for the authenticated customer
//	@Tags		resources
//	@Produce	json
//	@Security	BearerAuth
//	@Success	200	{array}		db.Resource
//	@Failure	401	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/resources [get]
func (s *Server) ListResources(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)

	resources, err := s.repos.Resources.ListByCustomerID(c.Request.Context(), customer.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resources)
}

// CreateResource godoc
//
// @Summary	Provision a new resource for the authenticated customer
// @Tags		resources
// @Accept		json
// @Produce	json
// @Security	BearerAuth
// @Param		body	body		createResourceRequest	true	"Resource details"
// @Success	201		{object}	db.Resource
// @Failure	400		{object}	map[string]string
// @Failure	401		{object}	map[string]string
// @Failure	500		{object}	map[string]string
// @Router		/resources [post]
func (s *Server) CreateResource(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)

	var req createResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	billingAccountID := req.BillingAccountID
	billingAccount, err := s.repos.BillingAccounts.GetByAccountID(c.Request.Context(), billingAccountID)
	if err != nil || billingAccount.CustomerID != customer.ID {
		c.JSON(404, gin.H{"error": "Billing account not found"})
		return
	}

	resourceType, err := s.repos.ResourceTypes.Get(req.ResourceType)
	if err != nil {
		c.JSON(404, gin.H{"error": "Resource type not found"})
		return
	}

	var storageGB *decimal.Decimal
	if resourceType.Pricing.Model == "per_gb_hour" {
		val := decimal.NewFromInt(1)
		storageGB = &val // Default to 1 GB for storage-based pricing
	}

	resource, err := s.repos.Resources.Create(c.Request.Context(), customer.ID, billingAccountID, req.ResourceType, storageGB)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, resource)
}

// SetStorageGB godoc
//
//	@Summary	Update the storage GB for a resource
//	@Tags		resources
//	@Accept		json
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	string				true	"Resource ID"
//	@Param		body	body	setStorageRequest	true	"Storage details"
//	@Success	204
//	@Failure	400	{object}	map[string]string
//	@Failure	401	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/resources/{id}/storage [put]
func (s *Server) SetStorageGB(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)
	resourceID := c.Param("id")

	var req setStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resource, err := s.repos.Resources.GetByID(c.Request.Context(), resourceID)
	if err != nil || resource.CustomerID != customer.ID {
		c.JSON(404, gin.H{"error": "Resource not found"})
		return
	}

	if err := s.repos.Resources.SetStorageGB(c.Request.Context(), resourceID, decimal.NewFromFloat(req.StorageGB)); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

// DeleteResource godoc
//
//	@Summary	Delete a resource by ID
//	@Tags		resources
//	@Produce	json
//	@Security	BearerAuth
//	@Param		id	path	string	true	"Resource ID"
//	@Success	204
//	@Failure	401	{object}	map[string]string
//	@Failure	404	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/resources/{id} [delete]
func (s *Server) DeleteResource(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)
	resourceID := c.Param("id")

	resource, err := s.repos.Resources.GetByID(c.Request.Context(), resourceID)
	if err != nil || resource.CustomerID != customer.ID {
		c.JSON(404, gin.H{"error": "Resource not found"})
		return
	}

	if err := s.repos.Resources.Delete(c.Request.Context(), resourceID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}
