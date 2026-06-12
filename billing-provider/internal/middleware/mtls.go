package middleware

import (
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/gin-gonic/gin"
)

const contextKeyCSP = "csp"

type AuthenticatedCSP struct {
	ID   string
	Name string
}

// MTLSAuth validates the client certificate against configured CSPs and
// sets the authenticated CSP in the request context.
func MTLSAuth(csps []config.CloudServiceProvider) gin.HandlerFunc {
	cnToCSP := make(map[string]config.CloudServiceProvider)
	for _, csp := range csps {
		cnToCSP[csp.MTLS.CommonName] = csp
	}

	return func(c *gin.Context) {
		if c.Request.TLS == nil || len(c.Request.TLS.PeerCertificates) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "client certificate required"})
			return
		}

		cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName

		csp, ok := cnToCSP[cn]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unknown counterparty"})
			return
		}

		c.Set(contextKeyCSP, &AuthenticatedCSP{
			ID:   csp.ID,
			Name: csp.Name,
		})
		c.Next()
	}
}

func CSPFromContext(c *gin.Context) *AuthenticatedCSP {
	val, _ := c.Get(contextKeyCSP)
	csp, _ := val.(*AuthenticatedCSP)
	return csp
}
