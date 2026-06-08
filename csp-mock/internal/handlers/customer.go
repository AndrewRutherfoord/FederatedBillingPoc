package handlers

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/middleware"
	"github.com/gin-gonic/gin"
)

var _ *db.Customer // swag type anchor

type registerCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// RegisterCustomer godoc
//
//	@Summary	Register a new customer
//	@Tags		customer
//	@Accept		json
//	@Produce	json
//	@Param		body	body		registerCustomerRequest	true	"Customer details"
//	@Success	201		{object}	db.Customer
//	@Failure	400		{object}	map[string]string
//	@Failure	500		{object}	map[string]string
//	@Router		/customer/register [post]
func (s *Server) RegisterCustomer(c *gin.Context) {
	var req registerCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := s.repos.Customers.Create(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

// GetCustomer godoc
//
//	@Summary	Get the authenticated customer
//	@Tags		customer
//	@Produce	json
//	@Security	BearerAuth
//	@Success	200		{object}	db.Customer
//	@Failure	401		{object}	map[string]string
//	@Router		/customer [get]
func (s *Server) GetCustomer(c *gin.Context) {
	customer := middleware.CustomerFromContext(c)
	c.JSON(http.StatusOK, customer)
}
