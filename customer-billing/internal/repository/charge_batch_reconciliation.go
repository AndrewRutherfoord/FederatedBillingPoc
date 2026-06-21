package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
)

// ChargeBatchReconciliationStatus describes whether a charge batch's billing provider and
// cloud service provider reports agree.
type ChargeBatchReconciliationStatus string

const (
	ChargeBatchStatusMatched        ChargeBatchReconciliationStatus = "matched"
	ChargeBatchStatusMismatched     ChargeBatchReconciliationStatus = "mismatched"
	ChargeBatchStatusMissingFromCSP ChargeBatchReconciliationStatus = "missing_from_csp"
	ChargeBatchStatusMissingFromBP  ChargeBatchReconciliationStatus = "missing_from_bp"
)

// ChargeBatchReport is one party's (the billing provider's or the cloud service provider's)
// reported figures for a charge batch.
type ChargeBatchReport struct {
	TotalItems     int32
	TotalCost      float64
	BilledCurrency string
	MerkleRoot     string
	BatchSignature string
	CreatedAt      time.Time
	ReceivedAt     time.Time
}

// MergedChargeBatch is a single charge batch with both parties' reports side by side, so
// the customer-facing side can display them together and spot any disagreement.
type MergedChargeBatch struct {
	BatchID                string
	BillingAccountID       string
	CloudServiceProviderID string
	BillingProviderReport  *ChargeBatchReport
	CloudProviderReport    *ChargeBatchReport
	Status                 ChargeBatchReconciliationStatus
}

type ChargeBatchReconciliationRepository interface {
	ListMergedBatches(ctx context.Context, billingAccountID string) ([]MergedChargeBatch, error)
}

type chargeBatchReconciliationRepo struct {
	db *db.DB
}

func newChargeBatchReconciliationRepo(database *db.DB) ChargeBatchReconciliationRepository {
	return &chargeBatchReconciliationRepo{db: database}
}

func (r *chargeBatchReconciliationRepo) ListMergedBatches(ctx context.Context, billingAccountID string) ([]MergedChargeBatch, error) {
	bpRows, err := r.db.Queries.ListChargeBatchesReportedByBillingProvider(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}

	result := make([]MergedChargeBatch, 0, len(bpRows))
	for _, row := range bpRows {
		merged := MergedChargeBatch{
			BatchID:                row.BatchID,
			BillingAccountID:       row.BillingAccountID,
			CloudServiceProviderID: row.CloudServiceProviderID,
			BillingProviderReport: &ChargeBatchReport{
				TotalItems:     int32(row.BpTotalItems),
				TotalCost:      row.BpTotalCost,
				BilledCurrency: row.BpBilledCurrency,
				MerkleRoot:     row.BpMerkleRoot,
				BatchSignature: row.BpBatchSignature,
				CreatedAt:      row.BpCreatedAt,
				ReceivedAt:     row.BpReceivedAt,
			},
		}

		if row.CspTotalItems.Valid {
			merged.CloudProviderReport = &ChargeBatchReport{
				TotalItems:     int32(row.CspTotalItems.Int64),
				TotalCost:      row.CspTotalCost.Float64,
				BilledCurrency: row.CspBilledCurrency.String,
				MerkleRoot:     row.CspMerkleRoot.String,
				BatchSignature: row.CspBatchSignature.String,
				CreatedAt:      row.CspCreatedAt.Time,
				ReceivedAt:     row.CspReceivedAt.Time,
			}
		}

		merged.Status = reconciliationStatus(merged.BillingProviderReport, merged.CloudProviderReport)
		result = append(result, merged)
	}

	cspOnlyRows, err := r.db.Queries.ListChargeBatchesOnlyReportedByCloudServiceProvider(ctx, billingAccountID)
	if err != nil {
		return nil, err
	}
	for _, row := range cspOnlyRows {
		result = append(result, MergedChargeBatch{
			BatchID:                row.ID,
			BillingAccountID:       row.BillingAccountID,
			CloudServiceProviderID: row.CloudServiceProviderID,
			CloudProviderReport: &ChargeBatchReport{
				TotalItems:     int32(row.TotalItems),
				TotalCost:      row.TotalCost,
				BilledCurrency: row.BilledCurrency,
				MerkleRoot:     row.MerkleRoot,
				BatchSignature: row.BatchSignature,
				CreatedAt:      row.CreatedAt,
				ReceivedAt:     row.ReceivedAt,
			},
			Status: ChargeBatchStatusMissingFromBP,
		})
	}

	return result, nil
}

func reconciliationStatus(bp *ChargeBatchReport, csp *ChargeBatchReport) ChargeBatchReconciliationStatus {
	if csp == nil {
		return ChargeBatchStatusMissingFromCSP
	}
	if bp == nil {
		return ChargeBatchStatusMissingFromBP
	}
	if bp.TotalItems != csp.TotalItems || bp.TotalCost != csp.TotalCost || bp.MerkleRoot != csp.MerkleRoot {
		return ChargeBatchStatusMismatched
	}
	return ChargeBatchStatusMatched
}
