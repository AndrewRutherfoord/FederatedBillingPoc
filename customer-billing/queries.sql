-- name: ListBillingProviders :many
SELECT id, name, api_endpoint_url FROM billing_providers;

-- name: GetBillingProvider :one
SELECT id, name, api_endpoint_url FROM billing_providers WHERE id = ?;

-- name: ListCloudProvidersByBillingProvider :many
SELECT id, billing_provider_id, name, api_endpoint_url FROM billing_provider_supported_cloud_provider WHERE billing_provider_id = ?;