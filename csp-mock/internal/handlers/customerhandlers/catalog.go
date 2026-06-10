package customerhandlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/gin-gonic/gin"
)

var _ config.ResourceType // swag type anchor

// ListResourceTypes godoc
//
//	@Summary	List all resource types offered by this CSP
//	@Tags		catalog
//	@Produce	json
//	@Success	200	{array}		config.ResourceType
//	@Router		/resource-types [get]
func (cs *CustomerServer) ListResourceTypes(c *gin.Context) {
	c.JSON(http.StatusOK, cs.repos.ResourceTypes.List(c.Request.Context()))
}

// GetResourceType godoc
//
//	@Summary	Get a resource type by ID
//	@Tags		catalog
//	@Produce	json
//	@Param		id	path		string	true	"Resource type ID"
//	@Success	200	{object}	config.ResourceType
//	@Failure	404	{object}	map[string]string
//	@Router		/resource-types/{id} [get]
func (cs *CustomerServer) GetResourceType(c *gin.Context) {
	id := c.Param("id")
	rt, err := cs.repos.ResourceTypes.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource type not found", "id": id})
		return
	}
	c.JSON(http.StatusOK, rt)
}
