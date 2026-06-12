package apibpadapter

import (
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/gin-gonic/gin"
)

type ApiBillingProviderAdapter struct {
	r              *gin.Engine
	config         *config.CspConfig
	repos          *repository.Repos
	clientRegistry *BPClientRegistry
}

func NewApiBillingProviderAdapter(r *gin.Engine, repos *repository.Repos, config *config.CspConfig) *ApiBillingProviderAdapter {

	registry, error := NewBPClientRegistry(config)
	if error != nil {
		log.Fatalf("Failed to create BP client registry: %v", error)
	}

	t := ApiBillingProviderAdapter{
		repos:          repos,
		config:         config,
		clientRegistry: registry,
	}
	t.RegisterRoutes(r)

	return &t
}

func (t *ApiBillingProviderAdapter) Test(c *gin.Context) {
	// curl --cert ./bp-1.crt --key bp-1.key -k https://localhost:8443/billing-provider/test
	bp := middleware.BPFromContext(c)

	log.Printf("Billing provider test endpoint called by %s", bp.Name)

	c.JSON(200, gin.H{
		"message": "hello from billing provider test endpoint",
	})
}

func (t *ApiBillingProviderAdapter) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/billing-provider")

	group.Use(middleware.MTLSAuth(t.config.BillingProviders))
	{
		group.GET("/test", t.Test)
		group.GET("/cost-batch-records", t.GetCostBatchRecords)
	}
}

func (t *ApiBillingProviderAdapter) Close() error {
	return nil
}

func (t *ApiBillingProviderAdapter) SendCostBatchRecord(record sharedmodels.AggregatedChargeRecord) error {
	return nil
}
