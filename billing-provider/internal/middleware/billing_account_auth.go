package middleware

import (
	"net/http"
	"strings"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	"github.com/gin-gonic/gin"
)

const contextKeyBillingAccount = "billing_account"

type AuthenticatedBillingAccount struct {
	ID string
}

// BillingAccountAuth is a placeholder for the customer API's bearer auth: the
// bearer token is the billing account ID itself, with no signature to verify.
// This will be replaced once POST /auth/token issues real signed JWTs.
func BillingAccountAuth(accounts repository.BillingAccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid bearer token"})
			return
		}

		account, err := accounts.Get(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown billing account"})
			return
		}

		c.Set(contextKeyBillingAccount, &AuthenticatedBillingAccount{ID: account.ID})
		c.Next()
	}
}

func BillingAccountFromContext(c *gin.Context) *AuthenticatedBillingAccount {
	val, _ := c.Get(contextKeyBillingAccount)
	account, _ := val.(*AuthenticatedBillingAccount)
	return account
}
