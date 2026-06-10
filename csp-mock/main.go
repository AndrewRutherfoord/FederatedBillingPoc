// Package main is the entry point for the CSP Mock API server.
//
//	@title			CSP Mock API
//	@version		1.0
//	@description	Mock Cloud Service Provider billing API implementing the FOCUS spec.
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8443
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer <account_id>
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers/billingproviderhandlers"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers/customerhandlers"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbPath := os.Getenv("CSP_DB_PATH")
	if dbPath == "" {
		dbPath = "csp-mock.sqlite"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	repos := repository.New(cfg, database)

	sched := scheduler.NewWithPersistence(repos.KeyValue)
	err = scheduler.RegisterJobs(sched, []scheduler.JobToRegister{
		scheduler.NewJobToRegister(
			scheduler.NewRecordMeteringAndCostJob("record-metering-and-cost", repos, cfg),
			"0 0 * * * *", // Every hour at :00 seconds
		),
	})
	if err != nil {
		log.Fatalf("failed to register jobs: %v", err)
	}

	clockHost := os.Getenv("MOCK_CLOCK_HOST")
	if clockHost == "" {
		clockHost = "http://localhost:9999"
	}

	// Mock clock that allows manual time advancement for testing. It uses a centralised mock time server. Later it can be swapped out for a regular clock that just returns the current time.
	// On the clock advancing it calls the scheduler's OnTimeAdvance method which triggers any jobs
	clock := shared.NewMockClock(clockHost, []shared.OnTimeAdvanceCallback{sched})

	r := gin.Default()
	r.Use(cors.Default())

	// Build the trusted cert pool for CSPs (billing provider certs)
	trustedPool := x509.NewCertPool()
	for _, bp := range cfg.BillingProviders {
		trustedCertPEM, err := os.ReadFile(bp.MTLS.CertPath)
		if err != nil {
			log.Fatalf("failed to read cert for billing provider %s: %v", bp.Name, err)
		}
		if !trustedPool.AppendCertsFromPEM(trustedCertPEM) {
			log.Fatalf("failed to append cert for billing provider %s", bp.Name)
		}
	}

	server := handlers.NewServer(cfg)
	server.RegisterRoutes(r, []handlers.SubServer{
		customerhandlers.NewCustomerServer(repos, clock),
		billingproviderhandlers.NewBillingProviderServer(cfg.BillingProviders, repos, clock),
	})

	tlsConfig := &tls.Config{
		ClientAuth: tls.RequestClientCert,
		ClientCAs:  trustedPool,
	}

	httpServer := &http.Server{
		Addr:      ":8443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting %s (%s) on :8443", repos.Provider.Name, repos.Provider.ID)
	if err := httpServer.ListenAndServeTLS(cfg.MTLSCertPath, cfg.MTLSKeyPath); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
