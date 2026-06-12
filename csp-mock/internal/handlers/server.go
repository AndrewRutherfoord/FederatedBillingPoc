package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server holds all application dependencies and owns route registration.
// Add new repositories to Repos as the API grows.
type Server struct {
	config *config.CspConfig
	repos  *repository.Repos
	clock  util.Clock
}

func NewServer(config *config.CspConfig, repos *repository.Repos, clock util.Clock) *Server {
	return &Server{
		config: config,
		repos:  repos,
		clock:  clock,
	}
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"provider_id": s.config.ProviderID,
	})
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Health)

	group := r.Group("/customer")

	group.GET("/resource-types", s.ListResourceTypes)
	group.GET("/resource-types/:id", s.GetResourceType)

	group.POST("/customer/register", s.RegisterCustomer)

	// Routes below require a valid customer in the Authorization header.
	authed := group.Group("/", middleware.Auth(s.repos.Customers))
	{
		authed.GET("/customer", s.GetCustomer)

		authed.GET("/resources", s.ListResources)
		authed.POST("/resources", s.CreateResource)
		authed.DELETE("/resources/:id", s.DeleteResource)

		authed.GET("/billing-accounts", s.ListBillingAccounts)
		authed.GET("/billing-accounts/:account_id", s.GetBillingAccount)
		authed.POST("/billing-accounts", s.CreateBillingAccount)
	}
}
