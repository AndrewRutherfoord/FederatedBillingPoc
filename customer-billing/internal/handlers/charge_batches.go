package handlers

import (
	"net/http"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/gin-gonic/gin"
)

type chargeBatchReportEntry struct {
	TotalItems     int32     `json:"total_items"`
	TotalCost      float64   `json:"total_cost"`
	BilledCurrency string    `json:"billed_currency"`
	MerkleRoot     string    `json:"merkle_root"`
	BatchSignature string    `json:"batch_signature"`
	CreatedAt      time.Time `json:"created_at"`
	ReceivedAt     time.Time `json:"received_at"`
}

type ChargeBatchEntry struct {
	BatchID                string                  `json:"batch_id"`
	BillingAccountID       string                  `json:"billing_account_id"`
	CloudServiceProviderID string                  `json:"cloud_service_provider_id"`
	Status                 string                  `json:"status"`
	BillingProviderReport  *chargeBatchReportEntry `json:"billing_provider_report"`
	CloudProviderReport    *chargeBatchReportEntry `json:"cloud_provider_report"`
}

// ListChargeBatches godoc
//
//	@Summary		List charge batches for a billing account, merging the billing provider's and cloud provider's reports
//	@Description	Each charge batch is returned once, with the billing provider's report and the cloud service provider's own report of the same batch ID side by side (either may be missing), plus a status indicating whether they agree.
//	@Tags			billing
//	@Produce		json
//	@Param			id	path		string	true	"Billing account ID"
//	@Success		200	{array}		ChargeBatchEntry
//	@Failure		500	{object}	map[string]string
//	@Router			/billing/accounts/{id}/charge-batches [get]
func (s *Server) ListChargeBatches(c *gin.Context) {
	batches, err := s.repos.ChargeBatchReconciliation.ListMergedBatches(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list charge batches"})
		return
	}

	result := make([]ChargeBatchEntry, len(batches))
	for i, b := range batches {
		result[i] = ChargeBatchEntry{
			BatchID:                b.BatchID,
			BillingAccountID:       b.BillingAccountID,
			CloudServiceProviderID: b.CloudServiceProviderID,
			Status:                 string(b.Status),
			BillingProviderReport:  toChargeBatchReportEntry(b.BillingProviderReport),
			CloudProviderReport:    toChargeBatchReportEntry(b.CloudProviderReport),
		}
	}

	c.JSON(http.StatusOK, result)
}

func toChargeBatchReportEntry(report *repository.ChargeBatchReport) *chargeBatchReportEntry {
	if report == nil {
		return nil
	}
	return &chargeBatchReportEntry{
		TotalItems:     report.TotalItems,
		TotalCost:      report.TotalCost,
		BilledCurrency: report.BilledCurrency,
		MerkleRoot:     report.MerkleRoot,
		BatchSignature: report.BatchSignature,
		CreatedAt:      report.CreatedAt,
		ReceivedAt:     report.ReceivedAt,
	}
}
