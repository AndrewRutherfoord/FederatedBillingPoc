package port

import (
	"context"
	"fmt"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

// CSPHandler handles incoming messages from a CSP.
// Implemented by CSPPortImpl and called by the inbound adapter.
type CSPHandler interface {
	OnChargeBatch(ctx context.Context, batch sharedmodels.ChargeBatch) error
}

// CSPPortImpl contains the application logic for handling inbound CSP messages.
type CSPPortImpl struct {
	repositories *repository.Repos
	clock        shared.Clock
}

func NewCSPPort(repositories *repository.Repos, clock shared.Clock) *CSPPortImpl {
	return &CSPPortImpl{repositories: repositories, clock: clock}
}

// OnChargeBatch handles a charge batch pushed by a CSP.
func (p *CSPPortImpl) OnChargeBatch(ctx context.Context, batch sharedmodels.ChargeBatch) error {
	if batch.BatchID == "" {
		return fmt.Errorf("batch_id is required")
	}
	if batch.BillingAccountID == "" {
		return fmt.Errorf("billing_account_id is required")
	}
	if batch.CloudServiceProviderID == "" {
		return fmt.Errorf("cloud_service_provider_id is required")
	}

	chargeBatch := &db.ChargeBatch{
		ID:                     batch.BatchID,
		BillingAccountID:       batch.BillingAccountID,
		CloudServiceProviderID: batch.CloudServiceProviderID,
		TotalItems:             batch.LineItemCount,
		TotalCost:              batch.TotalBilledCost,
		MerkleRoot:             batch.MerkleRoot,
		BatchSignature:         batch.BatchSignature,
		CreatedAt:              batch.CreatedAt,
		ReceivedAt:             p.clock.Now(),
	}

	return p.repositories.ChargeBatch.Create(ctx, chargeBatch)
}
