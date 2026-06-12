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

	_ "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/app"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/handlers"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	port := flag.String("port", ":8080", "port to listen on")
	flag.Parse()

	appInstance := app.NewApp(app.Config{ConfigPath: *configPath})

	server := handlers.NewServer(appInstance.Config, appInstance.Repos, appInstance.Clock)
	if err := server.Start(*port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
