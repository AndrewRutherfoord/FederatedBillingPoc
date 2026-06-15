package cloudserviceprovider

type SupportedCloudProviderMetadata struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Metadata struct {
	ID                        string                           `json:"id"`
	Name                      string                           `json:"name"`
	SupportedBillingProviders []SupportedCloudProviderMetadata `json:"supported_billing_providers"`
}
