package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

type InMemoryTenantRepo struct {
	mu     sync.RWMutex
	byID   map[uint]*models.Tenant
	bySlug map[string]*models.Tenant
	nextID uint
}

func NewInMemoryTenantRepo() *InMemoryTenantRepo {
	return &InMemoryTenantRepo{byID: make(map[uint]*models.Tenant), bySlug: make(map[string]*models.Tenant), nextID: 1}
}

func (r *InMemoryTenantRepo) FindByID(ctx context.Context, id uint) (*models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.byID[id]; ok {
		return t, nil
	}
	return nil, errors.New("not found")
}

func (r *InMemoryTenantRepo) FindBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.bySlug[slug]; ok {
		return t, nil
	}
	return nil, nil
}

func (r *InMemoryTenantRepo) Save(ctx context.Context, t *models.Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.ID == 0 {
		t.ID = r.nextID
		r.nextID++
	}
	copy := *t
	r.byID[t.ID] = &copy
	r.bySlug[t.Slug] = &copy
	return nil
}

var _ TenantRepo = (*InMemoryTenantRepo)(nil)

func (r *InMemoryTenantRepo) List(ctx context.Context) ([]models.Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]models.Tenant, 0, len(r.byID))
	for _, v := range r.byID {
		res = append(res, *v)
	}
	return res, nil
}
