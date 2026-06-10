package db

import (
	"time"
)

type BillingAccount struct {
	ID        string `gorm:"column:account_id;primaryKey"`
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
