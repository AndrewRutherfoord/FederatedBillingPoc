package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"gorm.io/gorm"
)

type KeyValueRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
	Delete(ctx context.Context, key string) error
}

type keyValueRepo struct {
	db *gorm.DB
}

func newKeyValueRepo(database *gorm.DB) KeyValueRepository {
	return &keyValueRepo{db: database}
}

func (r *keyValueRepo) Get(ctx context.Context, key string) (string, error) {
	var kv db.KeyValue
	if err := r.db.WithContext(ctx).First(&kv, "key = ?", key).Error; err != nil {
		return "", err
	}
	return kv.Value, nil
}

func (r *keyValueRepo) Set(ctx context.Context, key, value string) error {
	kv := db.KeyValue{
		Key:   key,
		Value: value,
	}
	return r.db.WithContext(ctx).Save(&kv).Error
}

func (r *keyValueRepo) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Delete(&db.KeyValue{}, "key = ?", key).Error
}
