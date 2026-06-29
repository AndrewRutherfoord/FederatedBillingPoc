-- +goose Up
-- Remove duplicate links created by repeated calls to the onboarding-complete
-- endpoint (e.g. the customer refreshing the redirect-back page), keeping the
-- earliest row for each billing account / cloud provider pair.
DELETE FROM cloud_service_provider_accounts
WHERE id NOT IN (
    SELECT MIN(id) FROM cloud_service_provider_accounts
    GROUP BY billing_account_id, cloud_service_provider_id
);

CREATE UNIQUE INDEX idx_cloud_service_provider_accounts_billing_account_csp
    ON cloud_service_provider_accounts(billing_account_id, cloud_service_provider_id);

-- +goose Down
DROP INDEX idx_cloud_service_provider_accounts_billing_account_csp;
