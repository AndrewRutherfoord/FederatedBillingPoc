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
