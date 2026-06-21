-- name: CreateCloudServiceProviderAccount :one
INSERT INTO cloud_service_provider_accounts (id, billing_account_id, cloud_service_provider_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetCloudServiceProviderAccount :one
SELECT * FROM cloud_service_provider_accounts
WHERE id = ?;

-- name: ListCloudServiceProviderAccountsByBillingAccountWithTotalCost :many
SELECT cspa.id, cspa.billing_account_id, cspa.cloud_service_provider_id, csp.name AS cloud_service_provider_name, csp.customer_api_endpoint_url AS customer_api_endpoint_url, SUM(bacb.total_cost) AS total_cost, bacb.billed_currency AS billing_currency
FROM cloud_service_provider_accounts cspa
JOIN cloud_service_providers csp ON csp.id = cspa.cloud_service_provider_id
LEFT JOIN billing_account_charge_batch bacb ON bacb.cloud_service_provider_id = cspa.cloud_service_provider_id AND bacb.billing_account_id = cspa.billing_account_id
WHERE cspa.billing_account_id = ?
GROUP BY cspa.id, cspa.billing_account_id, cspa.cloud_service_provider_id, csp.name, csp.customer_api_endpoint_url, bacb.billed_currency;
