package apicspadapter

import (
	"fmt"
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

type CSPClientRegistry struct {
	csps map[string]shared.HttpClient
}

func NewCSPClientRegistry(cfg *config.Config) (*CSPClientRegistry, error) {
	certPaths := make([]string, len(cfg.CloudServiceProviders))
	for i, csp := range cfg.CloudServiceProviders {
		certPaths[i] = csp.MTLS.CertPath
	}

	httpClient, err := shared.NewMtlsHttpClient(cfg.MTLSKeyPath, cfg.MTLSCertPath, certPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to create mTLS HTTP client: %w", err)
	}

	csps := make(map[string]shared.HttpClient)
	for _, csp := range cfg.CloudServiceProviders {
		csps[csp.ID] = shared.NewHttpClientWithCustomClient(csp.APIEndpointURL, httpClient)
	}

	log.Printf("Created CSP client registry for: %v", cfg.CloudServiceProviders)

	return &CSPClientRegistry{csps: csps}, nil
}

func (r *CSPClientRegistry) GetClient(cspID string) (shared.HttpClient, bool) {
	client, exists := r.csps[cspID]
	return client, exists
}
