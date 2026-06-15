-- +goose Up
CREATE TABLE billing_providers (
    id TEXT PRIMARY KEY, 
    name TEXT NOT NULL, 
    api_endpoint_url TEXT NOT NULL
);

CREATE TABLE billing_provider_supported_cloud_provider (
    id TEXT PRIMARY KEY, 
    billing_provider_id TEXT NOT NULL,
    name TEXT NOT NULL,
    api_endpoint_url TEXT NOT NULL,
    FOREIGN KEY (billing_provider_id) REFERENCES billing_providers(id)
);

CREATE TABLE billing_accounts (
    id TEXT PRIMARY KEY, 
    billing_provider_id TEXT NOT NULL,
    alias TEXT NOT NULL, 
    created_at DATETIME NOT NULL,
    onboarding_redirect_url TEXT,
    onboarding_complete DATETIME,
    FOREIGN KEY (billing_provider_id) REFERENCES billing_providers(id)
);

CREATE TABLE cloud_service_provider_accounts (
    id TEXT PRIMARY KEY, 
    billing_account_id TEXT NOT NULL,
    cloud_provider_id TEXT NOT NULL,
    FOREIGN KEY (billing_account_id) REFERENCES billing_accounts(id),
    FOREIGN KEY (cloud_provider_id) REFERENCES billing_provider_supported_cloud_provider(id)
);

-- +goose Down
DROP TABLE cloud_service_provider_accounts;
DROP TABLE billing_accounts;
DROP TABLE billing_provider_supported_cloud_provider;
DROP TABLE billing_providers;
