package config

import (
	"fmt"
	"os"
	"strings"

	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v3"
)

// normalizeLabel converts human-readable FOCUS labels (e.g. "Virtual Machines")
// to their snake_case enum values (e.g. "virtual_machines").
func normalizeLabel(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

type PricingTier struct {
	UpToGB    *int            `yaml:"up_to_gb"`
	UnitPrice decimal.Decimal `yaml:"unit_price"`
}

type PricingConfig struct {
	Model     string          `yaml:"model"`
	UnitPrice decimal.Decimal `yaml:"unit_price"`
	Tiers     []PricingTier   `yaml:"tiers,omitempty"`
}

type ResourceType struct {
	ID                 string                          `yaml:"id"`
	DisplayName        string                          `yaml:"display_name"`
	Description        string                          `yaml:"description,omitempty"`
	ServiceName        string                          `yaml:"service_name"`
	ServiceCategory    sharedmodels.ServiceCategory    `yaml:"service_category"`
	ServiceSubcategory sharedmodels.ServiceSubcategory `yaml:"service_subcategory"`
	SkuID              string                          `yaml:"sku_id"`
	BillingUnit        string                          `yaml:"billing_unit"`
	Pricing            PricingConfig                   `yaml:"pricing"`
	ConfigSchema       map[string]any                  `yaml:"config_schema,omitempty"`
	Tags               map[string]string               `yaml:"tags,omitempty"`
}

// UnmarshalYAML normalises service_category and service_subcategory from
// human-readable title-case (e.g. "Virtual Machines") to snake_case enum values.
func (r *ResourceType) UnmarshalYAML(value *yaml.Node) error {
	// Alias avoids infinite recursion when decoding into the same type.
	type raw ResourceType
	var tmp raw
	if err := value.Decode(&tmp); err != nil {
		return err
	}
	*r = ResourceType(tmp)
	r.ServiceCategory = sharedmodels.ServiceCategory(normalizeLabel(string(r.ServiceCategory)))
	r.ServiceSubcategory = sharedmodels.ServiceSubcategory(normalizeLabel(string(r.ServiceSubcategory)))
	return nil
}

type MTLSConfig struct {
	CommonName string `yaml:"common_name"`
	CertPath   string `yaml:"cert_path"`
}

type BillingProvider struct {
	ID          string     `yaml:"id"`
	Name        string     `yaml:"name"`
	ApiEndpoint string     `yaml:"api_endpoint"`
	MTLS        MTLSConfig `yaml:"mtls"`
}

// TODO: Remove this since it's not being used
type MeteringConfig struct {
	TickIntervalSeconds   int `yaml:"tick_interval_seconds"`
	SimulatedHoursPerTick int `yaml:"simulated_hours_per_tick"`
}

// CspConfig is the root structure for the mock CSP config.yaml.
type CspConfig struct {
	ProviderID             string            `yaml:"provider_id"`
	ProviderName           string            `yaml:"provider_name"`
	Currency               string            `yaml:"currency"`
	RegionID               string            `yaml:"region_id"`
	RegionName             string            `yaml:"region_name"`
	AvailabilityZone       string            `yaml:"availability_zone"`
	CustomerAPIEndpointURL string            `yaml:"customer_api_endpoint_url"` // Customer API endpoint URL that customer metering seric calls
	ResourceTypes          []ResourceType    `yaml:"resource_types"`
	BillingProviders       []BillingProvider `yaml:"billing_providers"`
	Metering               MeteringConfig    `yaml:"metering"`
	MTLSKeyPath            string            `yaml:"mtls_key_path"`
	MTLSCertPath           string            `yaml:"mtls_cert_path"`
}

// Load reads and parses the YAML config at path.
func Load(path string) (*CspConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}
	var cfg CspConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}
	return &cfg, nil
}
