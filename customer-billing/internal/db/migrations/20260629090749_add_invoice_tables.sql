-- +goose Up
CREATE TABLE invoice (
    id TEXT PRIMARY KEY,
    billing_account_id TEXT NOT NULL,
    billing_period_id TEXT NOT NULL,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL,
    issued_at DATETIME NOT NULL,
    due_at DATETIME NOT NULL,
    FOREIGN KEY (billing_account_id) REFERENCES billing_accounts(id)
);

CREATE TABLE invoice_provider_line_item (
    id TEXT PRIMARY KEY,
    invoice_id TEXT NOT NULL,
    cloud_service_provider_id TEXT NOT NULL,
    amount REAL NOT NULL,
    merkle_root TEXT NOT NULL,
    batch_count INTEGER NOT NULL,
    FOREIGN KEY (invoice_id) REFERENCES invoice(id),
    FOREIGN KEY (cloud_service_provider_id) REFERENCES cloud_service_providers(id)
);
CREATE INDEX idx_invoice_provider_line_item_invoice ON invoice_provider_line_item(invoice_id);

-- +goose Down
DROP INDEX IF EXISTS idx_invoice_provider_line_item_invoice;
DROP TABLE invoice_provider_line_item;
DROP TABLE invoice;
