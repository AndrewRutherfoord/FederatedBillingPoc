package apibpadapter

import (
	"fmt"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

type BPClientRegistry struct {
	bps map[string]shared.HttpClient
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

	bps := make(map[string]shared.HttpClient)
	for _, cspConfig := range config.BillingProviders {
		bps[cspConfig.ID] = *shared.NewHttpClient(cspConfig.ApiEndpoint, httpClient)
	}

	return &BPClientRegistry{bps: bps}, nil
}

func (r *BPClientRegistry) GetClient(bpID string) (shared.HttpClient, bool) {
	client, exists := r.bps[bpID]
	return client, exists
}
