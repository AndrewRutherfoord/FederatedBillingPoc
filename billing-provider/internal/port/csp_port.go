package port

import (
	"context"
	"fmt"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

// CSPHandler handles incoming messages from a CSP.
// Implemented by CSPPortImpl and called by the inbound adapter.
type CSPHandler interface {
	OnAggregatedChargeRecord(ctx context.Context, record sharedmodels.AggregatedChargeRecord) error
}

// CSPPortImpl contains the application logic for handling inbound CSP messages.
type CSPPortImpl struct {
	repositories *repository.Repos
}

func NewCSPPort(repositories *repository.Repos) *CSPPortImpl {
	return &CSPPortImpl{repositories: repositories}
}

// OnAggregatedChargeRecord handles a cost batch pushed by a CSP.
func (p *CSPPortImpl) OnAggregatedChargeRecord(ctx context.Context, record sharedmodels.AggregatedChargeRecord) error {
	if record.BatchID == "" {
		return fmt.Errorf("batch_id is required")
	}
	if record.BillingAccountID == "" {
		return fmt.Errorf("billing_account_id is required")
	}
	if record.ResourceProviderID == "" {
		return fmt.Errorf("resource_provider_id is required")
	}

	billingPeriodID := generateBillingPeriodID(record.CreatedAt)

	costBatch := &db.CostBatch{
		ID:                     record.BatchID,
		BillingAccountID:       record.BillingAccountID,
		BillingPeriodID:        billingPeriodID,
		CloudServiceProviderID: record.ResourceProviderID,
		TotalItems:             record.LineItemCount,
		TotalCost:              record.TotalBilledCost,
		MerkelRoot:             record.BatchHash,
		Signature:              record.BatchSignature,
		CreatedAt:              record.CreatedAt,
	}

	return p.repositories.CostBatch.Create(ctx, costBatch)
}

func generateBillingPeriodID(t time.Time) string {
	return t.Format("2006-01")
}
