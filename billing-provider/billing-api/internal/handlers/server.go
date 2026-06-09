package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	repos *repository.Repos
	clock shared.Clock
}

func NewServer(repos *repository.Repos, clock shared.Clock) *Server {
	return &Server{repos: repos, clock: clock}
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

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", s.Health)

	r.GET("/cloud-service-providers", s.ListCloudServiceProviders)
}
