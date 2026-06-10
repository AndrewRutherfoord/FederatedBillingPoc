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
	"context"
	"flag"
	"log"
	"os"

	_ "github.com/andrewrutherfoord/fed-bill-poc/billing-provider/billing-api/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/billing-api/internal/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/config"
	cspclient "github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/csp_client"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/repository"
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

	// Print loaded config for verification
	log.Printf("Loaded config: ProviderID=%s, ProviderName=%s", cfg.ProviderID, cfg.ProviderName)

	dbPath := os.Getenv("CSP_DB_PATH")
	if dbPath == "" {
		dbPath = "billing-provider.sqlite"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	repos := repository.New(cfg, database)

	mockClock := shared.NewMockClock("http://localhost:9999", []shared.OnTimeAdvanceCallback{})

	log.Printf("Creating CSP client registry")
	clientRegistries, err := cspclient.NewCSPClientRegistry(cfg)
	if err != nil {
		log.Fatalf("failed to create CSP client registry: %v", err)
	}

	client, err := clientRegistries.GetCSPClient("mock-csp-1")
	if err != nil {
		log.Fatalf("failed to get CSP client: %v", err)
	}
	res, err := client.GetTest(context.TODO())
	if err != nil {
		log.Fatalf("failed to call GetTest on CSP client: %v", err)
	}
	log.Printf("GetTest response from CSP client: %s", res)

	r := gin.Default()
	r.Use(cors.Default())
	handlers.NewServer(repos, mockClock).RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting %s (%s) on :%s", repos.Provider.Name, repos.Provider.ID, port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
