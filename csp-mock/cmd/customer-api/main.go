// Package main is the entry point for the CSP Mock Customer API server.
//
//	@title			CSP Mock Customer API
//	@version		1.0
//	@description	Mock Cloud Service Provider customer billing API implementing the FOCUS spec.
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
	"net/http"

	_ "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/app"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	port := flag.String("port", ":8080", "port to listen on")
	flag.Parse()

	// Initialize app with shared setup
	appInstance := app.NewApp(app.Config{
		ConfigPath: *configPath,
	})

	r := gin.Default()
	r.Use(cors.Default())

	// Customer routes
	server := handlers.NewServer(appInstance.Config, appInstance.Repos, appInstance.Clock)
	server.RegisterRoutes(r)

	httpServer := &http.Server{
		Addr:    *port,
		Handler: r,
	}

	log.Printf("Starting %s (%s) customer API on %s", appInstance.Repos.Provider.Name, appInstance.Repos.Provider.ID, *port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
