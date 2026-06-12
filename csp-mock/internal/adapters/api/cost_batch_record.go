package apibpadapter

// import (
// 	"time"

// 	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
// 	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
// 	"github.com/gin-gonic/gin"
// )

// func parseTimeParam(c *gin.Context, key string) (time.Time, bool) {
// 	param := c.Query(key)

// 	if param == "" {
// 		c.JSON(400, gin.H{"error": "Missing '" + key + "' query parameter"})
// 		return time.Time{}, false
// 	}

// 	t, err := time.Parse(time.RFC3339, param)
// 	if err != nil {
// 		c.JSON(400, gin.H{"error": "Invalid '" + key + "' parameter"})
// 		return time.Time{}, false
// 	}

// 	return t, true
// }

// func (t *ApiBillingProviderAdapter) GetCostBatchRecords(c *gin.Context) {
// 	// curl --cert ./bp-1.crt --key bp-1.key -k https://localhost:8443/billing-provider/test
// 	bp := middleware.BPFromContext(c)

// 	// Get since and until query params
// 	since, ok := parseTimeParam(c, "since")
// 	if !ok {
// 		return
// 	}
// 	until, ok := parseTimeParam(c, "until")
// 	if !ok {
// 		return
// 	}

// 	batches, err := t.repos.CostBatch.ListByBillingProvider(c.Request.Context(), bp.ID, since, until)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Convert it to the message format defined in the shared package
// 	batchesMessages := make([]sharedmodels.AggregatedChargeRecord, len(batches))
// 	for _, batch := range batches {
// 		batchesMessages = append(batchesMessages, sharedmodels.AggregatedChargeRecord{
// 			BillingRecord: sharedmodels.BillingRecord{
// 				BillingProviderID:  batch.BillingProviderID,
// 				ResourceProviderID: batch.ResourceProviderID,
// 				BillingAccountID:   batch.BillingAccountID,
// 			},
// 			BatchID:         batch.ID,
// 			TotalBilledCost: batch.TotalCost,
// 			BilledCurrency:  batch.BilledCurrency,
// 			LineItemCount:   batch.TotalItems,
// 			BatchHash:       batch.MerkelRoot,
// 			BatchSignature:  batch.BatchSignature,
// 		})
// 	}

// 	c.JSON(200, gin.H{
// 		"cost_batches": batchesMessages,
// 	})
// }
