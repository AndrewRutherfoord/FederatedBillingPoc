package app

import (
	"fmt"
	"log"
	"os"

	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/config"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/db"
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
	"github.com/andrewrutherfoord/fed-bill-poc/shared"
)

// App holds all shared infrastructure that every service needs.
type App struct {
	Config *config.Config
	DB     *db.DB
	Repos  *repository.Repos
	Clock  *shared.MockClock
}

func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	log.Printf("Loaded config: ProviderID=%s, ProviderName=%s", cfg.ProviderID, cfg.ProviderName)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "billing-provider.sqlite"
	}
	database, err := db.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	clock := shared.NewMockClock("http://localhost:9999", []shared.OnTimeAdvanceCallback{})

	return &App{
		Config: cfg,
		DB:     database,
		Repos:  repository.New(cfg, database),
		Clock:  clock,
	}, nil
}

func (a *App) Close() error {
	return a.DB.Close()
}
