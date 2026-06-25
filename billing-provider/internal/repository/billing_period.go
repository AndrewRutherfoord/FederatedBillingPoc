package repository

import (
	"context"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
)

type BillingPeriod struct {
	ID    string
	Start *time.Time
	End   *time.Time
}

type BillingPeriodRepository interface {
	GetBillingAccountCurrentPeriod(ctx context.Context, billingAccountID string) (*BillingPeriod, error)
}

type billingPeriodRepo struct {
	db *db.DB
}

func newBillingPeriodRepo(database *db.DB) BillingPeriodRepository {
	return &billingPeriodRepo{db: database}
}

func (r *billingPeriodRepo) GetBillingAccountCurrentPeriod(ctx context.Context, billingAccountID string) (*BillingPeriod, error) {
	// Placeholder for actual implementation stored in DB
	startOfCurrentMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfCurrentMonth := startOfCurrentMonth.AddDate(0, 1, -1)

	return &BillingPeriod{
		ID:    "current-period",
		Start: &startOfCurrentMonth,
		End:   &endOfCurrentMonth,
	}, nil
}
