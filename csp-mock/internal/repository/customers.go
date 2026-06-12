package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/db"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	List(ctx context.Context) ([]db.Customer, error)
	GetByID(ctx context.Context, id string) (*db.Customer, error)
	Create(ctx context.Context, name, email string) (*db.Customer, error)
}

type customerRepo struct {
	db *db.DB
}

func newCustomerRepo(database *db.DB) CustomerRepository {
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

func (r *customerRepo) List(ctx context.Context) ([]db.Customer, error) {
	var customers []db.Customer
	if err := r.db.WithContext(ctx).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}
