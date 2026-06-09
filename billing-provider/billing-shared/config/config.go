package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type CloudServiceProvider struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	APIEndpoint string `yaml:"api_endpoint" json:"api_endpoint"`
}

type BillingProviderConfig struct {
	ProviderID            string                 `yaml:"provider_id" json:"provider_id"`
	ProviderName          string                 `yaml:"provider_name" json:"provider_name"`
	CloudServiceProviders []CloudServiceProvider `yaml:"cloud_service_providers" json:"cloud_service_providers"`
}

func Load(path string) (*BillingProviderConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %q: %w", path, err)
	}
	var cfg BillingProviderConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %q: %w", path, err)
	}
	return &cfg, nil
}
