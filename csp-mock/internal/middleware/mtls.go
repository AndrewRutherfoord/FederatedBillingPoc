package middleware

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/gin-gonic/gin"
)

const contextKeyBP = "csp"

type AuthenticatedBP struct {
	ID   string
	Name string
}

func MTLSAuth(bps []config.BillingProvider) gin.HandlerFunc {
	// Build a lookup map from CN -> CSP config at startup, not per request
	cnToBp := make(map[string]config.BillingProvider)
	for _, bp := range bps {
		cnToBp[bp.MTLS.CommonName] = bp
	}

	return func(c *gin.Context) {
		if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client certificate required"})
			return
		}

		cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName

		bp, ok := cnToBp[cn]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown counterparty"})
			return
		}

		c.Set(contextKeyBP, &AuthenticatedBP{
			ID:   bp.ID,
			Name: bp.Name,
		})
		c.Next()
	}
}

func BPFromContext(c *gin.Context) *AuthenticatedBP {
	val, _ := c.Get(contextKeyBP)
	csp, _ := val.(*AuthenticatedBP)
	return csp
}
