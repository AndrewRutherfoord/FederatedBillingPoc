-- name: ListBillingAccountChargeBatches :many
SELECT * FROM billing_account_charge_batch;

-- name: ListBillingAccountChargeBatchesByAccount :many
SELECT * FROM billing_account_charge_batch
WHERE billing_account_id = ?;

-- name: GetLatestBillingAccountChargeBatchByAccount :one
SELECT * FROM billing_account_charge_batch
WHERE billing_account_id = ?
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateBillingAccountChargeBatch :one
INSERT INTO billing_account_charge_batch (id, billing_account_id, billing_period_id, cloud_service_provider_id, total_items, total_cost, billed_currency, merkle_root, batch_signature, created_at, received_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: ListChargeBatchesReportedByBillingProvider :many
-- Every charge batch the billing provider reported for this billing account, with the cloud
-- service provider's own report of the same batch ID alongside it (nullable, since the CSP
-- side may not have been fetched yet, or may simply disagree on whether the batch exists).
SELECT
    bp.id AS batch_id,
    bp.billing_account_id AS billing_account_id,
    bp.cloud_service_provider_id AS cloud_service_provider_id,
    bp.total_items AS bp_total_items,
    bp.total_cost AS bp_total_cost,
    bp.billed_currency AS bp_billed_currency,
    bp.merkle_root AS bp_merkle_root,
    bp.batch_signature AS bp_batch_signature,
    bp.created_at AS bp_created_at,
    bp.received_at AS bp_received_at,
    csp.total_items AS csp_total_items,
    csp.total_cost AS csp_total_cost,
    csp.billed_currency AS csp_billed_currency,
    csp.merkle_root AS csp_merkle_root,
    csp.batch_signature AS csp_batch_signature,
    csp.created_at AS csp_created_at,
    csp.received_at AS csp_received_at
FROM billing_account_charge_batch bp
LEFT JOIN cloud_service_provider_charge_batch csp ON csp.id = bp.id
WHERE bp.billing_account_id = ?
ORDER BY bp.created_at DESC;

-- name: ListChargeBatchesOnlyReportedByCloudServiceProvider :many
-- Charge batches the cloud service provider reported directly that the billing provider has
-- not (yet) reported for this billing account.
SELECT csp.* FROM cloud_service_provider_charge_batch csp
WHERE csp.billing_account_id = ?
  AND NOT EXISTS (
    SELECT 1 FROM billing_account_charge_batch bp WHERE bp.id = csp.id
  )
ORDER BY csp.created_at DESC;
