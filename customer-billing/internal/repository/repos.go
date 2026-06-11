package repository

import "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db"

type Repos struct {
}

func New(database *db.DB) *Repos {
	return &Repos{}
}
