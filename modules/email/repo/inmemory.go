package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// InMemoryEmailRepo is a lightweight in-memory EmailRepo for unit tests.
type InMemoryEmailRepo struct {
	mu     sync.RWMutex
	byID   map[uint]*models.Email
	list   []*models.Email
	nextID uint
}

func NewInMemoryEmailRepo() *InMemoryEmailRepo {
	return &InMemoryEmailRepo{byID: make(map[uint]*models.Email), list: make([]*models.Email, 0), nextID: 1}
}

func (r *InMemoryEmailRepo) List(ctx context.Context, offset, limit int, status string) ([]models.Email, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res []models.Email
	for _, e := range r.list {
		if status == "" || e.Status == status {
			res = append(res, *e)
		}
	}
	total := int64(len(res))
	// apply offset/limit
	start := offset
	if start > len(res) {
		return []models.Email{}, total, nil
	}
	end := start + limit
	if limit <= 0 || end > len(res) {
		end = len(res)
	}
	slice := res[start:end]
	return slice, total, nil
}

func (r *InMemoryEmailRepo) FindByID(ctx context.Context, id uint) (*models.Email, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if e, ok := r.byID[id]; ok {
		copy := *e
		return &copy, nil
	}
	return nil, errors.New("not found")
}

func (r *InMemoryEmailRepo) Create(ctx context.Context, e *models.Email) error {
	if e == nil {
		return errors.New("nil email")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if e.ID == 0 {
		e.ID = r.nextID
		r.nextID++
	}
	copy := *e
	r.byID[e.ID] = &copy
	// prepend to list to emulate ORDER BY created_at DESC
	r.list = append([]*models.Email{&copy}, r.list...)
	return nil
}

func (r *InMemoryEmailRepo) UpdateStatus(ctx context.Context, id uint, status, errorMessage string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.byID[id]; ok {
		e.Status = status
		e.ErrorMessage = errorMessage
		return nil
	}
	return errors.New("not found")
}

func (r *InMemoryEmailRepo) Stats(ctx context.Context) (map[string]int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stats := map[string]int64{"total": 0, "pending": 0, "sent": 0, "delivered": 0, "failed": 0}
	for _, e := range r.list {
		stats["total"]++
		stats[e.Status]++
	}
	return stats, nil
}

// ensure interface compliance
var _ EmailRepo = (*InMemoryEmailRepo)(nil)
