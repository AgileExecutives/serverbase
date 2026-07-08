package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
)

type InMemoryCustomerRepo struct {
	mu      sync.RWMutex
	byID    map[uint]*models.Customer
	byEmail map[string]*models.Customer
	nextID  uint
}

func NewInMemoryCustomerRepo() *InMemoryCustomerRepo {
	return &InMemoryCustomerRepo{byID: make(map[uint]*models.Customer), byEmail: make(map[string]*models.Customer), nextID: 1}
}

func (r *InMemoryCustomerRepo) FindByID(ctx context.Context, id uint) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if c, ok := r.byID[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found")
}

func (r *InMemoryCustomerRepo) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if c, ok := r.byEmail[email]; ok {
		return c, nil
	}
	return nil, nil
}

func (r *InMemoryCustomerRepo) Save(ctx context.Context, c *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.nextID
		r.nextID++
	}
	copy := *c
	r.byID[c.ID] = &copy
	if c.Email != "" {
		r.byEmail[c.Email] = &copy
	}
	return nil
}

var _ CustomerRepo = (*InMemoryCustomerRepo)(nil)

func (r *InMemoryCustomerRepo) FindByTenant(ctx context.Context, tenantID uint) ([]models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]models.Customer, 0)
	for _, c := range r.byID {
		if c.TenantID == tenantID {
			res = append(res, *c)
		}
	}
	return res, nil
}
