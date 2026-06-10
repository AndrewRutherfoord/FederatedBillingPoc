package bpclient

import (
	"fmt"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

// Interface for interacting with a CSP's billing adapter. This can allow the transport layer to be swapped out in principle
type BPClient interface {
	// GetBillingRecords(ctx context.Context, from, to time.Time) ([]BillingRecord, error)
	// ... other operations
}

// Implementation of BPClient that uses HTTP with a
type httpBPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPBPClient(baseURL string, client *http.Client) BPClient {
	return &httpBPClient{
		client:  client,
		baseURL: baseURL,
	}
}

type BPClientRegistry struct {
	bps map[string]BPClient
}

func NewBPClientRegistry(config *config.CspConfig) (*BPClientRegistry, error) {
	cspCertPaths := make([]string, len(config.BillingProviders))
	for i, csp := range config.BillingProviders {
		cspCertPaths[i] = csp.MTLSCertPath
	}
	httpClient, err := shared.NewMtlsHttpClient(config.MTLSKeyPath, config.MTLSCertPath, cspCertPaths)

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	bps := make(map[string]BPClient)
	for _, cspConfig := range config.BillingProviders {
		bps[cspConfig.ID] = NewHTTPBPClient(cspConfig.ApiEndpoint, httpClient)
	}

	return &BPClientRegistry{bps: bps}, nil
}
