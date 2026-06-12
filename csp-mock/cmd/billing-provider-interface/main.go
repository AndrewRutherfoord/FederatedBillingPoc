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
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/port"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	sharedscheduler "github.com/andrewrutherfoord/fed-bill-poc/shared/scheduler"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	apiPort := flag.String("port", ":8443", "port to listen on")
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

	handler := port.NewBPPort(appInstance.Repos)

	// Adapter is the inbound HTTP server (receives from BPs) and outbound sender (sends to BPs).
	// Pass handler to adapter so incoming BP messages are routed to application logic.
	// Pass adapter as port.BPSender to the scheduler when wiring it up.
	adapter := apibpadapter.NewApiBillingProviderAdapter(r, handler, appInstance.Config)
	defer adapter.Close()

	costBatchJob := scheduler.NewRecordMeteringAndCostJob("cost-batch", appInstance.Repos, appInstance.Config, adapter)
	costBatchSched, err := sharedscheduler.NewCronSchedule("0 0 * * * *") // every hour on the hour
	if err != nil {
		log.Fatalf("failed to create cost batch schedule: %v", err)
	}
	if err := appInstance.Sched.Register(costBatchJob, costBatchSched); err != nil {
		log.Fatalf("failed to register cost batch job: %v", err)
	}

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  trustedPool,
	}

	httpServer := &http.Server{
		Addr:      *apiPort,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting %s (%s) adapter server on %s", appInstance.Repos.Provider.Name, appInstance.Repos.Provider.ID, *apiPort)
	if err := httpServer.ListenAndServeTLS(appInstance.Config.MTLSCertPath, appInstance.Config.MTLSKeyPath); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
