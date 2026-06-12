// Package main is the entry point for the CSP Mock Adapter server.
// The adapter serves as the API endpoint for billing providers to interact with the CSP.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"

	apibpadapter "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/adapters/api"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/app"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	port := flag.String("port", ":8443", "port to listen on")
	flag.Parse()

	// Initialize app with shared setup
	appInstance := app.NewApp(app.Config{
		ConfigPath: *configPath,
	})

	r := gin.Default()

	// Build the trusted cert pool for CSPs (billing provider certs)
	trustedPool := x509.NewCertPool()
	for _, bp := range appInstance.Config.BillingProviders {
		trustedCertPEM, err := os.ReadFile(bp.MTLS.CertPath)
		if err != nil {
			log.Fatalf("failed to read cert for billing provider %s: %v", bp.Name, err)
		}
		if !trustedPool.AppendCertsFromPEM(trustedCertPEM) {
			log.Fatalf("failed to append cert for billing provider %s", bp.Name)
		}
	}

	// Adapter routes with mTLS authentication
	adapter := apibpadapter.NewApiBillingProviderAdapter(r, appInstance.Repos, appInstance.Config)
	defer adapter.Close()

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  trustedPool,
	}

	httpServer := &http.Server{
		Addr:      *port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting %s (%s) adapter server on %s", appInstance.Repos.Provider.Name, appInstance.Repos.Provider.ID, *port)
	if err := httpServer.ListenAndServeTLS(appInstance.Config.MTLSCertPath, appInstance.Config.MTLSKeyPath); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
