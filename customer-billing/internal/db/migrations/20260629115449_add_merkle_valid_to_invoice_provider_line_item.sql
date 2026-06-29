-- +goose Up
ALTER TABLE invoice_provider_line_item ADD COLUMN merkle_valid BOOLEAN NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE invoice_provider_line_item DROP COLUMN merkle_valid;
