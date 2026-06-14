package billingprovider

type SupportedCloudProviderMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	APIEndpoint string `json:"api_endpoint"` // The API endpoint that the customer uses to interact with the cloud provider
}

type Metadata struct {
	ID                      string                           `json:"id"`
	Name                    string                           `json:"name"`
	SupportedCloudProviders []SupportedCloudProviderMetadata `json:"supported_cloud_providers"`
}
