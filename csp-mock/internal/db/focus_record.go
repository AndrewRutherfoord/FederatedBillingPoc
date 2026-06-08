package db

import (
	sharedmodels "github.com/andrewrutherfoord/fed-bill-poc/shared-models"
	"gorm.io/gorm"
)

// FocusRecord is the GORM model for a persisted FOCUS 1.3 line item.
type FocusRecord struct {
	gorm.Model
	sharedmodels.FocusLineItem `gorm:"embedded"`
}
