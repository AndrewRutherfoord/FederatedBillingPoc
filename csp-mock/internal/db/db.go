package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func Open(path string) (*DB, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := database.AutoMigrate(&FocusRecord{}, &BillingAccount{}, &Customer{}, &Resource{}, &KeyValue{}, &CostBatch{}, &OnboardingSession{}); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return &DB{database}, nil
}

func (d *DB) Close() error {
	sql, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("getting underlying sql.DB: %w", err)
	}
	return sql.Close()
}
