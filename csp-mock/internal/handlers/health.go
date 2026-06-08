package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
		"provider_id": s.repos.Provider.ID,
	})
}
