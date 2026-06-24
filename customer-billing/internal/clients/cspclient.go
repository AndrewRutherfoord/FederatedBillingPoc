package clients

import (
	"fmt"
	"net/url"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
)

type CloudServiceProviderClient interface {
	GetMetadata() (cspsharedmodels.Metadata, error)
	RegisterCloudProviderAccount(billingProviderID string, billingAccountID string, returnURL string) (cspsharedmodels.RegisterLinkedCloudProviderResponse, error)
	GetChargeBatch(batchID string, billingAccountID string) (sharedmodels.ChargeBatchDetail, error)
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

func (c *cloudServiceProviderClient) RegisterCloudProviderAccount(billingProviderID string, billingAccountID string, returnURL string) (cspsharedmodels.RegisterLinkedCloudProviderResponse, error) {
	payload := cspsharedmodels.RegisterLinkedCloudProviderRequest{
		BillingProviderID: billingProviderID,
		BillingAccountID:  billingAccountID,
		ReturnURL:         returnURL,
	}

	var response cspsharedmodels.RegisterLinkedCloudProviderResponse
	err := c.SendJSON("/api/v1/billing/onboarding", payload, &response)
	if err != nil {
		return cspsharedmodels.RegisterLinkedCloudProviderResponse{}, fmt.Errorf("failed to register linked cloud provider: %w", err)
	}
	return response, nil
}

// GetChargeBatch fetches CSP report of a charge batch. Includes raw FOCUS line items. Must already know the Batch ID after getting it from the BP
func (c *cloudServiceProviderClient) GetChargeBatch(batchID string, billingAccountID string) (sharedmodels.ChargeBatchDetail, error) {
	var response sharedmodels.ChargeBatchDetail

	err := c.FetchJSONWithAuth("/api/v1/billing/charge-batches/"+url.PathEscape(batchID), billingAccountID, &response)
	if err != nil {
		return sharedmodels.ChargeBatchDetail{}, fmt.Errorf("failed to fetch charge batch %s: %w", batchID, err)
	}

	return response, nil
}
