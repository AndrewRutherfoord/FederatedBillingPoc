package clients

import "github.com/andrewrutherfoord/fed-bill-poc/shared"

type CloudServiceProviderClient interface {
	GetMetadata() (map[string]string, error)
}

type cloudServiceProviderClient struct {
	shared.HttpClient
}

func NewCloudServiceProviderClient(baseURL string) *cloudServiceProviderClient {
	return &cloudServiceProviderClient{
		HttpClient: shared.NewHttpClient(baseURL),
	}
}
