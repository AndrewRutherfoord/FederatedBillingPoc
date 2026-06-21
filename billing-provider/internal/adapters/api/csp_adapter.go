package apicspadapter

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/port"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/gin-gonic/gin"
)

// ApiCSPAdapter is both the inbound mTLS HTTP server (receiving cost records pushed by CSPs)
// and the outbound HTTP client (querying CSPs), implementing port.CSPHandler via delegation.
type ApiCSPAdapter struct {
	handler        port.CSPHandler
	config         *config.Config
	repos          *repository.Repos
	clientRegistry *CSPClientRegistry
}

func NewApiCSPAdapter(handler port.CSPHandler, cfg *config.Config, repos *repository.Repos) *ApiCSPAdapter {
	registry, err := NewCSPClientRegistry(cfg)
	if err != nil {
		log.Fatalf("failed to create CSP client registry: %v", err)
	}
	return &ApiCSPAdapter{
		handler:        handler,
		config:         cfg,
		repos:          repos,
		clientRegistry: registry,
	}
}

func (a *ApiCSPAdapter) RegisterRoutes(r *gin.Engine) {
	r.GET("/health", a.Health)

	// All CSP-protocol routes are versioned and require a valid mTLS client cert.
	v1 := r.Group("/api/v1/billing", middleware.MTLSAuth(a.config.CloudServiceProviders))
	{
		v1.GET("/cloud-service-providers", a.ListCloudServiceProviders)
		v1.POST("/cost-records", a.ReceiveCostRecord)
	}
}

// Start builds the mTLS http.Server and blocks serving requests on addr.
func (a *ApiCSPAdapter) Start(addr string) error {
	r := gin.Default()
	a.RegisterRoutes(r)

	trustedPool := x509.NewCertPool()
	for _, csp := range a.config.CloudServiceProviders {
		pem, err := os.ReadFile(csp.MTLS.CertPath)
		if err != nil {
			return fmt.Errorf("read cert for CSP %s: %w", csp.Name, err)
		}
		if !trustedPool.AppendCertsFromPEM(pem) {
			return fmt.Errorf("append cert for CSP %s failed", csp.Name)
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

	log.Printf("Starting %s (%s) CSP API on %s", a.config.ProviderName, a.config.ProviderID, addr)
	return srv.ListenAndServeTLS(a.config.MTLSCertPath, a.config.MTLSKeyPath)
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (a *ApiCSPAdapter) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "provider_id": a.repos.Provider.ID})
}

// ListCloudServiceProviders godoc
//
//	@Summary	List all cloud service providers configured for this billing provider
//	@Tags		cloud-service-providers
//	@Produce	json
//	@Success	200	{array}		config.CloudServiceProvider
//	@Failure	500	{object}	map[string]string
//	@Router		/cloud-service-providers [get]
func (a *ApiCSPAdapter) ListCloudServiceProviders(c *gin.Context) {
	csps, err := a.repos.CloudServiceProviders.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cloud service providers"})
		return
	}
	c.JSON(http.StatusOK, csps)
}

// ReceiveCostRecord godoc
//
//	@Summary	Receive a charge batch from a CSP
//	@Tags		cost-records
//	@Accept		json
//	@Produce	json
//	@Param		batch	body		sharedmodels.ChargeBatch	true	"Charge batch"
//	@Success	202
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/cost-records [post]
func (a *ApiCSPAdapter) ReceiveCostRecord(c *gin.Context) {
	csp := middleware.CSPFromContext(c)
	log.Printf("Received charge batch from CSP %s", csp.Name)

	var batch sharedmodels.ChargeBatch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := a.handler.OnChargeBatch(c.Request.Context(), batch); err != nil {
		log.Printf("failed to process charge batch from CSP %s: %v", csp.Name, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process charge batch"})
		return
	}

	c.Status(http.StatusAccepted)
}

func (a *ApiCSPAdapter) Close() error {
	return nil
}
