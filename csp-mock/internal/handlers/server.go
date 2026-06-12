package handlers

import (
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
	"github.com/gin-contrib/cors"
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

	r.GET("/resource-types", s.ListResourceTypes)
	r.GET("/resource-types/:id", s.GetResourceType)

	r.POST("/customer/register", s.RegisterCustomer)

	// Routes below require a valid customer in the Authorization header.
	authed := r.Group("/", middleware.Auth(s.repos.Customers))
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

// Start creates the gin engine, registers routes, and blocks serving on addr.
func (s *Server) Start(addr string) error {
	r := gin.Default()
	r.Use(cors.Default())
	s.RegisterRoutes(r)

	srv := &http.Server{Addr: addr, Handler: r}
	log.Printf("Starting %s (%s) customer API on %s", s.config.ProviderName, s.config.ProviderID, addr)
	return srv.ListenAndServe()
}
