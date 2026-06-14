package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	repos *repository.Repos
}

func NewServer(repos *repository.Repos) *Server {
	return &Server{
		repos: repos,
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
		"status": "ok",
	})
}

func (s *Server) RegisterRoutes(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", s.Health)

	r.POST("/billing/accounts/register", s.RegisterBillingAccount)
	r.GET("/billing/accounts", s.ListBillingAccounts)
	r.GET("/billing/accounts/:id", s.GetBillingAccount)
}
