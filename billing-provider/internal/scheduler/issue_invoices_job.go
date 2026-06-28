package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
)

// IssueInvoicesJob invoices elapsed billing periods for accounts whose cycle is in `cycles`.
// Register one instance per cron schedule (weekly/monthly/annual) with the matching cycles.
type IssueInvoicesJob struct {
	id     string
	repos  *repository.Repos
	config *config.Config
	cycles map[db.BillingCycle]bool
}

func NewIssueInvoicesJob(id string, repos *repository.Repos, cfg *config.Config, cycles ...db.BillingCycle) *IssueInvoicesJob {
	set := make(map[db.BillingCycle]bool, len(cycles))
	for _, c := range cycles {
		set[c] = true
	}
	return &IssueInvoicesJob{id: id, repos: repos, config: cfg, cycles: set}
}

func (j *IssueInvoicesJob) ID() string {
	return j.id
}

func (j *IssueInvoicesJob) Execute(ctx context.Context, startTime time.Time) error {
	accounts, err := j.repos.BillingAccounts.List(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		cycle := account.BillingCycle
		if !cycle.Valid() {
			cycle = db.BillingCycleMonthly
		}
		if !j.cycles[cycle] {
			continue
		}

		if err := j.processAccount(ctx, account, cycle, startTime); err != nil {
			log.Printf("IssueInvoicesJob: error processing billing account %s: %v", account.ID, err)
		}
	}
	return nil
}

func (j *IssueInvoicesJob) processAccount(ctx context.Context, account db.BillingAccount, cycle db.BillingCycle, now time.Time) error {
	unassociated, err := j.repos.ChargeBatch.GetUnassociatedByBillingAccount(ctx, account.ID)
	if err != nil {
		return err
	}
	if len(unassociated) == 0 {
		return nil
	}

	current, err := j.repos.BillingPeriod.GetBillingAccountCurrentPeriod(ctx, account.ID)
	if err != nil {
		return err
	}
	if current == nil {
		current, err = j.repos.BillingPeriod.GetOrCreatePeriodForTime(ctx, account.ID, cycle, unassociated[0].CreatedAt)
		if err != nil {
			return err
		}
	}

	for !current.End.After(now) {
		var batchesInPeriod []db.ChargeBatch
		var batchIDs []string
		for _, b := range unassociated {
			if !b.CreatedAt.Before(current.Start) && b.CreatedAt.Before(current.End) {
				batchesInPeriod = append(batchesInPeriod, b)
				batchIDs = append(batchIDs, b.ID)
			}
		}

		if len(batchesInPeriod) > 0 {
			if _, err := j.repos.Invoice.CreateForPeriod(ctx, account.ID, current.ID, j.config.BillingCurrency, batchesInPeriod, now); err != nil {
				return err
			}
			if err := j.repos.ChargeBatch.AssociateWithPeriod(ctx, batchIDs, current.ID); err != nil {
				return err
			}
		}

		current, err = j.repos.BillingPeriod.GetOrCreatePeriodForTime(ctx, account.ID, cycle, current.End)
		if err != nil {
			return err
		}
	}

	return nil
}
