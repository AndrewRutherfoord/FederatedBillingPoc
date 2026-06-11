package cspclient

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

// Interface for interacting with a CSP's billing adapter. This can allow the transport layer to be swapped out in principle
type CSPClient interface {
	// GetBillingRecords(ctx context.Context, from, to time.Time) ([]BillingRecord, error)
	// ... other operations
	GetTest(ctx context.Context) (string, error)
	// GetCostRecords(ctx context.Context, from, to string) ([]sharedmodels.CostBatchMessage, error)
}

type httpCSPClient struct {
	*shared.HttpClient
}

func NewHTTPCSPClient(baseURL string, client *http.Client) CSPClient {
	return &httpCSPClient{
		HttpClient: &shared.HttpClient{
			Client:  client,
			BaseURL: baseURL,
		},
	}
}

func (r *httpCSPClient) GetTest(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/test", r.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.Client.Do(req)
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

// func (r *httpCSPClient) GetCostRecords(ctx context.Context, from, to string) ([]sharedmodels.CostBatchMessage, error) {
// 	resp, err := r.makeRequest(ctx, "GET", fmt.Sprintf("cost-records?from=%s&to=%s", from, to))
// 	if err != nil {
// 		return nil, err
// 	}

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %w", err)
// 	}

// 	log.Printf("Received cost records: %s", string(body))

// 	// In a real implementation, we would unmarshal the body into the appropriate struct
// 	// For this example, we'll just return an empty slice
// 	return []sharedmodels.CostBatchMessage{}, nil
// }

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

func NewCSPClientRegistry(config *config.Config) (*CSPClientRegistry, error) {
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
