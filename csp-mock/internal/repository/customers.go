package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, name, email string) (*db.Customer, error)
	GetByID(ctx context.Context, id string) (*db.Customer, error)
}

type customerRepo struct {
	db *gorm.DB
}

func newCustomerRepo(database *gorm.DB) CustomerRepository {
	return &customerRepo{db: database}
}

func (r *customerRepo) Create(ctx context.Context, name, email string) (*db.Customer, error) {
	customer := db.Customer{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
	}
	if err := r.db.WithContext(ctx).Create(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepo) GetByID(ctx context.Context, id string) (*db.Customer, error) {
	var customer db.Customer
	if err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}
