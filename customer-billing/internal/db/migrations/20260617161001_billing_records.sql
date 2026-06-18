-- +goose Up
CREATE TABLE billing_account_cost_batch (
    id TEXT PRIMARY KEY, 
    billing_account_id TEXT NOT NULL,
    billing_period_id TEXT NOT NULL,
    cloud_service_provider_id TEXT NOT NULL,
    total_items INTEGER NOT NULL,
    total_cost REAL NOT NULL,
    billed_currency TEXT NOT NULL DEFAULT '',
    merkel_root TEXT NOT NULL,
    batch_signature TEXT NOT NULL DEFAULT '',
    created_at DATETIME NOT NULL,
    FOREIGN KEY (billing_account_id) REFERENCES billing_accounts(id),
    FOREIGN KEY (cloud_service_provider_id) REFERENCES billing_provider_supported_cloud_provider(id)
);

-- +goose Down
DROP TABLE billing_account_cost_batch;