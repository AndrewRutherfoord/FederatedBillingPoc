package apicspadapter

import (
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/port"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ApiCustomerAdapter is the inbound HTTP server for customer-facing endpoints.
type ApiCustomerAdapter struct {
	handler *port.CustomerPortImpl
	config  *config.Config
	repos   *repository.Repos
}

func NewApiCustomerAdapter(handler *port.CustomerPortImpl, cfg *config.Config, repos *repository.Repos) *ApiCustomerAdapter {
	return &ApiCustomerAdapter{handler: handler, config: cfg, repos: repos}
}

func (a *ApiCustomerAdapter) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", a.Health)
	r.GET("/billing-accounts", a.ListBillingAccounts)
}

// Start creates the gin engine, registers routes, and blocks serving on addr.
func (a *ApiCustomerAdapter) Start(addr string) error {
	r := gin.Default()
	r.Use(cors.Default())
	a.RegisterRoutes(r)

	srv := &http.Server{Addr: addr, Handler: r}
	log.Printf("Starting %s (%s) customer API on %s", a.config.ProviderName, a.config.ProviderID, addr)
	return srv.ListenAndServe()
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (a *ApiCustomerAdapter) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "provider_id": a.repos.Provider.ID})
}

// ListBillingAccounts godoc
//
//	@Summary	List billing accounts for the authenticated customer
//	@Tags		billing-accounts
//	@Produce	json
//	@Success	200	{array}		map[string]interface{}
//	@Failure	500	{object}	map[string]string
//	@Router		/billing-accounts [get]
func (a *ApiCustomerAdapter) ListBillingAccounts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"billing_accounts": []any{}})
}
