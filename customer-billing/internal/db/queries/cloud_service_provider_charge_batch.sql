-- name: ListCloudServiceProviderChargeBatchesByAccount :many
SELECT * FROM cloud_service_provider_charge_batch
WHERE billing_account_id = ?;

-- name: GetLatestCloudServiceProviderChargeBatchByAccount :one
SELECT * FROM cloud_service_provider_charge_batch
WHERE billing_account_id = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateCloudServiceProviderChargeBatch :one
INSERT INTO cloud_service_provider_charge_batch (id, billing_account_id, cloud_service_provider_id, total_items, total_cost, billed_currency, merkle_root, batch_signature, created_at, received_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;
