-- name: CreateCloudServiceProviderAccount :one
INSERT INTO cloud_service_provider_accounts (id, billing_account_id, cloud_provider_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetCloudServiceProviderAccount :one
SELECT * FROM cloud_service_provider_accounts
WHERE id = ?;

-- name: ListCloudServiceProviderAccountsByBillingAccount :many
SELECT cspa.id, cspa.billing_account_id, cspa.cloud_provider_id, bpsc.name AS cloud_provider_name
FROM cloud_service_provider_accounts cspa
JOIN billing_provider_supported_cloud_provider bpsc ON bpsc.id = cspa.cloud_provider_id
WHERE cspa.billing_account_id = ?;

-- name: ListCloudServiceProviderAccountsByBillingAccountWithTotalCost :many
SELECT cspa.id, cspa.billing_account_id, cspa.cloud_provider_id, bpsc.name AS cloud_provider_name, SUM(bacb.total_cost) AS total_cost, bacb.billed_currency AS billing_currency
FROM cloud_service_provider_accounts cspa
JOIN billing_provider_supported_cloud_provider bpsc ON bpsc.id = cspa.cloud_provider_id
LEFT JOIN billing_account_cost_batch bacb ON bacb.cloud_service_provider_id = cspa.cloud_provider_id AND bacb.billing_account_id = cspa.billing_account_id
WHERE cspa.billing_account_id = ?
GROUP BY cspa.id, cspa.billing_account_id, cspa.cloud_provider_id, bpsc.name, bacb.billed_currency;