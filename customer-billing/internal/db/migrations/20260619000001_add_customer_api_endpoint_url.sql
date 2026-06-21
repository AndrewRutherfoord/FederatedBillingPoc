-- +goose Up
ALTER TABLE billing_provider_supported_cloud_provider ADD COLUMN customer_api_endpoint_url TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE billing_provider_supported_cloud_provider DROP COLUMN customer_api_endpoint_url;
