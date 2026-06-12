package apibpadapter

import (
	"context"
	"fmt"
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/port"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/gin-gonic/gin"
)

// ApiBillingProviderAdapter is both the inbound HTTP server (receiving messages from BPs)
// and the outbound HTTP client (sending to BPs), implementing port.BPSender.
type ApiBillingProviderAdapter struct {
	handler        port.BPHandler
	config         *config.CspConfig
	clientRegistry *BPClientRegistry
}

func NewApiBillingProviderAdapter(r *gin.Engine, handler port.BPHandler, config *config.CspConfig) *ApiBillingProviderAdapter {
	registry, err := NewBPClientRegistry(config)
	if err != nil {
		log.Fatalf("Failed to create BP client registry: %v", err)
	}

	t := ApiBillingProviderAdapter{
		handler:        handler,
		config:         config,
		clientRegistry: registry,
	}
	t.RegisterRoutes(r)

	return &t
}

func (t *ApiBillingProviderAdapter) Test(c *gin.Context) {
	bp := middleware.BPFromContext(c)
	log.Printf("Billing provider test endpoint called by %s", bp.Name)
	c.JSON(200, gin.H{"message": "hello from billing provider test endpoint"})
}

func (t *ApiBillingProviderAdapter) RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.MTLSAuth(t.config.BillingProviders))
	{
		r.GET("/test", t.Test)
		// r.GET("/cost-batch-records", t.GetCostBatchRecords)
	}
}

func (t *ApiBillingProviderAdapter) Close() error {
	return nil
}

// SendAggregatedChargeRecord implements port.BPSender. It POSTs the record to the
// appropriate billing provider's API endpoint.
func (t *ApiBillingProviderAdapter) SendAggregatedChargeRecord(ctx context.Context, record sharedmodels.AggregatedChargeRecord) error {
	log.Printf("Sending aggregated charge record for batch %s to billing provider %s", record.BatchID, record.BillingRecord.BillingProviderID)
	client, ok := t.clientRegistry.GetClient(record.BillingRecord.BillingProviderID)
	if !ok {
		log.Printf("No client found for billing provider %s", record.BillingRecord.BillingProviderID)
		return fmt.Errorf("no client registered for billing provider %s", record.BillingRecord.BillingProviderID)
	}
	log.Printf("Found client for billing provider %s, sending record to %s", record.BillingRecord.BillingProviderID, client.BaseURL)
	return client.SendJSON(client.BaseURL+"/cost-records", record)
}
