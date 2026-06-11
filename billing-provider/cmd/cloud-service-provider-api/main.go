// Package main is the entry point for the Billing Provider API server
//
//	@title			Billing Provider API
//	@version		1.0
//	@description	Implementation of a Billing Provider
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8080
//	@BasePath	/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/api/cloud-service-provider/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Build the trusted cert pool for CSPs (billing provider certs)
	trustedPool := x509.NewCertPool()
	for _, csp := range a.Config.CloudServiceProviders {
		trustedCertPEM, err := os.ReadFile(csp.MTLS.CertPath)
		if err != nil {
			log.Fatalf("failed to read cert for billing provider %s: %v", csp.Name, err)
		}
		if !trustedPool.AppendCertsFromPEM(trustedCertPEM) {
			log.Fatalf("failed to append cert for billing provider %s", csp.Name)
		}
	}

	handlers.NewServer(a.Repos, a.Clock).RegisterRoutes(r)

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  trustedPool,
	}

	httpServer := &http.Server{
		Addr:      ":8444",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting %s (%s) on :8443", a.Config.ProviderName, a.Config.ProviderID)
	if err := httpServer.ListenAndServeTLS(a.Config.MTLSCertPath, a.Config.MTLSKeyPath); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
