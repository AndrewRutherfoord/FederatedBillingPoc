package handlers

import (
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/middleware"
	bpmodels "github.com/andrewrutherfoord/fed-bill-poc/shared/models/billing_provider"
	"github.com/gin-gonic/gin"
)

func (s *Server) GetCreditBalance(c *gin.Context) {
	account := middleware.BillingAccountFromContext(c)
	billingPeriod, err := s.repos.BillingPeriod.GetBillingAccountCurrentPeriod(c.Request.Context(), account.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get billing period"})
		return
	}

	var periodStart, periodEnd string
	if billingPeriod != nil {
		periodStart = billingPeriod.Start.Format(time.RFC3339)
		periodEnd = billingPeriod.End.Format(time.RFC3339)
	}

	c.JSON(200,
		bpmodels.CreditBalanceResponse{
			BillingAccountID:   account.ID,
			CreditAvailable:    100, // TODO: Placeholder
			CreditUsed:         0,
			PaymentModel:       s.config.PaymentModel,
			CreditCurrency:     s.config.BillingCurrency,
			BillingPeriodStart: periodStart,
			BillingPeriodEnd:   periodEnd,
		},
	)

}
