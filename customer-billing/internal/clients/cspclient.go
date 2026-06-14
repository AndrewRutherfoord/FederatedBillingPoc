package clients

import (
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


func (c *cloudServiceProviderClient) GetMetadata() (map[string]string, error) {


func (c *cloudServiceProviderClient) GetMetadata() (cspsharedmodels.Metadata, error) {
	var metadataResponse cspsharedmodels.Metadata
	err := c.FetchJSON("/.well-known/cloud-service-provider", &metadataResponse)
	if err != nil {
		return cspsharedmodels.Metadata{}, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	return metadataResponse, nil
}
}