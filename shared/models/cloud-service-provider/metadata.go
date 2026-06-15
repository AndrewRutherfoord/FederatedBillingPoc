package cloudserviceprovider

type SupportedCloudProviderMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Metadata struct {
	ID                        string                           `json:"id"`
	Name                      string                           `json:"name"`
	APIEndpointURL            string                           `json:"api_endpoint_url"`
	SupportedBillingProviders []SupportedCloudProviderMetadata `json:"supported_billing_providers"`
}
