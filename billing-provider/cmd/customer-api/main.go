// Package main is the entry point for the Billing Provider customer-facing API server.
// It exposes billing accounts, invoices, and payment endpoints to authenticated customers.
//
//	@title			Billing Provider – Customer API
//	@version		1.0
//	@description	Customer-facing REST API for the Billing Provider
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8081
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
	apiPort := flag.String("port", ":8081", "port to listen on")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	handler := port.NewCustomerPort(a.Repos)
	adapter := apicspadapter.NewApiCustomerAdapter(handler, a.Config, a.Repos)

	if err := adapter.Start(*apiPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
