-- +goose Up
-- Replace the narrow line-item table (a few extracted columns + a JSON payload) with a
-- full FOCUS 1.3 record, one column per FocusLineItem field. This is the raw data fetched
-- directly from the CSP, so resource-level views can query/aggregate on any FOCUS
-- dimension (resource, service, region, SKU, ...) without deserializing JSON.
DROP TABLE cloud_service_provider_charge_line_item;

CREATE TABLE cloud_service_provider_focus_record (
    id TEXT PRIMARY KEY,
    charge_batch_id TEXT NOT NULL,

    -- FOCUS required fields
    billed_cost            REAL NOT NULL,
    billing_account_id     TEXT NOT NULL,
    billing_account_type   TEXT NOT NULL,
    billing_currency       TEXT NOT NULL,
    billing_period_end     DATETIME NOT NULL,
    billing_period_start   DATETIME NOT NULL,
    charge_category        TEXT NOT NULL,
    charge_frequency       TEXT NOT NULL,
    charge_period_end      DATETIME NOT NULL,
    charge_period_start    DATETIME NOT NULL,
    contracted_cost        REAL NOT NULL,
    effective_cost         REAL NOT NULL,
    host_provider_name     TEXT NOT NULL,
    invoice_issuer_name    TEXT NOT NULL,
    list_cost              REAL NOT NULL,
    service_category       TEXT NOT NULL,
    service_name           TEXT NOT NULL,
    service_provider_name  TEXT NOT NULL,
    service_subcategory    TEXT NOT NULL,

    -- FOCUS optional fields
    allocated_method_details TEXT,
    allocated_method_id      TEXT,
    allocated_resource_id    TEXT,
    allocated_resource_name  TEXT,
    allocated_tags            TEXT,
    availability_zone         TEXT,
    region_id                  TEXT,
    region_name                TEXT,
    billing_account_name       TEXT,
    capacity_reservation_id     TEXT,
    capacity_reservation_status TEXT,
    charge_class                TEXT,
    charge_description           TEXT,
    commitment_discount_category TEXT,
    commitment_discount_id        TEXT,
    commitment_discount_name      TEXT,
    commitment_discount_quantity  REAL,
    commitment_discount_status    TEXT,
    commitment_discount_type      TEXT,
    commitment_discount_unit      TEXT,
    consumed_quantity              REAL,
    consumed_unit                  TEXT,
    contract_applied               TEXT,
    contracted_unit_price          REAL,
    invoice_id                     TEXT,
    list_unit_price                REAL,
    pricing_category               TEXT,
    pricing_currency                TEXT,
    pricing_currency_contracted_unit_price REAL,
    pricing_currency_effective_cost        REAL,
    pricing_currency_list_unit_price       REAL,
    pricing_quantity                       REAL,
    pricing_unit                           TEXT,
    resource_id                            TEXT,
    resource_name                          TEXT,
    resource_type                          TEXT,
    sku_id                                 TEXT,
    sku_meter                              TEXT,
    sku_price_details                      TEXT,
    sku_price_id                           TEXT,
    sub_account_id                         TEXT,
    sub_account_name                       TEXT,
    sub_account_type                       TEXT,
    tags                                   TEXT,

    FOREIGN KEY (charge_batch_id) REFERENCES cloud_service_provider_charge_batch(id)
);
CREATE INDEX idx_csp_focus_record_charge_batch ON cloud_service_provider_focus_record(charge_batch_id);
CREATE INDEX idx_csp_focus_record_billing_account ON cloud_service_provider_focus_record(billing_account_id);
CREATE INDEX idx_csp_focus_record_resource ON cloud_service_provider_focus_record(resource_id);

-- +goose Down
DROP INDEX IF EXISTS idx_csp_focus_record_resource;
DROP INDEX IF EXISTS idx_csp_focus_record_billing_account;
DROP INDEX IF EXISTS idx_csp_focus_record_charge_batch;
DROP TABLE cloud_service_provider_focus_record;

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
