// Package main is the entry point for the Billing Provider customer-facing API server.
// It exposes billing accounts, invoices, and payment endpoints to authenticated customers.
//
//	@title			Billing Provider – Customer API
//	@version		1.0
//	@description	Customer-facing REST API for the Billing Provider
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8081
//	@BasePath	/
package main

import (
	"flag"
	"log"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/app"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/scheduler"
	sharedscheduler "github.com/andrewrutherfoord/fed-bill-poc/shared/scheduler"
)

// Periods close on different days, so one IssueInvoicesJob instance is registered per cron schedule below.
var invoiceJobSchedules = []struct {
	id       string
	cronExpr string
	cycles   []db.BillingCycle
}{
	{"issue-invoices-weekly", "0 0 0 * * 1", []db.BillingCycle{db.BillingCycleWeekly}},                             // every Monday
	{"issue-invoices-monthly", "0 0 0 2 * *", []db.BillingCycle{db.BillingCycleMonthly, db.BillingCycleQuarterly}}, // 2nd of every month
	{"issue-invoices-annual", "0 0 0 2 1 *", []db.BillingCycle{db.BillingCycleAnnual}},                             // 2nd of January
}

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	apiPort := flag.String("port", ":8081", "port to listen on")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		log.Fatalf("failed to initialise app: %v", err)
	}

	for _, js := range invoiceJobSchedules {
		job := scheduler.NewIssueInvoicesJob(js.id, a.Repos, a.Config, js.cycles...)
		sched, err := sharedscheduler.NewCronSchedule(js.cronExpr)
		if err != nil {
			log.Fatalf("failed to create schedule for %q: %v", js.id, err)
		}
		if err := a.Sched.Register(job, sched); err != nil {
			log.Fatalf("failed to register job %q: %v", js.id, err)
		}
	}

	server := handlers.NewServer(a.Config, a.Repos)
	if err := server.Start(*apiPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
