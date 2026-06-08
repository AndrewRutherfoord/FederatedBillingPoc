package db

type Customer struct {
	ID    string `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"uniqueIndex;not null"`
}
