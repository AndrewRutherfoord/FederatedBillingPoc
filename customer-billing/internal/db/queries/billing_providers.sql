-- name: GetBillingProvider :one
SELECT * FROM billing_providers
WHERE id = ?;

-- name: CreateBillingProvider :exec
INSERT INTO billing_providers (id, name, api_endpoint_url)
VALUES (?, ?, ?);

-- name: UpsertBillingProvider :exec
INSERT INTO billing_providers (id, name, api_endpoint_url)
VALUES (?, ?, ?)
ON CONFLICT(id) DO UPDATE SET name = excluded.name, api_endpoint_url = excluded.api_endpoint_url;

-- name: ListBillingProviders :many
SELECT * FROM billing_providers;

-- name: GetBillingProviderSupportedCSP :one
SELECT id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url FROM billing_provider_supported_cloud_provider
WHERE id = ?;

-- name: UpdateBillingProviderSupportedCSPCustomerEndpoint :one
UPDATE billing_provider_supported_cloud_provider
SET name = ?, customer_api_endpoint_url = ?
WHERE id = ?
RETURNING id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url;

-- name: ListBillingProviderSupportedCSPs :many
SELECT id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url FROM billing_provider_supported_cloud_provider
WHERE billing_provider_id = ?;

-- name: DeleteBillingProviderSupportedCSPs :exec
DELETE FROM billing_provider_supported_cloud_provider
WHERE billing_provider_id = ?;

-- name: CreateBillingProviderSupportedCSP :exec
INSERT INTO billing_provider_supported_cloud_provider (id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url)
VALUES (?, ?, ?, ?, ?);