package cspclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/shared/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

// Interface for interacting with a CSP's billing adapter. This can allow the transport layer to be swapped out in principle
type CSPClient interface {
	// GetBillingRecords(ctx context.Context, from, to time.Time) ([]BillingRecord, error)
	// ... other operations
	GetTest(ctx context.Context) (string, error)
}

// Implementation of CSPClient that uses HTTP with a
type httpCSPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPCSPClient(baseURL string, client *http.Client) CSPClient {
	return &httpCSPClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (r *httpCSPClient) GetTest(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/test", r.baseURL), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Printf("Received message: %s", string(body))

	return string(body), nil
}

type CSPClientRegistry struct {
	csps map[string]CSPClient
}

func (r *CSPClientRegistry) GetCSPClient(id string) (CSPClient, error) {
	csp, ok := r.csps[id]
	if !ok {
		return nil, fmt.Errorf("CSP client not found for id: %s", id)
	}
	return csp, nil
}

func NewCSPClientRegistry(config *config.BillingProviderConfig) (*CSPClientRegistry, error) {
	cspCertPaths := make([]string, len(config.CloudServiceProviders))
	for i, csp := range config.CloudServiceProviders {
		cspCertPaths[i] = csp.MTLSCertPath
	}
	httpClient, err := shared.NewMtlsHttpClient(config.MTLSKeyPath, config.MTLSCertPath, cspCertPaths)

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	csps := make(map[string]CSPClient)
	for _, cspConfig := range config.CloudServiceProviders {
		csps[cspConfig.ID] = NewHTTPCSPClient(cspConfig.APIEndpoint, httpClient)
	}

	return &CSPClientRegistry{csps: csps}, nil
}
