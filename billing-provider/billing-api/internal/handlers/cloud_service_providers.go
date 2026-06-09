package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListCloudServiceProviders godoc
//
//	@Summary	List all cloud service providers supported by this billing provider
//	@Tags		cloud-service-providers
//	@Produce	json
//	@Success	200	{array}	map[string]interface{}
//	@Failure	500	{object}	map[string]string
//	@Router		/cloud-service-providers [get]
func (s *Server) ListCloudServiceProviders(c *gin.Context) {
	csps, err := s.repos.CloudServiceProviders.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list cloud service providers"})
		return
	}
	c.JSON(http.StatusOK, csps)
}
