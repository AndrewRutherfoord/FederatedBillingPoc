package db

import (
	"time"
)

type BillingPeriod struct {
	ID          string `gorm:"column:id;primaryKey"`
	PeriodStart time.Time
	PeriodEnd   time.Time
	Status      string
	OpenedAt    time.Time
	ClosedAt    time.Time
}
