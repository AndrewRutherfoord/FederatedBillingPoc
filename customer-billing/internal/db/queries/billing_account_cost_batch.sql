-- name: ListBillingAccountCostBatches :many
SELECT * FROM billing_account_cost_batch;

-- name: ListBillingAccountCostBatchesByAccount :many
SELECT * FROM billing_account_cost_batch
WHERE billing_account_id = ?;

-- name: GetLatestBillingAccountCostBatchByAccount :one
SELECT * FROM billing_account_cost_batch
WHERE billing_account_id = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateBillingAccountCostBatch :one
INSERT INTO billing_account_cost_batch (id, billing_account_id, billing_period_id, cloud_service_provider_id, total_items, total_cost, billed_currency, merkel_root, batch_signature, created_at, received_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;