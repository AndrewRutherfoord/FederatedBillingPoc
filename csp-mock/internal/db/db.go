package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	if err := database.AutoMigrate(&FocusRecord{}, &BillingAccount{}, &Customer{}, &Resource{}, &KeyValue{}); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}
	return database, nil
}
