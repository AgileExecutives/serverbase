package repo

import (
	"context"
	"errors"
	"sync"

	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
)

type InMemoryContactRepo struct {
	mu          sync.RWMutex
	byID        map[uint]*models.Contact
	nextID      uint
	newsletters map[string]*basemodels.Newsletter
}

func NewInMemoryContactRepo() *InMemoryContactRepo {
	return &InMemoryContactRepo{byID: make(map[uint]*models.Contact), nextID: 1, newsletters: make(map[string]*basemodels.Newsletter)}
}

func (r *InMemoryContactRepo) ListContacts(ctx context.Context, offset, limit int, active *bool, contactType string) ([]models.Contact, int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []models.Contact
	for _, v := range r.byID {
		if active != nil && v.Active != *active {
			continue
		}
		if contactType != "" && v.Type != contactType {
			continue
		}
		out = append(out, *v)
	}
	total := int64(len(out))
	// simple offset/limit
	start := offset
	if start > len(out) {
		return []models.Contact{}, total, nil
	}
	end := start + limit
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], total, nil
}

func (r *InMemoryContactRepo) FindByID(ctx context.Context, id uint) (*models.Contact, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if v, ok := r.byID[id]; ok {
		copy := *v
		return &copy, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *InMemoryContactRepo) CreateContact(ctx context.Context, c *models.Contact) error {
	if c == nil {
		return errors.New("nil contact")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if c.ID == 0 {
		c.ID = r.nextID
		r.nextID++
	}
	copy := *c
	r.byID[c.ID] = &copy
	return nil
}

func (r *InMemoryContactRepo) UpdateContact(ctx context.Context, c *models.Contact) error {
	if c == nil {
		return errors.New("nil contact")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[c.ID]; !ok {
		return errors.New("not found")
	}
	copy := *c
	r.byID[c.ID] = &copy
	return nil
}

func (r *InMemoryContactRepo) DeleteContact(ctx context.Context, c *models.Contact) error {
	if c == nil {
		return errors.New("nil contact")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.byID, c.ID)
	return nil
}

func (r *InMemoryContactRepo) UpsertNewsletter(ctx context.Context, n *basemodels.Newsletter) (bool, error) {
	if n == nil {
		return false, errors.New("nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *n
	r.newsletters[n.Email] = &copy
	return true, nil
}

func (r *InMemoryContactRepo) ListNewsletters(ctx context.Context) ([]basemodels.Newsletter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]basemodels.Newsletter, 0, len(r.newsletters))
	for _, v := range r.newsletters {
		out = append(out, *v)
	}
	return out, nil
}

func (r *InMemoryContactRepo) DeleteNewsletterByEmail(ctx context.Context, email string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.newsletters[email]; !ok {
		return 0, nil
	}
	delete(r.newsletters, email)
	return 1, nil
}

var _ ContactRepo = (*InMemoryContactRepo)(nil)
