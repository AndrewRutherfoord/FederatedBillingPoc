package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server holds all application dependencies and owns route registration.
// Add new repositories to Repos as the API grows.
type Server struct {
	config *config.CspConfig
}

type SubServer interface {
	RegisterRoutes(r *gin.Engine)
}

func NewServer(config *config.CspConfig) *Server {
	return &Server{
		config: config,
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

func (s *Server) RegisterRoutes(r *gin.Engine, subservers []SubServer) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Health)

	for _, ss := range subservers {
		ss.RegisterRoutes(r)
	}
}
