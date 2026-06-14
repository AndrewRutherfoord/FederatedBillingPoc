package repository

import (
	"context"

	"github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/config"
)

type Customer struct {
	Name  string
	Email string
}

type CustomerRepository interface {
	GetCustomerDetails(ctx context.Context) (*Customer, error)
}

type customerRepo struct {
	config *config.Config
}

func newCustomerRepo(config *config.Config) CustomerRepository {
	return &customerRepo{config: config}
}

func (r *customerRepo) GetCustomerDetails(ctx context.Context) (*Customer, error) {
	return &Customer{
		Name:  r.config.Name,
		Email: r.config.Email,
	}, nil
}
