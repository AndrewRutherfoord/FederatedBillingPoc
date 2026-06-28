package repository

import (
	"context"
	"errors"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BillingPeriod struct {
	ID    string
	Start time.Time
	End   time.Time
}

type BillingPeriodRepository interface {
	GetBillingAccountCurrentPeriod(ctx context.Context, billingAccountID string) (*BillingPeriod, error)
	GetOrCreatePeriodForTime(ctx context.Context, billingAccountID string, cycle db.BillingCycle, at time.Time) (*BillingPeriod, error)
}

type billingPeriodRepo struct {
	db *db.DB
}

func newBillingPeriodRepo(database *db.DB) BillingPeriodRepository {
	return &billingPeriodRepo{db: database}
}

func (r *billingPeriodRepo) GetBillingAccountCurrentPeriod(ctx context.Context, billingAccountID string) (*BillingPeriod, error) {
	var period db.BillingPeriod
	err := r.db.WithContext(ctx).
		Where("billing_account_id = ?", billingAccountID).
		Order("period_start DESC").
		First(&period).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return toRepoBillingPeriod(&period), nil
}

func (r *billingPeriodRepo) GetOrCreatePeriodForTime(ctx context.Context, billingAccountID string, cycle db.BillingCycle, at time.Time) (*BillingPeriod, error) {
	var covering db.BillingPeriod
	err := r.db.WithContext(ctx).
		Where("billing_account_id = ? AND period_start <= ? AND period_end > ?", billingAccountID, at, at).
		First(&covering).Error
	if err == nil {
		return toRepoBillingPeriod(&covering), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var latest db.BillingPeriod
	err = r.db.WithContext(ctx).
		Where("billing_account_id = ?", billingAccountID).
		Order("period_start DESC").
		First(&latest).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// First charge ever for this account anchors its first period.
		first := db.BillingPeriod{
			ID:               uuid.New().String(),
			BillingAccountID: billingAccountID,
			PeriodStart:      at,
			PeriodEnd:        cycle.Advance(at),
			Status:           "open",
			OpenedAt:         at,
		}
		if err := r.db.WithContext(ctx).Create(&first).Error; err != nil {
			return nil, err
		}
		return toRepoBillingPeriod(&first), nil
	}
	if err != nil {
		return nil, err
	}

	if at.Before(latest.PeriodStart) {
		var earliest db.BillingPeriod
		if err := r.db.WithContext(ctx).
			Where("billing_account_id = ?", billingAccountID).
			Order("period_start ASC").
			First(&earliest).Error; err != nil {
			return nil, err
		}
		return toRepoBillingPeriod(&earliest), nil
	}

	current := latest
	for !current.PeriodEnd.After(at) {
		current.Status = "closed"
		current.ClosedAt = current.PeriodEnd
		if err := r.db.WithContext(ctx).Save(&current).Error; err != nil {
			return nil, err
		}

		next := db.BillingPeriod{
			ID:               uuid.New().String(),
			BillingAccountID: billingAccountID,
			PeriodStart:      current.PeriodEnd,
			PeriodEnd:        cycle.Advance(current.PeriodEnd),
			Status:           "open",
			OpenedAt:         current.PeriodEnd,
		}
		if err := r.db.WithContext(ctx).Create(&next).Error; err != nil {
			return nil, err
		}
		current = next
	}

	return toRepoBillingPeriod(&current), nil
}

func toRepoBillingPeriod(p *db.BillingPeriod) *BillingPeriod {
	return &BillingPeriod{ID: p.ID, Start: p.PeriodStart, End: p.PeriodEnd}
}
