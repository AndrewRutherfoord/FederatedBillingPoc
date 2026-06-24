package clients

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	billingprovidermodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
)

type BillingProviderClient interface {
	GetMetadata() (billingprovidermodels.Metadata, error)
	RegisterBillingAccount(returnURL string) (billingprovidermodels.RegisterBillingAccountResponse, error)
	GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.ChargeBatch, error)
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

func (c *billingProviderClient) GetBillingAccountRecords(billingAccountID string, from time.Time, to time.Time) ([]sharedmodels.ChargeBatch, error) {
	path := fmt.Sprintf("/billing/accounts/charge-batches?from=%s&to=%s",
		url.QueryEscape(from.Format(time.RFC3339)),
		url.QueryEscape(to.Format(time.RFC3339)),
	)

	var response billingprovidermodels.GetBillingAccountRecordsResponse
	err := c.FetchJSONWithAuth(path, billingAccountID, &response)
	log.Printf("Fetched %d charge batches for billing account %s from %s to %s", len(response.Batches), billingAccountID, from.Format(time.RFC3339), to.Format(time.RFC3339))
	if err != nil {
		return []sharedmodels.ChargeBatch{}, fmt.Errorf("failed to fetch billing account records: %w", err)
	}

	if len(response.Batches) != response.Count {
		return []sharedmodels.ChargeBatch{}, fmt.Errorf("response count mismatch: expected %d, got %d", response.Count, len(response.Batches))
	}

	return response.Batches, nil
}
