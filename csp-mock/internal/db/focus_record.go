package db

import (
	shared "github.com/andrewrutherfoord/fed-bill-poc/shared"
	"gorm.io/gorm"
)

// FocusRecord is the GORM model for a persisted FOCUS 1.3 line item.
type FocusRecord struct {
	gorm.Model
	ID                   string `gorm:"uniqueIndex" json:"id"`
	shared.FocusLineItem `gorm:"embedded"`
	BatchID              string `gorm:"index"`
}
