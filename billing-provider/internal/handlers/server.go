package handlers

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	bpsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed templates/*
var templatesFS embed.FS

type Server struct {
	config *config.Config
	repos  *repository.Repos
}

func NewServer(config *config.Config, repos *repository.Repos) *Server {
	return &Server{config: config, repos: repos}
}

// Health godoc
//
//	@Summary	Health check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "provider_id": s.repos.Provider.ID})
}

func (s *Server) WellKnown(c *gin.Context) {
	csps, err := s.repos.CloudServiceProviders.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list cloud service providers"})
		return
	}
	var cspMetas []bpsharedmodels.SupportedCloudProviderMetadata
	for _, csp := range csps {
		cspMetas = append(cspMetas, bpsharedmodels.SupportedCloudProviderMetadata{
			ID:          csp.ID,
			Name:        csp.Name,
			APIEndpoint: csp.APIEndpointURL,
		})
	}
	metadata := bpsharedmodels.Metadata{
		ID:                      s.repos.Provider.ID,
		Name:                    s.repos.Provider.Name,
		SupportedCloudProviders: cspMetas,
	}
	c.JSON(http.StatusOK, metadata)
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Health)
	r.GET("/.well-known/billing-provider", s.WellKnown)
	r.POST("/billing/accounts", s.RegisterBillingAccount)
	r.GET("/billing/accounts/:id/onboard", s.OnboardForm)
	r.POST("/billing/accounts/:id/onboard", s.OnboardSubmit)
	r.POST("/billing/accounts/records", s.GetBillingAccountRecords)
}

// Start creates the gin engine, registers routes, and blocks serving on addr.
func (s *Server) Start(addr string) error {
	r := gin.Default()
	r.Use(cors.Default())

	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	s.RegisterRoutes(r)

	srv := &http.Server{Addr: addr, Handler: r}
	log.Printf("Starting %s (%s) customer API on %s", s.config.ProviderName, s.config.ProviderID, addr)
	return srv.ListenAndServe()
}
