package port

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

// BPSender is what the application calls to send things to the BP.
// Implemented by the outbound adapter.
type BPSender interface {
	SendChargeBatch(ctx context.Context, batch sharedmodels.ChargeBatch) error
}

// BPHandler handles incoming messages from the BP.
// Implemented by BPPortImpl and called by the inbound adapter.
type BPHandler interface {
	OnCreditUpdate(ctx context.Context, update sharedmodels.CreditUpdate) error
}

// BPPortImpl contains the application logic for handling inbound BP messages.
type BPPortImpl struct {
	repositories *repository.Repos
}

func NewBPPort(repositories *repository.Repos) *BPPortImpl {
	return &BPPortImpl{repositories: repositories}
}

// OnCreditUpdate handles an incoming credit update from the BP.
func (p *BPPortImpl) OnCreditUpdate(ctx context.Context, update sharedmodels.CreditUpdate) error {
	return nil
}
