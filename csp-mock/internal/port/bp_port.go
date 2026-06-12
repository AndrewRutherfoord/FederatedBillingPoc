package port

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/adapters"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

type BPPort interface {
	// Outbound - CSP sends to BP
	SendAggregatedChargeRecord(ctx context.Context, record sharedmodels.AggregatedChargeRecord) error

	// Inbound - BP sends to CSP, port receives and surfaces to application
	OnCreditUpdate(ctx context.Context, update sharedmodels.CreditUpdate) error
}

type BPPortImpl struct {
	repositories *repository.Repos
	adapter      adapters.BillingProviderAdapter
}

func NewBPPort(
	repositories *repository.Repos,
	adapter adapters.BillingProviderAdapter,
) *BPPortImpl {
	return &BPPortImpl{
		repositories: repositories,
		adapter:      adapter,
	}
}

func (p *BPPortImpl) SendAggregatedChargeRecord(ctx context.Context, record sharedmodels.AggregatedChargeRecord) error {
	return nil
}

func (p *BPPortImpl) OnCreditUpdate(ctx context.Context, update sharedmodels.CreditUpdate) error {
	return nil
}
