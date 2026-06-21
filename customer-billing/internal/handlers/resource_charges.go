package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResourceChargeEntry struct {
	ResourceID             *string `json:"resource_id"`
	ResourceName           *string `json:"resource_name"`
	ResourceType           *string `json:"resource_type"`
	ServiceName            string  `json:"service_name"`
	ServiceCategory        string  `json:"service_category"`
	BillingCurrency        string  `json:"billing_currency"`
	CloudServiceProviderID string  `json:"cloud_service_provider_id"`
	TotalBilledCost        float64 `json:"total_billed_cost"`
	LineItemCount          int32   `json:"line_item_count"`
}

// ListResourceCharges godoc
//
//	@Summary		List aggregated billed cost per resource
//	@Description	Aggregates billed cost per resource entirely from the FOCUS line items fetched directly from the CSP (not from anything reported by the billing provider), so each row also shows which cloud service provider it came from.
//	@Tags			billing
//	@Produce		json
//	@Param			id	path		string	true	"Billing account ID"
//	@Success		200	{array}		ResourceChargeEntry
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id}/resource-charges [get]
func (s *Server) ListResourceCharges(c *gin.Context) {
	charges, err := s.repos.CloudServiceProviderFocusRecord.ListAggregatedResourceCharges(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list resource charges"})
		return
	}

	result := make([]ResourceChargeEntry, len(charges))
	for i, charge := range charges {
		result[i] = ResourceChargeEntry{
			ResourceID:             charge.ResourceID,
			ResourceName:           charge.ResourceName,
			ResourceType:           charge.ResourceType,
			ServiceName:            charge.ServiceName,
			ServiceCategory:        charge.ServiceCategory,
			BillingCurrency:        charge.BillingCurrency,
			CloudServiceProviderID: charge.CloudServiceProviderID,
			TotalBilledCost:        charge.TotalBilledCost,
			LineItemCount:          charge.LineItemCount,
		}
	}

	c.JSON(http.StatusOK, result)
}
