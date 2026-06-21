package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
	"github.com/gin-gonic/gin"
)

// GetBillingAccountRecords godoc
//
//	@Summary	Get charge batches (including raw FOCUS line items) for a billing account
//	@Description	Called directly by the customer-billing service to independently verify what this CSP reported, without going through the billing provider.
//	@Tags		billing
//	@Accept		json
//	@Produce	json
//	@Param		request	body		cspsharedmodels.GetBillingAccountRecordsRequest	true	"Request body"
//	@Success	200		{object}	cspsharedmodels.GetBillingAccountRecordsResponse
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/billing/accounts/records [post]
func (s *Server) GetBillingAccountRecords(c *gin.Context) {
	var req cspsharedmodels.GetBillingAccountRecordsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batches, err := s.repos.ChargeBatch.ListByBillingAccount(c.Request.Context(), req.BillingAccountID, req.From, req.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch charge batches"})
		return
	}

	details := make([]sharedmodels.ChargeBatchDetail, 0, len(batches))
	for _, batch := range batches {
		lineItems, err := s.repos.Focus.List(c.Request.Context(), repository.FocusFilter{
			BillingAccountID: req.BillingAccountID,
			BatchID:          batch.ID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch line items"})
			return
		}

		details = append(details, sharedmodels.ChargeBatchDetail{
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

	c.JSON(http.StatusOK, cspsharedmodels.GetBillingAccountRecordsResponse{
		Batches: details,
		Count:   len(details),
		From:    req.From,
		To:      req.To,
	})
}
