package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/gin-gonic/gin"
)

// GetChargeBatch godoc
//
//	@Summary	Get the details of a specific charge batch
//	@Description	Called directly by the customer-billing service to independently verify a charge batch it already knows about from the billing provider.
//	@Tags		billing
//	@Produce	json
//	@Param		batch_id	path		string	true	"Charge batch ID"
//	@Success	200			{object}	sharedmodels.ChargeBatchDetail
//	@Failure	404			{object}	map[string]string
//	@Router		/billing/charge-batches/{batch_id} [get]
func (s *Server) GetChargeBatch(c *gin.Context) {
	batchID := c.Param("batch_id")

	batch, err := s.repos.ChargeBatch.GetByID(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "charge batch not found"})
		return
	}

	account := middleware.BillingAccountFromContext(c)
	if batch.BillingAccountID != account.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "charge batch not found"})
		return
	}

	lineItems, err := s.repos.Focus.List(c.Request.Context(), repository.FocusFilter{
		BillingAccountID: batch.BillingAccountID,
		BatchID:          batch.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch line items"})
		return
	}

	c.JSON(http.StatusOK, sharedmodels.ChargeBatchDetail{
		ChargeBatch: sharedmodels.ChargeBatch{
			BillingContext: sharedmodels.BillingContext{
				BillingProviderID:      batch.BillingProviderID,
				CloudServiceProviderID: batch.CloudServiceProviderID,
				BillingAccountID:       batch.BillingAccountID,
			},
			BatchID:         batch.ID,
			TotalBilledCost: batch.TotalCost,
			BilledCurrency:  batch.BilledCurrency,
			LineItemCount:   batch.TotalItems,
			MerkleRoot:      batch.MerkleRoot,
			BatchSignature:  batch.BatchSignature,
			CreatedAt:       batch.CreatedAt,
		},
		LineItems: lineItems,
	})
}