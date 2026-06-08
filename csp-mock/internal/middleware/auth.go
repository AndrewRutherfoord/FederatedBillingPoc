package middleware

import (
	"net/http"
	"strings"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/gin-gonic/gin"
)

const contextKeyCustomer = "customer"

// Auth resolves the caller's customer from the Authorization header.
// Expected format: "Bearer <customer_id>"
func Auth(customers repository.CustomerRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}

		customerID := strings.TrimPrefix(header, "Bearer ")
		if customerID == header || customerID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be: Bearer <customer_id>"})
			return
		}

		customer, err := customers.GetByID(c.Request.Context(), customerID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown customer"})
			return
		}

		c.Set(contextKeyCustomer, customer)
		c.Next()
	}
}

// CustomerFromContext retrieves the authenticated customer set by Auth middleware.
func CustomerFromContext(c *gin.Context) *db.Customer {
	val, _ := c.Get(contextKeyCustomer)
	customer, _ := val.(*db.Customer)
	return customer
}
