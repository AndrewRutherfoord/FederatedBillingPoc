package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed templates/*
var templatesFS embed.FS

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

// WellKnown godoc
//
//	@Summary	Cloud Service Provider metadata
//	@Description	Returns metadata about the cloud service provider, such as provider ID and name. This is used by the billing provider or the customer metering service to identify and display information about the CSP.
//	@Tags		well-known
//	@Produce	json
//	@Success	200	{object}	cspsharedmodels.Metadata
//	@Router		/.well-known/cloud-service-provider [get]
func (s *Server) WellKnown(c *gin.Context) {
	bps, err := s.repos.BillingProviders.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list billing providers"})
		return
	}
	var bpMetas []cspsharedmodels.SupportedCloudProviderMetadata
	for _, bp := range bps {
		bpMetas = append(bpMetas, cspsharedmodels.SupportedCloudProviderMetadata{
			ID:   bp.ID,
			Name: bp.Name,
		})
	}
	metadata := cspsharedmodels.Metadata{
		ID:                        s.config.ProviderID,
		Name:                      s.config.ProviderName,
		APIEndpointURL:            s.config.CustomerAPIEndpointURL,
		SupportedBillingProviders: bpMetas,
	}
	c.JSON(http.StatusOK, metadata)
}

func (s *Server) Jwks(c *gin.Context) {
	// jwks := s.config.Jwks
	jwks := map[string]interface{}{}
	c.JSON(http.StatusOK, jwks)
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	billingGroup := api.Group("/billing")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", s.Health)
	r.GET("/.well-known/cloud-service-provider", s.WellKnown)
	r.GET("/.well-known/jwks.json", s.Jwks)

	// These routes all act as the actual CSP and are not for access by the CMS.
	// TODO: Move this to a seperate service eventually since this is the billing adapter and not the customer API.

	api.GET("/resource-types", s.ListResourceTypes)
	api.GET("/resource-types/:id", s.GetResourceType)

	api.POST("/customer/register", s.RegisterCustomer)

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

	// Called by the customer-billing service to initiate onboarding (no customer auth).
	billingGroup.POST("/onboarding", s.InitiateBillingAccountOnboarding)

	billingAuth := billingGroup.Group("/", middleware.BillingAccountAuth(s.repos.BillingAccounts))
	{
		// Called by the customer-billing service to independently fetch this CSP's view of
		// a charge batch it already knows about from the billing provider, bypassing the
		// billing provider. Includes the line items for the batch.
		billingAuth.GET("/charge-batches/:batch_id", s.GetChargeBatch)
	}

	// SSR routes for the onboarding forms
	r.GET("/onboarding/:session_id", s.OnboardingForm)
	r.POST("/onboarding/:session_id", s.OnboardingSubmit)
	r.GET("/onboarding/:session_id/complete", s.OnboardingComplete)
	r.GET("/onboarding/:session_id/complete/download", s.OnboardingCompleteDownload)
}

// Start creates the gin engine, registers routes, and blocks serving on addr.
func (s *Server) Start(addr string) error {
	r := gin.Default()
	r.Use(cors.Default())
	r.SetHTMLTemplate(mustParseTemplates())

	s.RegisterRoutes(r)

	srv := &http.Server{Addr: addr, Handler: r}
	log.Printf("Starting %s (%s) customer API on %s", s.config.ProviderName, s.config.ProviderID, addr)
	return srv.ListenAndServe()
}

func mustParseTemplates() *template.Template {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		panic("failed to parse templates: " + err.Error())
	}
	return tmpl
}
