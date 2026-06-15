-- name: CreateBillingAccount :exec
INSERT INTO billing_accounts (id, billing_provider_id, alias, created_at)
VALUES (?, ?, ?, ?);

-- name: GetBillingAccount :one
SELECT * FROM billing_accounts
WHERE id = ?;

-- name: ListBillingAccountsByProvider :many
SELECT * FROM billing_accounts
WHERE billing_provider_id = ?;

-- name: UpdateBillingAccountOnboarding :exec
UPDATE billing_accounts
SET onboarding_redirect_url = ?, onboarding_complete = ?
WHERE id = ?;

-- name: ListBillingAccounts :many
SELECT ba.id, ba.billing_provider_id, ba.alias, ba.created_at, ba.onboarding_redirect_url, ba.onboarding_complete, bp.name AS billing_provider_name
FROM billing_accounts ba
LEFT JOIN billing_providers bp ON bp.id = ba.billing_provider_id;