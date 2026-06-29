-- name: GetInvoice :one
SELECT * FROM invoice WHERE id = ?;

-- name: CreateInvoice :one
INSERT INTO invoice (id, billing_account_id, billing_period_id, amount, currency, status, issued_at, due_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: ListInvoicesByBillingAccount :many
SELECT * FROM invoice WHERE billing_account_id = ? ORDER BY issued_at DESC;

-- name: CreateInvoiceProviderLineItem :one
INSERT INTO invoice_provider_line_item (id, invoice_id, cloud_service_provider_id, amount, merkle_root, batch_count, merkle_valid)
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: ListInvoiceProviderLineItemsByInvoiceIDs :many
SELECT * FROM invoice_provider_line_item WHERE invoice_id IN (sqlc.slice('invoice_ids'));
