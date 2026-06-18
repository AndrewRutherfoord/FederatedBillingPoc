package clients

import (
	"fmt"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	billingprovidermodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
)

type BillingProviderClient interface {
	GetMetadata() (billingprovidermodels.Metadata, error)
	RegisterBillingAccount(returnURL string) (billingprovidermodels.RegisterBillingAccountResponse, error)
	GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.BillingRecord, error)
}

type billingProviderClient struct {
	shared.HttpClient
}

func NewBillingProviderClient(baseURL string) *billingProviderClient {
	return &billingProviderClient{
		HttpClient: shared.NewHttpClient(baseURL),
	}
}

func (c *billingProviderClient) GetMetadata() (billingprovidermodels.Metadata, error) {
	var metadataResponse billingprovidermodels.Metadata
	err := c.FetchJSON("/.well-known/billing-provider", &metadataResponse)
	if err != nil {
		return billingprovidermodels.Metadata{}, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	return metadataResponse, nil
}

func (c *billingProviderClient) RegisterBillingAccount(returnURL string) (billingprovidermodels.RegisterBillingAccountResponse, error) {
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

func (c *billingProviderClient) GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.AggregatedChargeRecord, error) {
	payload := billingprovidermodels.GetBillingAccountRecordsRequest{
		BillingAccountID: billingAccountID,
		From:             from,
		To:               to,
	}

	var response billingprovidermodels.GetBillingAccountRecordsResponse
	err := c.SendJSON("/billing/accounts/records", payload, &response)
	if err != nil {
		return []sharedmodels.AggregatedChargeRecord{}, fmt.Errorf("failed to fetch billing account records: %w", err)
	}

	if len(response.Records) != response.Count {
		return []sharedmodels.AggregatedChargeRecord{}, fmt.Errorf("response count mismatch: expected %d, got %d", response.Count, len(response.Records))
	}

	return response.Records, nil
}
