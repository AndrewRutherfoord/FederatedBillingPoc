package db

type KeyValue struct {
	Key   string `gorm:"primaryKey"`
	Value string `gorm:"not null"`
}
