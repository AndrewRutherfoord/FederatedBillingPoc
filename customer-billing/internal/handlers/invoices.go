package handlers

import (
	"net/http"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/clients"
	"github.com/gin-gonic/gin"
)

type InvoiceProviderLineItemEntry struct {
	CloudServiceProviderID string  `json:"cloud_service_provider_id"`
	Amount                 float64 `json:"amount"`
	MerkleRoot             string  `json:"merkle_root"`
	BatchCount             int     `json:"batch_count"`
}

type InvoiceEntry struct {
	ID                string                         `json:"id"`
	BillingPeriodID   string                         `json:"billing_period_id"`
	Amount            float64                        `json:"amount"`
	Currency          string                         `json:"currency"`
	Status            string                         `json:"status"`
	IssuedAt          time.Time                      `json:"issued_at"`
	DueAt             time.Time                      `json:"due_at"`
	ProviderLineItems []InvoiceProviderLineItemEntry `json:"provider_line_items"`
}

// ListInvoices godoc
//
//	@Summary		List invoices for a billing account
//	@Description	Fetches the billing account's invoices live from its billing provider.
//	@Tags			billing
//	@Produce		json
//	@Param			id	path		string	true	"Billing account ID"
//	@Success		200	{array}		InvoiceEntry
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id}/invoices [get]
func (s *Server) ListInvoices(c *gin.Context) {
	id := c.Param("id")

	account, err := s.repos.BillingAccount.GetBillingAccountByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing account not found"})
		return
	}

	client := clients.NewBillingProviderClient(account.BillingProvider.BaseURL)
	invoices, err := client.GetInvoices(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch invoices"})
		return
	}

	result := make([]InvoiceEntry, len(invoices))
	for i, inv := range invoices {
		lineItems := make([]InvoiceProviderLineItemEntry, len(inv.ProviderLineItems))
		for j, li := range inv.ProviderLineItems {
			lineItems[j] = InvoiceProviderLineItemEntry{
				CloudServiceProviderID: li.CloudServiceProviderID,
				Amount:                 li.Amount,
				MerkleRoot:             li.MerkleRoot,
				BatchCount:             li.BatchCount,
			}
		}
		result[i] = InvoiceEntry{
			ID:                inv.ID,
			BillingPeriodID:   inv.BillingPeriodID,
			Amount:            inv.Amount,
			Currency:          inv.Currency,
			Status:            inv.Status,
			IssuedAt:          inv.IssuedAt,
			DueAt:             inv.DueAt,
			ProviderLineItems: lineItems,
		}
	}

	c.JSON(http.StatusOK, result)
}
