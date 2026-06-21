-- +goose Up

-- Rename the messy/inconsistent names left over from earlier iterations:
--   * "billing_provider_supported_cloud_provider" is really just the CSP catalog.
--   * "cloud_provider_id" vs "cloud_service_provider_id" was used interchangeably
--     for the same foreign key across tables.
--   * "cost_batch"/"cost"-prefixed naming clashed with the FOCUS spec's own use of
--     "cost" for monetary amounts on a charge; these tables hold charge batches.
ALTER TABLE billing_provider_supported_cloud_provider RENAME TO cloud_service_providers;

ALTER TABLE cloud_service_provider_accounts RENAME COLUMN cloud_provider_id TO cloud_service_provider_id;

ALTER TABLE billing_account_cost_batch RENAME TO billing_account_charge_batch;
ALTER TABLE billing_account_charge_batch RENAME COLUMN merkel_root TO merkle_root;

-- billing_account_charge_batch already holds what the billing provider reported for a
-- charge batch. This table holds what the cloud service provider reported directly for
-- the same batch ID, fetched independently of the billing provider, so the two reports
-- can be compared to spot discrepancies.
CREATE TABLE cloud_service_provider_charge_batch (
    id TEXT PRIMARY KEY,
    billing_account_id TEXT NOT NULL,
    cloud_service_provider_id TEXT NOT NULL,
    total_items INTEGER NOT NULL,
    total_cost REAL NOT NULL,
    billed_currency TEXT NOT NULL DEFAULT '',
    merkle_root TEXT NOT NULL,
    batch_signature TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    received_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (billing_account_id) REFERENCES billing_accounts(id),
    FOREIGN KEY (cloud_service_provider_id) REFERENCES cloud_service_providers(id)
);

-- Raw FOCUS line items underlying a cloud_service_provider_charge_batch, fetched directly
-- from the CSP so the customer-facing side has access to the underlying detail, not just
-- the aggregated batch totals. Only the fields useful for lookups/auditing are broken out
-- as columns; the full FOCUS line item is kept verbatim in `payload`.
CREATE TABLE cloud_service_provider_charge_line_item (
    id TEXT PRIMARY KEY,
    charge_batch_id TEXT NOT NULL,
    billing_account_id TEXT NOT NULL,
    billed_cost REAL NOT NULL,
    service_name TEXT NOT NULL,
    charge_period_start DATETIME NOT NULL,
    charge_period_end DATETIME NOT NULL,
    payload TEXT NOT NULL
);
CREATE INDEX idx_csp_charge_line_item_batch ON cloud_service_provider_charge_line_item(charge_batch_id);

-- +goose Down
DROP INDEX IF EXISTS idx_csp_charge_line_item_batch;
DROP TABLE cloud_service_provider_charge_line_item;
DROP TABLE cloud_service_provider_charge_batch;

ALTER TABLE billing_account_charge_batch RENAME COLUMN merkle_root TO merkel_root;
ALTER TABLE billing_account_charge_batch RENAME TO billing_account_cost_batch;

ALTER TABLE cloud_service_provider_accounts RENAME COLUMN cloud_service_provider_id TO cloud_provider_id;

ALTER TABLE cloud_service_providers RENAME TO billing_provider_supported_cloud_provider;
