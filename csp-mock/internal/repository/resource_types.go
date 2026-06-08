package repository

import (
	"context"
	"fmt"

	"github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"
)

type ResourceTypeRepository interface {
	List(ctx context.Context) []config.ResourceType
	GetById(ctx context.Context, id string) (*config.ResourceType, error)
}

type resourceTypeRepo struct {
	items []config.ResourceType
}

func newResourceTypeRepo(items []config.ResourceType) ResourceTypeRepository {
	return &resourceTypeRepo{items: items}
}

func (r *resourceTypeRepo) List(ctx context.Context) []config.ResourceType {
	return r.items
}

func (r *resourceTypeRepo) GetById(ctx context.Context, id string) (*config.ResourceType, error) {
	for i := range r.items {
		if r.items[i].ID == id {
			return &r.items[i], nil
		}
	}
	return nil, fmt.Errorf("resource type not found")
}
