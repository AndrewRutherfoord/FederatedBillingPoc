package repository

import (
	"context"
	"database/sql"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	sqlcdb "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ResourceCharge is the total billed cost for a single resource (or resource-less charge
// group), aggregated entirely from the FOCUS records fetched directly from the CSP.
type ResourceCharge struct {
	ResourceID             *string
	ResourceName           *string
	ResourceType           *string
	ServiceName            string
	ServiceCategory        string
	BillingCurrency        string
	CloudServiceProviderID string
	TotalBilledCost        float64
	LineItemCount          int32
}

// CloudServiceProviderFocusRecordRepository stores the raw FOCUS line items fetched
// directly from a CSP for a charge batch, giving the customer-facing side access to the
// underlying detail rather than just the batch's aggregated totals.
type CloudServiceProviderFocusRecordRepository interface {
	CreateBatch(ctx context.Context, chargeBatchID string, billingAccountID string, items []sharedmodels.FocusLineItem) error
	ListByChargeBatch(ctx context.Context, chargeBatchID string) ([]sharedmodels.FocusLineItem, error)
	ListAggregatedResourceCharges(ctx context.Context, billingAccountID string) ([]ResourceCharge, error)
}

type cloudServiceProviderFocusRecordRepo struct {
	db *db.DB
}

func newCloudServiceProviderFocusRecordRepo(database *db.DB) CloudServiceProviderFocusRecordRepository {
	return &cloudServiceProviderFocusRecordRepo{db: database}
}

func (r *cloudServiceProviderFocusRecordRepo) CreateBatch(ctx context.Context, chargeBatchID string, billingAccountID string, items []sharedmodels.FocusLineItem) error {
	return r.db.WithTx(ctx, func(q *sqlcdb.Queries) error {
		for _, item := range items {
			if err := q.CreateCloudServiceProviderFocusRecord(ctx, toCreateFocusRecordParams(uuid.NewString(), chargeBatchID, item)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *cloudServiceProviderFocusRecordRepo) ListByChargeBatch(ctx context.Context, chargeBatchID string) ([]sharedmodels.FocusLineItem, error) {
	rows, err := r.db.Queries.ListCloudServiceProviderFocusRecordsByChargeBatch(ctx, chargeBatchID)
	if err != nil {
		return nil, err
	}
	items := make([]sharedmodels.FocusLineItem, len(rows))
	for i, row := range rows {
		items[i] = fromFocusRecordRow(row)
	}
	return items, nil
}

func (r *cloudServiceProviderFocusRecordRepo) ListAggregatedResourceCharges(ctx context.Context, billingAccountID string) ([]ResourceCharge, error) {
	rows, err := r.db.Queries.ListAggregatedResourceChargesByBillingAccount(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}
	result := make([]ResourceCharge, len(rows))
	for i, row := range rows {
		result[i] = ResourceCharge{
			ResourceID:             stringPtr(row.ResourceID),
			ResourceName:           stringPtr(row.ResourceName),
			ResourceType:           stringPtr(row.ResourceType),
			ServiceName:            row.ServiceName,
			ServiceCategory:        row.ServiceCategory,
			BillingCurrency:        row.BillingCurrency,
			CloudServiceProviderID: row.CloudServiceProviderID,
			TotalBilledCost:        row.TotalBilledCost.Float64,
			LineItemCount:          int32(row.LineItemCount),
		}
	}
	return result, nil
}

func nullableString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

func stringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	v := ns.String
	return &v
}

func nullableFloat(d *decimal.Decimal) sql.NullFloat64 {
	if d == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: d.InexactFloat64(), Valid: true}
}

func decimalPtr(nf sql.NullFloat64) *decimal.Decimal {
	if !nf.Valid {
		return nil
	}
	v := decimal.NewFromFloat(nf.Float64)
	return &v
}

func toCreateFocusRecordParams(id string, chargeBatchID string, item sharedmodels.FocusLineItem) sqlcdb.CreateCloudServiceProviderFocusRecordParams {
	return sqlcdb.CreateCloudServiceProviderFocusRecordParams{
		ID:                                 id,
		ChargeBatchID:                      chargeBatchID,
		BilledCost:                         item.BilledCost.InexactFloat64(),
		BillingAccountID:                   item.BillingAccountID,
		BillingAccountType:                 item.BillingAccountType,
		BillingCurrency:                    item.BillingCurrency,
		BillingPeriodEnd:                   item.BillingPeriodEnd,
		BillingPeriodStart:                 item.BillingPeriodStart,
		ChargeCategory:                     item.ChargeCategory,
		ChargeFrequency:                    item.ChargeFrequency,
		ChargePeriodEnd:                    item.ChargePeriodEnd,
		ChargePeriodStart:                  item.ChargePeriodStart,
		ContractedCost:                     item.ContractedCost.InexactFloat64(),
		EffectiveCost:                      item.EffectiveCost.InexactFloat64(),
		HostProviderName:                   item.HostProviderName,
		InvoiceIssuerName:                  item.InvoiceIssuerName,
		ListCost:                           item.ListCost.InexactFloat64(),
		ServiceCategory:                    string(item.ServiceCategory),
		ServiceName:                        item.ServiceName,
		ServiceProviderName:                item.ServiceProviderName,
		ServiceSubcategory:                 string(item.ServiceSubcategory),
		AllocatedMethodDetails:             nullableString(item.AllocatedMethodDetails),
		AllocatedMethodID:                  nullableString(item.AllocatedMethodID),
		AllocatedResourceID:                nullableString(item.AllocatedResourceID),
		AllocatedResourceName:              nullableString(item.AllocatedResourceName),
		AllocatedTags:                      nullableString(item.AllocatedTags),
		AvailabilityZone:                   nullableString(item.AvailabilityZone),
		RegionID:                           nullableString(item.RegionID),
		RegionName:                         nullableString(item.RegionName),
		BillingAccountName:                 nullableString(item.BillingAccountName),
		CapacityReservationID:              nullableString(item.CapacityReservationID),
		CapacityReservationStatus:          nullableCapacityReservationStatus(item.CapacityReservationStatus),
		ChargeClass:                        nullableString(item.ChargeClass),
		ChargeDescription:                  nullableString(item.ChargeDescription),
		CommitmentDiscountCategory:         nullableCommitmentDiscountCategory(item.CommitmentDiscountCategory),
		CommitmentDiscountID:               nullableString(item.CommitmentDiscountID),
		CommitmentDiscountName:             nullableString(item.CommitmentDiscountName),
		CommitmentDiscountQuantity:         nullableFloat(item.CommitmentDiscountQuantity),
		CommitmentDiscountStatus:           nullableCommitmentDiscountStatus(item.CommitmentDiscountStatus),
		CommitmentDiscountType:             nullableString(item.CommitmentDiscountType),
		CommitmentDiscountUnit:             nullableString(item.CommitmentDiscountUnit),
		ConsumedQuantity:                   nullableFloat(item.ConsumedQuantity),
		ConsumedUnit:                       nullableString(item.ConsumedUnit),
		ContractApplied:                    nullableString(item.ContractApplied),
		ContractedUnitPrice:                nullableFloat(item.ContractedUnitPrice),
		InvoiceID:                          nullableString(item.InvoiceID),
		ListUnitPrice:                      nullableFloat(item.ListUnitPrice),
		PricingCategory:                    nullableString(item.PricingCategory),
		PricingCurrency:                    nullableString(item.PricingCurrency),
		PricingCurrencyContractedUnitPrice: nullableFloat(item.PricingCurrencyContractedUnitPrice),
		PricingCurrencyEffectiveCost:       nullableFloat(item.PricingCurrencyEffectiveCost),
		PricingCurrencyListUnitPrice:       nullableFloat(item.PricingCurrencyListUnitPrice),
		PricingQuantity:                    nullableFloat(item.PricingQuantity),
		PricingUnit:                        nullableString(item.PricingUnit),
		ResourceID:                         nullableString(item.ResourceID),
		ResourceName:                       nullableString(item.ResourceName),
		ResourceType:                       nullableString(item.ResourceType),
		SkuID:                              nullableString(item.SkuID),
		SkuMeter:                           nullableString(item.SkuMeter),
		SkuPriceDetails:                    nullableString(item.SkuPriceDetails),
		SkuPriceID:                         nullableString(item.SkuPriceID),
		SubAccountID:                       nullableString(item.SubAccountID),
		SubAccountName:                     nullableString(item.SubAccountName),
		SubAccountType:                     nullableString(item.SubAccountType),
		Tags:                               nullableString(item.Tags),
	}
}

func fromFocusRecordRow(row sqlcdb.CloudServiceProviderFocusRecord) sharedmodels.FocusLineItem {
	return sharedmodels.FocusLineItem{
		BilledCost:                         decimal.NewFromFloat(row.BilledCost),
		BillingAccountID:                   row.BillingAccountID,
		BillingAccountType:                 row.BillingAccountType,
		BillingCurrency:                    row.BillingCurrency,
		BillingPeriodEnd:                   row.BillingPeriodEnd,
		BillingPeriodStart:                 row.BillingPeriodStart,
		ChargeCategory:                     row.ChargeCategory,
		ChargeFrequency:                    row.ChargeFrequency,
		ChargePeriodEnd:                    row.ChargePeriodEnd,
		ChargePeriodStart:                  row.ChargePeriodStart,
		ContractedCost:                     decimal.NewFromFloat(row.ContractedCost),
		EffectiveCost:                      decimal.NewFromFloat(row.EffectiveCost),
		HostProviderName:                   row.HostProviderName,
		InvoiceIssuerName:                  row.InvoiceIssuerName,
		ListCost:                           decimal.NewFromFloat(row.ListCost),
		ServiceCategory:                    sharedmodels.ServiceCategory(row.ServiceCategory),
		ServiceName:                        row.ServiceName,
		ServiceProviderName:                row.ServiceProviderName,
		ServiceSubcategory:                 sharedmodels.ServiceSubcategory(row.ServiceSubcategory),
		AllocatedMethodDetails:             stringPtr(row.AllocatedMethodDetails),
		AllocatedMethodID:                  stringPtr(row.AllocatedMethodID),
		AllocatedResourceID:                stringPtr(row.AllocatedResourceID),
		AllocatedResourceName:              stringPtr(row.AllocatedResourceName),
		AllocatedTags:                      stringPtr(row.AllocatedTags),
		AvailabilityZone:                   stringPtr(row.AvailabilityZone),
		RegionID:                           stringPtr(row.RegionID),
		RegionName:                         stringPtr(row.RegionName),
		BillingAccountName:                 stringPtr(row.BillingAccountName),
		CapacityReservationID:              stringPtr(row.CapacityReservationID),
		CapacityReservationStatus:          capacityReservationStatusPtr(row.CapacityReservationStatus),
		ChargeClass:                        stringPtr(row.ChargeClass),
		ChargeDescription:                  stringPtr(row.ChargeDescription),
		CommitmentDiscountCategory:         commitmentDiscountCategoryPtr(row.CommitmentDiscountCategory),
		CommitmentDiscountID:               stringPtr(row.CommitmentDiscountID),
		CommitmentDiscountName:             stringPtr(row.CommitmentDiscountName),
		CommitmentDiscountQuantity:         decimalPtr(row.CommitmentDiscountQuantity),
		CommitmentDiscountStatus:           commitmentDiscountStatusPtr(row.CommitmentDiscountStatus),
		CommitmentDiscountType:             stringPtr(row.CommitmentDiscountType),
		CommitmentDiscountUnit:             stringPtr(row.CommitmentDiscountUnit),
		ConsumedQuantity:                   decimalPtr(row.ConsumedQuantity),
		ConsumedUnit:                       stringPtr(row.ConsumedUnit),
		ContractApplied:                    stringPtr(row.ContractApplied),
		ContractedUnitPrice:                decimalPtr(row.ContractedUnitPrice),
		InvoiceID:                          stringPtr(row.InvoiceID),
		ListUnitPrice:                      decimalPtr(row.ListUnitPrice),
		PricingCategory:                    stringPtr(row.PricingCategory),
		PricingCurrency:                    stringPtr(row.PricingCurrency),
		PricingCurrencyContractedUnitPrice: decimalPtr(row.PricingCurrencyContractedUnitPrice),
		PricingCurrencyEffectiveCost:       decimalPtr(row.PricingCurrencyEffectiveCost),
		PricingCurrencyListUnitPrice:       decimalPtr(row.PricingCurrencyListUnitPrice),
		PricingQuantity:                    decimalPtr(row.PricingQuantity),
		PricingUnit:                        stringPtr(row.PricingUnit),
		ResourceID:                         stringPtr(row.ResourceID),
		ResourceName:                       stringPtr(row.ResourceName),
		ResourceType:                       stringPtr(row.ResourceType),
		SkuID:                              stringPtr(row.SkuID),
		SkuMeter:                           stringPtr(row.SkuMeter),
		SkuPriceDetails:                    stringPtr(row.SkuPriceDetails),
		SkuPriceID:                         stringPtr(row.SkuPriceID),
		SubAccountID:                       stringPtr(row.SubAccountID),
		SubAccountName:                     stringPtr(row.SubAccountName),
		SubAccountType:                     stringPtr(row.SubAccountType),
		Tags:                               stringPtr(row.Tags),
	}
}

func nullableCapacityReservationStatus(s *sharedmodels.CapacityReservationStatus) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

func nullableCommitmentDiscountCategory(s *sharedmodels.CommitmentDiscountCategory) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

func nullableCommitmentDiscountStatus(s *sharedmodels.CommitmentDiscountStatus) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

func capacityReservationStatusPtr(ns sql.NullString) *sharedmodels.CapacityReservationStatus {
	if !ns.Valid {
		return nil
	}
	v := sharedmodels.CapacityReservationStatus(ns.String)
	return &v
}

func commitmentDiscountCategoryPtr(ns sql.NullString) *sharedmodels.CommitmentDiscountCategory {
	if !ns.Valid {
		return nil
	}
	v := sharedmodels.CommitmentDiscountCategory(ns.String)
	return &v
}

func commitmentDiscountStatusPtr(ns sql.NullString) *sharedmodels.CommitmentDiscountStatus {
	if !ns.Valid {
		return nil
	}
	v := sharedmodels.CommitmentDiscountStatus(ns.String)
	return &v
}
