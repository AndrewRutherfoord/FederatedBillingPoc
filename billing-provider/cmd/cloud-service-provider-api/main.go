// Package main is the entry point for the Billing Provider CSP-facing API server.
// It receives cost records pushed by CSPs over mTLS and exposes CSP management endpoints.
//
//	@title			Billing Provider – CSP API
//	@version		1.0
//	@description	Inbound mTLS API for Cloud Service Provider interactions
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8444
//	@BasePath	/
package main

import (
	"flag"
	"log"

	apicspadapter "github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/adapters/api"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/app"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/port"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	apiPort := flag.String("port", ":8444", "port to listen on")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	handler := port.NewCSPPort(a.Repos, a.Clock)
	adapter := apicspadapter.NewApiCSPAdapter(handler, a.Config, a.Repos)
	defer adapter.Close()

	if err := adapter.Start(*apiPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
