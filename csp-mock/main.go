// Package main is the entry point for the CSP Mock API server.
//
//	@title			CSP Mock API
//	@version		1.0
//	@description	Mock Cloud Service Provider billing API implementing the FOCUS spec.
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer <account_id>
package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/util"
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
	mockClock := util.NewMockClock(1767225600, repos.KeyValue) // Jan 1, 2026 00:00:00 UTC

	r := gin.Default()
	r.Use(cors.Default())
	handlers.NewServer(repos, mockClock).RegisterRoutes(r)

	log.Printf("Starting %s (%s) on :8080", repos.Provider.Name, repos.Provider.ID)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
