package db

import (
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared-models"
	"gorm.io/gorm"
)

// FocusRecord is the GORM model for a persisted FOCUS 1.3 line item.
type FocusRecord struct {
	gorm.Model
	ID                         string `gorm:"uniqueIndex" json:"id"`
	sharedmodels.FocusLineItem `gorm:"embedded"`
}
