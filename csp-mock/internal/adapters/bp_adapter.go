package adapters

import (
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models"
)

type BillingProviderAdapter interface {
	SendCostBatchRecord(record sharedmodels.AggregatedChargeRecord) error
	Close() error
}
