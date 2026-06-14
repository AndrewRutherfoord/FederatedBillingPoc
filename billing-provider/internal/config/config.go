package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type SettlementFrequency string

const (
	SettlementDaily   SettlementFrequency = "daily"
	SettlementWeekly  SettlementFrequency = "weekly"
	SettlementMonthly SettlementFrequency = "monthly"
)

type MTLSConfig struct {
	CommonName string `yaml:"common_name"`
	CertPath   string `yaml:"cert_path"`
}

type CloudServiceProvider struct {
	ID                  string              `yaml:"id" json:"id"`
	Name                string              `yaml:"name" json:"name"`
	APIEndpointURL      string              `yaml:"api_endpoint_url" json:"api_endpoint_url"`
	CustomerEndpointURL string              `yaml:"customer_endpoint_url" json:"customer_endpoint_url"`
	SettlementFrequency SettlementFrequency `yaml:"settlement_frequency" json:"settlement_frequency"`
	MTLS                MTLSConfig          `yaml:"mtls" json:"mtls"`
}

type Config struct {
	ProviderID            string                 `yaml:"provider_id" json:"provider_id"`
	ProviderName          string                 `yaml:"provider_name" json:"provider_name"`
	CloudServiceProviders []CloudServiceProvider `yaml:"cloud_service_providers" json:"cloud_service_providers"`
	MTLSKeyPath           string                 `yaml:"mtls_key_path" json:"mtls_key_path"`
	MTLSCertPath          string                 `yaml:"mtls_cert_path" json:"mtls_cert_path"`
}

func (s SettlementFrequency) settlementFrequencyIsValid() bool {
	switch s {
	case SettlementDaily, SettlementWeekly, SettlementMonthly:
		return true
	default:
		return false
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}

	// Validate settlement frequencies
	for _, csp := range cfg.CloudServiceProviders {
		if !csp.SettlementFrequency.settlementFrequencyIsValid() {
			return nil, fmt.Errorf("invalid settlement frequency %q for CSP %q", csp.SettlementFrequency, csp.ID)
		}
	}
	return &cfg, nil
}
