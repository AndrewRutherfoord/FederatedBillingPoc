// Package main is the entry point for the CSP Mock Adapter server.
// The adapter serves as the API endpoint for billing providers to interact with the CSP.
package main

import (
	"flag"
	"log"

	apibpadapter "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/adapters/api"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/app"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/port"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	sharedscheduler "github.com/andrewrutherfoord/fed-bill-poc/shared/scheduler"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	apiPort := flag.String("port", ":8443", "port to listen on")
	flag.Parse()

	appInstance := app.NewApp(app.Config{ConfigPath: *configPath})

	handler := port.NewBPPort(appInstance.Repos)
	adapter := apibpadapter.NewApiBillingProviderAdapter(handler, appInstance.Config)
	defer adapter.Close()

	costBatchJob := scheduler.NewRecordMeteringAndCostJob("cost-batch", appInstance.Repos, appInstance.Config, adapter)
	costBatchSched, err := sharedscheduler.NewCronSchedule("0 0 * * * *")
	if err != nil {
		log.Fatalf("failed to create cost batch schedule: %v", err)
	}
	if err := appInstance.Sched.Register(costBatchJob, costBatchSched); err != nil {
		log.Fatalf("failed to register cost batch job: %v", err)
	}

	if err := adapter.Start(*apiPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
