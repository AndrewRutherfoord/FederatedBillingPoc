package apibpadapter

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/port"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/gin-gonic/gin"
)

// ApiBillingProviderAdapter is both the inbound mTLS HTTP server (receiving messages from BPs)
// and the outbound HTTP client (sending to BPs), implementing port.BPSender.
type ApiBillingProviderAdapter struct {
	handler        port.BPHandler
	config         *config.CspConfig
	clientRegistry *BPClientRegistry
}

func NewApiBillingProviderAdapter(handler port.BPHandler, config *config.CspConfig) *ApiBillingProviderAdapter {
	registry, err := NewBPClientRegistry(config)
	if err != nil {
		log.Fatalf("Failed to create BP client registry: %v", err)
	}
	return &ApiBillingProviderAdapter{
		handler:        handler,
		config:         config,
		clientRegistry: registry,
	}
}

func (t *ApiBillingProviderAdapter) RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.MTLSAuth(t.config.BillingProviders))
	{
		r.GET("/test", t.Test)
		// r.GET("/cost-batch-records", t.GetCostBatchRecords)
	}
}

// Start builds the mTLS http.Server and blocks serving requests on addr.
func (t *ApiBillingProviderAdapter) Start(addr string) error {
	r := gin.Default()
	t.RegisterRoutes(r)

	trustedPool := x509.NewCertPool()
	for _, bp := range t.config.BillingProviders {
		pem, err := os.ReadFile(bp.MTLS.CertPath)
		if err != nil {
			return fmt.Errorf("read cert for billing provider %s: %w", bp.Name, err)
		}
		if !trustedPool.AppendCertsFromPEM(pem) {
			return fmt.Errorf("append cert for billing provider %s failed", bp.Name)
		}
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequestClientCert,
			ClientCAs:  trustedPool,
		},
	}

	log.Printf("Starting %s (%s) adapter server on %s", t.config.ProviderName, t.config.ProviderID, addr)
	return srv.ListenAndServeTLS(t.config.MTLSCertPath, t.config.MTLSKeyPath)
}

func (t *ApiBillingProviderAdapter) Test(c *gin.Context) {
	bp := middleware.BPFromContext(c)
	log.Printf("Billing provider test endpoint called by %s", bp.Name)
	c.JSON(200, gin.H{"message": "hello from billing provider test endpoint"})
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
