package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

type InMemoryOrganizationRepo struct {
	mu            sync.Mutex
	nextID        uint
	organizations map[uint]models.Organization
}

func NewInMemoryOrganizationRepo() *InMemoryOrganizationRepo {
	return &InMemoryOrganizationRepo{nextID: 1, organizations: make(map[uint]models.Organization)}
}

func (r *InMemoryOrganizationRepo) Create(ctx context.Context, o *models.Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = r.nextID
	r.nextID++
	r.organizations[o.ID] = *o
	return nil
}

func (r *InMemoryOrganizationRepo) GetByID(ctx context.Context, id, tenantID uint) (*models.Organization, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.organizations[id]
	if !ok || o.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	return &o, nil
}

func (r *InMemoryOrganizationRepo) ListByTenant(ctx context.Context, offset, limit int, tenantID uint) ([]models.Organization, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []models.Organization
	for _, o := range r.organizations {
		if o.TenantID == tenantID {
			list = append(list, o)
		}
	}
	total := int64(len(list))
	start := offset
	end := offset + limit
	if start > len(list) {
		return []models.Organization{}, total, nil
	}
	if end > len(list) {
		end = len(list)
	}
	return list[start:end], total, nil
}

func (r *InMemoryOrganizationRepo) Update(ctx context.Context, o *models.Organization) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.organizations[o.ID]
	if !ok {
		return errors.New("not found")
	}
	r.organizations[o.ID] = *o
	return nil
}

func (r *InMemoryOrganizationRepo) Delete(ctx context.Context, id, tenantID uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	o, ok := r.organizations[id]
	if !ok || o.TenantID != tenantID {
		return errors.New("not found")
	}
	delete(r.organizations, id)
	return nil
}
