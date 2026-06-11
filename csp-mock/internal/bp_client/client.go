package bpclient

import (
	"fmt"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

// Interface for interacting with a CSP's billing adapter. This can allow the transport layer to be swapped out in principle
type BPClient interface {
	SendAggregatedChargeRecord(record sharedmodels.AggregatedChargeRecord) error
	// GetBillingRecords(ctx context.Context, from, to time.Time) ([]BillingRecord, error)
	// ... other operations
}

type HttpBPClient struct {
	*shared.HttpClient
}

func NewHTTPBPClient(baseURL string, client *http.Client) BPClient {
	return &HttpBPClient{
		HttpClient: &shared.HttpClient{
			Client:  client,
			BaseURL: baseURL,
		},
	}
}

func (c *HttpBPClient) SendAggregatedChargeRecord(record sharedmodels.AggregatedChargeRecord) error {
	url := fmt.Sprintf("%s/aggregated-charge-records", c.BaseURL)
	return c.SendJSON(url, record)
}

type BPClientRegistry struct {
	bps map[string]BPClient
}

func NewBPClientRegistry(config *config.CspConfig) (*BPClientRegistry, error) {
	cspCertPaths := make([]string, len(config.BillingProviders))
	for i, csp := range config.BillingProviders {
		cspCertPaths[i] = csp.MTLS.CertPath
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
