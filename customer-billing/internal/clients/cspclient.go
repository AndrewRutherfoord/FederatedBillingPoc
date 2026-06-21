package clients

import (
	"fmt"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	cspsharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/cloud-service-provider"
)

type CloudServiceProviderClient interface {
	GetMetadata() (cspsharedmodels.Metadata, error)
	RegisterCloudProviderAccount(billingProviderID string, billingAccountID string, returnURL string) (cspsharedmodels.RegisterLinkedCloudProviderResponse, error)
	GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.ChargeBatchDetail, error)
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
	err := c.SendJSON("/billing/accounts", payload, &response)
	if err != nil {
		return cspsharedmodels.RegisterLinkedCloudProviderResponse{}, fmt.Errorf("failed to register linked cloud provider: %w", err)
	}
	return response, nil
}

// GetBillingAccountRecords fetches this CSP's own report of a billing account's charge
// batches, including the raw FOCUS line items, independently of the billing provider.
func (c *cloudServiceProviderClient) GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.ChargeBatchDetail, error) {
	payload := cspsharedmodels.GetBillingAccountRecordsRequest{
		BillingAccountID: billingAccountID,
		From:             from,
		To:               to,
	}

	var response cspsharedmodels.GetBillingAccountRecordsResponse
	err := c.SendJSON("/billing/accounts/records", payload, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch billing account records: %w", err)
	}

	if len(response.Batches) != response.Count {
		return nil, fmt.Errorf("response count mismatch: expected %d, got %d", response.Count, len(response.Batches))
	}

	return response.Batches, nil
}
