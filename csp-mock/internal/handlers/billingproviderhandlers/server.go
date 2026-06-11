package billingproviderhandlers

import (
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
	"github.com/gin-gonic/gin"
)

type BillingProviderServer struct {
	billingProviders []config.BillingProvider // Needed for MTLSAuth middleware to look up the billing provider based on the client cert's CN
	repos            *repository.Repos
	clock            util.Clock
}

func NewBillingProviderServer(billingProviders []config.BillingProvider, repos *repository.Repos, clock util.Clock) *BillingProviderServer {
	return &BillingProviderServer{billingProviders: billingProviders, repos: repos, clock: clock}
}

func (bs *BillingProviderServer) Test(c *gin.Context) {
	// curl --cert ./bp-1.crt --key bp-1.key -k https://localhost:8443/billing-provider/test
	bp := middleware.BPFromContext(c)

	log.Printf("Billing provider test endpoint called by %s", bp.Name)

	c.JSON(200, gin.H{
		"message": "hello from billing provider test endpoint",
	})
}

func (bs *BillingProviderServer) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/billing-provider")

	group.Use(middleware.MTLSAuth(bs.billingProviders))
	{
		group.GET("/test", bs.Test)
		group.GET("/cost-batch-records", bs.GetCostBatchRecords)
	}
}
