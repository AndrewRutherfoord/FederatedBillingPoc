package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
)

type OnboardingSessionRepository interface {
	Create(ctx context.Context, billingProviderID, billingAccountID, returnURL string) (*db.OnboardingSession, error)
	GetByID(ctx context.Context, id string) (*db.OnboardingSession, error)
}

type onboardingSessionRepo struct {
	db *db.DB
}

func newOnboardingSessionRepo(database *db.DB) OnboardingSessionRepository {
	return &onboardingSessionRepo{db: database}
}

func (r *onboardingSessionRepo) Create(ctx context.Context, billingProviderID, billingAccountID, returnURL string) (*db.OnboardingSession, error) {
	session := db.OnboardingSession{
		ID:                uuid.NewString(),
		BillingProviderID: billingProviderID,
		BillingAccountID:  billingAccountID,
		ReturnURL:         returnURL,
	}
	if err := r.db.WithContext(ctx).Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *onboardingSessionRepo) GetByID(ctx context.Context, id string) (*db.OnboardingSession, error) {
	var session db.OnboardingSession
	if err := r.db.WithContext(ctx).First(&session, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}
