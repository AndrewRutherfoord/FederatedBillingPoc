package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/middleware"
	bpmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetInvoices(c *gin.Context) {
	account := middleware.BillingAccountFromContext(c)

	invoices, err := s.repos.Invoice.ListByBillingAccount(c.Request.Context(), account.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list invoices"})
		return
	}

	response := bpmodels.GetInvoicesResponse{
		Invoices: make([]bpmodels.Invoice, len(invoices)),
		Count:    len(invoices),
	}
	for i, inv := range invoices {
		lineItems := make([]bpmodels.InvoiceProviderLineItem, len(inv.LineItems))
		for j, li := range inv.LineItems {
			lineItems[j] = bpmodels.InvoiceProviderLineItem{
				CloudServiceProviderID: li.CloudServiceProviderID,
				Amount:                 li.Amount,
				MerkleRoot:             li.MerkleRoot,
				BatchCount:             li.BatchCount,
			}
		}
		response.Invoices[i] = bpmodels.Invoice{
			ID:                inv.Invoice.ID,
			BillingAccountID:  inv.Invoice.BillingAccountID,
			BillingPeriodID:   inv.Invoice.BillingPeriodID,
			Amount:            inv.Invoice.Amount,
			Currency:          inv.Invoice.Currency,
			Status:            inv.Invoice.Status,
			IssuedAt:          inv.Invoice.IssuedAt,
			DueAt:             inv.Invoice.DueAt,
			ProviderLineItems: lineItems,
		}
	}

	c.JSON(http.StatusOK, response)
}
