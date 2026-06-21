-- name: CreateCloudServiceProviderFocusRecord :exec
INSERT INTO cloud_service_provider_focus_record (
    id, charge_batch_id,
    billed_cost, billing_account_id, billing_account_type, billing_currency,
    billing_period_end, billing_period_start, charge_category, charge_frequency,
    charge_period_end, charge_period_start, contracted_cost, effective_cost,
    host_provider_name, invoice_issuer_name, list_cost, service_category,
    service_name, service_provider_name, service_subcategory,
    allocated_method_details, allocated_method_id, allocated_resource_id, allocated_resource_name, allocated_tags,
    availability_zone, region_id, region_name,
    billing_account_name,
    capacity_reservation_id, capacity_reservation_status,
    charge_class, charge_description,
    commitment_discount_category, commitment_discount_id, commitment_discount_name, commitment_discount_quantity, commitment_discount_status, commitment_discount_type, commitment_discount_unit,
    consumed_quantity, consumed_unit,
    contract_applied, contracted_unit_price,
    invoice_id,
    list_unit_price,
    pricing_category, pricing_currency, pricing_currency_contracted_unit_price, pricing_currency_effective_cost, pricing_currency_list_unit_price, pricing_quantity, pricing_unit,
    resource_id, resource_name, resource_type,
    sku_id, sku_meter, sku_price_details, sku_price_id,
    sub_account_id, sub_account_name, sub_account_type,
    tags
) VALUES (
    ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?, ?, ?,
    ?, ?, ?,
    ?,
    ?, ?,
    ?, ?,
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?,
    ?, ?,
    ?,
    ?,
    ?, ?, ?, ?, ?, ?, ?,
    ?, ?, ?,
    ?, ?, ?, ?,
    ?, ?, ?,
    ?
);

-- name: ListCloudServiceProviderFocusRecordsByChargeBatch :many
SELECT * FROM cloud_service_provider_focus_record
WHERE charge_batch_id = ?;

-- name: ListAggregatedResourceChargesByBillingAccount :many
-- Aggregated billed cost per resource, sourced entirely from the FOCUS records fetched
-- directly from the CSP (not from anything reported by the billing provider).
SELECT
    fr.resource_id AS resource_id,
    fr.resource_name AS resource_name,
    fr.resource_type AS resource_type,
    fr.service_name AS service_name,
    fr.service_category AS service_category,
    fr.billing_currency AS billing_currency,
    cspb.cloud_service_provider_id AS cloud_service_provider_id,
    SUM(fr.billed_cost) AS total_billed_cost,
    COUNT(*) AS line_item_count
FROM cloud_service_provider_focus_record fr
JOIN cloud_service_provider_charge_batch cspb ON cspb.id = fr.charge_batch_id
WHERE fr.billing_account_id = ?
GROUP BY fr.resource_id, fr.resource_name, fr.resource_type, fr.service_name, fr.service_category, fr.billing_currency, cspb.cloud_service_provider_id
ORDER BY total_billed_cost DESC;
