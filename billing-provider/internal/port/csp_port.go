package port

import (
	"context"

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
	return nil
}
