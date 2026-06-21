-- name: GetCloudServiceProvider :one
SELECT id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url FROM cloud_service_providers
WHERE id = ?;

-- name: UpdateCloudServiceProviderCustomerEndpoint :one
UPDATE cloud_service_providers
SET name = ?, customer_api_endpoint_url = ?
WHERE id = ?
RETURNING id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url;

-- name: ListCloudServiceProvidersByBillingProvider :many
SELECT id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url FROM cloud_service_providers
WHERE billing_provider_id = ?;

-- name: DeleteCloudServiceProvidersByBillingProvider :exec
DELETE FROM cloud_service_providers
WHERE billing_provider_id = ?;

-- name: CreateCloudServiceProvider :exec
INSERT INTO cloud_service_providers (id, billing_provider_id, name, api_endpoint_url, customer_api_endpoint_url)
VALUES (?, ?, ?, ?, ?);
