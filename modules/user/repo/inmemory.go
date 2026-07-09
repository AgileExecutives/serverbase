package repo

import (
	"context"
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

// InMemoryUserRepo is a simple in-memory implementation of UserRepository
type InMemoryUserRepo struct {
	mu          sync.RWMutex
	byID        map[uint]*models.User
	byEmail     map[string]*models.User
	nextID      uint
	newsletters []*models.Newsletter
	blacklists  []*models.TokenBlacklist
	// When true, SaveNewsletter will deduplicate by email instead of appending
	dedupeNewsletters bool
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

func (r *InMemoryUserRepo) FindByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// check email map first
	if u, ok := r.byEmail[identifier]; ok {
		return u, nil
	}
	// fall back to scanning by username
	for _, u := range r.byID {
		if u != nil && u.Username == identifier {
			return u, nil
		}
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

func (r *InMemoryUserRepo) SaveNewsletter(ctx context.Context, n *models.Newsletter) error {
	if n == nil {
		return errors.New("nil newsletter")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.dedupeNewsletters {
		for i, ex := range r.newsletters {
			if ex != nil && ex.Email == n.Email {
				// replace existing entry
				copy := *n
				r.newsletters[i] = &copy
				return nil
			}
		}
		copy := *n
		r.newsletters = append(r.newsletters, &copy)
		return nil
	}

	copy := *n
	r.newsletters = append(r.newsletters, &copy)
	return nil
}

// SetDeduplicateNewsletters enables or disables de-duplication behavior for newsletters
func (r *InMemoryUserRepo) SetDeduplicateNewsletters(enable bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.dedupeNewsletters = enable
}

func (r *InMemoryUserRepo) SaveTokenBlacklist(ctx context.Context, tb *models.TokenBlacklist) error {
	if tb == nil {
		return errors.New("nil token blacklist")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blacklists = append(r.blacklists, tb)
	return nil
}

// ListNewsletters returns a copy of stored newsletters for inspection in tests
func (r *InMemoryUserRepo) ListNewsletters() []*models.Newsletter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*models.Newsletter, len(r.newsletters))
	copy(res, r.newsletters)
	return res
}

// ListBlacklists returns a copy of stored token blacklist entries for tests
func (r *InMemoryUserRepo) ListBlacklists() []*models.TokenBlacklist {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]*models.TokenBlacklist, len(r.blacklists))
	copy(res, r.blacklists)
	return res
}

// Ensure interface compliance
var _ UserRepo = (*InMemoryUserRepo)(nil)
var _ NewsletterRepo = (*InMemoryUserRepo)(nil)
var _ TokenBlacklistRepo = (*InMemoryUserRepo)(nil)
