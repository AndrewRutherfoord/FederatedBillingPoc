package clients

import (
	"fmt"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
)

type CloudServiceProviderClient interface {
	GetMetadata() (cspsharedmodels.Metadata, error)
}

type cloudServiceProviderClient struct {
	shared.HttpClient
}

func NewCloudServiceProviderClient(baseURL string) CloudServiceProviderClient {
	return &cloudServiceProviderClient{
		HttpClient: shared.NewHttpClient(baseURL),
	}
}

func (c *cloudServiceProviderClient) GetMetadata() (cspsharedmodels.Metadata, error) {
	var metadataResponse cspsharedmodels.Metadata
	err := c.FetchJSON("/.well-known/cloud-service-provider", &metadataResponse)
	if err != nil {
		return cspsharedmodels.Metadata{}, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	return metadataResponse, nil
}

func (c *cloudServiceProviderClient) RegisterCloudProviderLink(billingProviderID string, billingAccountID string, returnURL string) (cspsharedmodels.RegisterLinkedCloudProviderResponse, error) {
	payload := cspsharedmodels.RegisterLinkedCloudProviderRequest{
		BillingProviderID: billingProviderID,
		BillingAccountID:  billingAccountID,
		ReturnURL:         returnURL,
	}

	var response cspsharedmodels.RegisterLinkedCloudProviderResponse
	err := c.SendJSON("/billing/accounts", payload, &response)
	if err != nil {
		return cspsharedmodels.RegisterLinkedCloudProviderResponse{}, fmt.Errorf("failed to register linked cloud provider: %w", err)
	}
	return response, nil
}
