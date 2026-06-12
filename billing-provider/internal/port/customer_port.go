package port

import (
	"github.com/andrewrutherfoord/fed-bill-poc/billing-provider/internal/repository"
)

// CustomerPortImpl contains the application logic for the customer-facing boundary.
type CustomerPortImpl struct {
	repositories *repository.Repos
}

func NewCustomerPort(repositories *repository.Repos) *CustomerPortImpl {
	return &CustomerPortImpl{repositories: repositories}
}
