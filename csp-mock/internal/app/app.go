package app

import (
	"log"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/scheduler"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
	sharedscheduler "github.com/andrewrutherfoord/fed-bill-poc/shared/scheduler"
)

type App struct {
	Config *config.CspConfig
	DB     *db.DB
	Repos  *repository.Repos
	Clock  shared.Clock
	Sched  *sharedscheduler.Scheduler
}

type Config struct {
	ConfigPath string
	DBPath     string
	ClockHost  string
}

func NewApp(cfg Config) *App {
	// Load config
	cspCfg, err := config.Load(cfg.ConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Open database
	dbPath := cfg.DBPath
	if dbPath == "" {
		dbPath = os.Getenv("DB_PATH")
		if dbPath == "" {
			dbPath = "csp-mock.sqlite"
		}
	}
	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Create repositories
	repos := repository.New(cspCfg, database)

	// Setup scheduler
	sched := sharedscheduler.NewWithPersistence(scheduler.NewSchedulerPersistence(repos.KeyValue))

	// Setup clock
	clockHost := cfg.ClockHost
	if clockHost == "" {
		clockHost = os.Getenv("MOCK_CLOCK_HOST")
		if clockHost == "" {
			clockHost = "http://localhost:9999"
		}
	}
	clock := shared.NewMockClock(clockHost, []shared.OnTimeAdvanceCallback{sched})

	return &App{
		Config: cspCfg,
		DB:     database,
		Repos:  repos,
		Clock:  clock,
		Sched:  sched,
	}
}
