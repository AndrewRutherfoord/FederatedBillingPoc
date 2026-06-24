package middleware

import (
	"net/http"
	"strings"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/gin-gonic/gin"
)

const contextKeyBillingAccount = "billing_account"

type AuthenticatedBillingAccount struct {
	ID string
}

// Placeholder for the bearer auth: token is the billing account ID itself, with no signature to verify currently.
// TODO: Replace with this with CSP issued JWT
func BillingAccountAuth(accounts repository.BillingAccountRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid bearer token"})
			return
		}

		account, err := accounts.GetByAccountID(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown billing account"})
			return
		}

		c.Set(contextKeyBillingAccount, &AuthenticatedBillingAccount{ID: account.AccountID})
		c.Next()
	}
}

func BillingAccountFromContext(c *gin.Context) *AuthenticatedBillingAccount {
	val, _ := c.Get(contextKeyBillingAccount)
	account, _ := val.(*AuthenticatedBillingAccount)
	return account
}
