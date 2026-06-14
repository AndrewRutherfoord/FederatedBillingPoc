package clients

import (
	"fmt"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	billingprovidermodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
)

type BillingProviderClient interface {
	RegisterBillingAccount(returnURL string) (billingprovidermodels.RegisterBillingAccountResponse, error)
}

type BillingProviderClientImpl struct {
	shared.HttpClient
}

func NewBillingProviderClient(baseURL string) *BillingProviderClientImpl {
	return &BillingProviderClientImpl{
		HttpClient: shared.NewHttpClient(baseURL),
	}
}

func (c *BillingProviderClientImpl) GetMetadata() (billingprovidermodels.Metadata, error) {
	var metadataResponse billingprovidermodels.Metadata
	err := c.FetchJSON("/.well-known/billing-provider", &metadataResponse)
	if err != nil {
		return billingprovidermodels.Metadata{}, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	return metadataResponse, nil
}

func (c *BillingProviderClientImpl) RegisterBillingAccount(returnURL string) (billingprovidermodels.RegisterBillingAccountResponse, error) {
	payload := billingprovidermodels.RegisterBillingAccountRequest{
		ReturnURL: returnURL,
	}

	var response billingprovidermodels.RegisterBillingAccountResponse
	err := c.SendJSON("/billing/accounts", payload, &response)
	if err != nil {
		return billingprovidermodels.RegisterBillingAccountResponse{}, fmt.Errorf("failed to register billing account: %w", err)
	}
	return response, nil
}
