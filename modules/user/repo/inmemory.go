package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// InMemoryUserRepo is a simple in-memory implementation of UserRepository
type InMemoryUserRepo struct {
	mu      sync.RWMutex
	byID    map[uint]*models.User
	byEmail map[string]*models.User
	nextID  uint
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{byID: make(map[uint]*models.User), byEmail: make(map[string]*models.User), nextID: 1}
}

func (r *InMemoryUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.byID[id]; ok {
		return u, nil
	}
	return nil, errors.New("not found")
}

func (r *InMemoryUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (r *InMemoryUserRepo) Save(ctx context.Context, u *models.User) error {
	if u == nil {
		return errors.New("nil user")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if u.ID == 0 {
		u.ID = r.nextID
		r.nextID++
	}
	copy := *u
	r.byID[u.ID] = &copy
	r.byEmail[u.Email] = &copy
	return nil
}

// Ensure interface compliance
var _ UserRepo = (*InMemoryUserRepo)(nil)
