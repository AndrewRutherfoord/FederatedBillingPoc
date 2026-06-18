-- +goose Up
ALTER TABLE billing_account_cost_batch ADD COLUMN received_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- +goose Down
ALTER TABLE billing_account_cost_batch DROP COLUMN received_at;
