package ingestpipeline

import (
	"context"
	"flag"
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/app"
	cspclient "github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/csp_client"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	clientRegistry, err := cspclient.NewCSPClientRegistry(a.Config)
	if err != nil {
		log.Fatalf("failed to create CSP client registry: %v", err)
	}

	ctx := context.Background()
	csps, err := a.Repos.CloudServiceProviders.List(ctx)
	if err != nil {
		log.Fatalf("failed to list cloud service providers: %v", err)
	}
	for _, csp := range csps {
		log.Printf("CSP: %s, API Endpoint: %s", csp.Name, csp.APIEndpoint)
		client, err := clientRegistry.GetCSPClient(csp.ID)
		if err != nil {
			log.Printf("failed to get client for CSP %s: %v", csp.Name, err)
			continue
		}
		// Example usage of the client - replace with actual logic as needed
		err = client.GetTest(ctx)
		if err != nil {
			log.Printf("failed to connect to CSP %s: %v", csp.Name, err)
		} else {
			log.Printf("successfully connected to CSP %s", csp.Name)
		}
	}

}
