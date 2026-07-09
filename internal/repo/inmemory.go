package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

type InMemoryCustomerRepo struct {
	mu        sync.Mutex
	nextID    uint
	customers map[uint]models.Customer
}

func NewInMemoryCustomerRepo() *InMemoryCustomerRepo {
	return &InMemoryCustomerRepo{nextID: 1, customers: make(map[uint]models.Customer)}
}

func (r *InMemoryCustomerRepo) ListByTenant(ctx context.Context, tenantID uint, offset, limit int, active *bool) ([]models.Customer, int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []models.Customer
	for _, c := range r.customers {
		if c.TenantID != tenantID {
			continue
		}
		if active != nil && c.Active != *active {
			continue
		}
		list = append(list, c)
	}
	total := int64(len(list))
	// simple pagination
	start := offset
	end := offset + limit
	if start > len(list) {
		return []models.Customer{}, total, nil
	}
	if end > len(list) {
		end = len(list)
	}
	return list[start:end], total, nil
}

func (r *InMemoryCustomerRepo) GetByID(ctx context.Context, id, tenantID uint) (*models.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.customers[id]
	if !ok || c.TenantID != tenantID {
		return nil, errors.New("not found")
	}
	return &c, nil
}

func (r *InMemoryCustomerRepo) Create(ctx context.Context, c *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c.ID = r.nextID
	r.nextID++
	r.customers[c.ID] = *c
	return nil
}

func (r *InMemoryCustomerRepo) Update(ctx context.Context, c *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.customers[c.ID]
	if !ok {
		return errors.New("not found")
	}
	r.customers[c.ID] = *c
	return nil
}

func (r *InMemoryCustomerRepo) Delete(ctx context.Context, c *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.customers, c.ID)
	return nil
}

func (r *InMemoryCustomerRepo) PlanExists(ctx context.Context, planID uint) (bool, error) {
	// Simplified: assume any non-zero plan exists in tests
	return planID != 0, nil
}
