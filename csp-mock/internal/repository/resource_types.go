package repository

import "github.com/andrewrutherfoord/fed-bill-poc/csp-mock/internal/config"

type ResourceTypeRepository interface {
	List() []config.ResourceType
	Get(id string) (*config.ResourceType, bool)
}

type resourceTypeRepo struct {
	items []config.ResourceType
}

func newResourceTypeRepo(items []config.ResourceType) ResourceTypeRepository {
	return &resourceTypeRepo{items: items}
}

func (r *resourceTypeRepo) List() []config.ResourceType {
	return r.items
}

func (r *resourceTypeRepo) Get(id string) (*config.ResourceType, bool) {
	for i := range r.items {
		if r.items[i].ID == id {
			return &r.items[i], true
		}
	}
	return nil, false
}
