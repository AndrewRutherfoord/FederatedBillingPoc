package main

// Package main is the entry point for the CSP Mock API server.
//
//	@title			CSP Mock API
//	@version		1.0
//	@description	Mock Cloud Service Provider billing API implementing the FOCUS spec.
//
//	@contact.name	Andrew Rutherfoord
//
//	@host		localhost:8082
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Bearer <account_id>

import (
	"flag"
	"log"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/handlers"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/scheduler"
	_ "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/docs"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedscheduler "github.com/andrewrutherfoord/fed-bill-poc/shared/scheduler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbPath := os.Getenv("CB_DB_PATH")
	if dbPath == "" {
		dbPath = "customer-billing.sqlite"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	repos := repository.New(config, database)

	// clientRegistry, err := bpclient.NewBPClientRegistry(cfg)
	// if err != nil {
	// 	log.Fatalf("failed to create billing provider client registry: %v", err)
	// }

	sched := sharedscheduler.NewWithPersistence(scheduler.NewSchedulerPersistence())
	err = sharedscheduler.RegisterJobs(sched, []sharedscheduler.JobToRegister{
		// sharedscheduler.NewJobToRegister(
		// scheduler.NewRecordMeteringAndCostJob("record-metering-and-cost", repos, cfg, clientRegistry),
		// "0 0 * * * *", // Every hour at :00 seconds
		// ),
	})
	if err != nil {
		log.Fatalf("failed to register jobs: %v", err)
	}

	clockHost := os.Getenv("MOCK_CLOCK_HOST")
	if clockHost == "" {
		clockHost = "http://localhost:9999"
	}

	// Mock clock that allows manual time advancement for testing. It uses a centralised mock time server. Later it can be swapped out for a regular clock that just returns the current time.
	// On the clock advancing it calls the scheduler's OnTimeAdvance method which triggers any jobs
	_ = shared.NewMockClock(clockHost, []shared.OnTimeAdvanceCallback{sched})

	r := gin.Default()
	r.Use(cors.Default())

	server := handlers.NewServer(repos)
	server.RegisterRoutes(r)

	port := os.Getenv("API_PORT")
	log.Printf("Starting customer billing service on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
