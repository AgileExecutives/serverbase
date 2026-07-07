package testutils

import "fmt"

// User is a tiny test user struct used by memory repos.
type User struct {
	ID    uint
	Email string
	Name  string
}

// MemoryUserRepo is a simple in-memory repository for tests.
type MemoryUserRepo struct {
	data map[uint]*User
}

// NewMemoryUserRepo creates a new in-memory user repo.
func NewMemoryUserRepo() *MemoryUserRepo {
	return &MemoryUserRepo{data: make(map[uint]*User)}
}

// Save stores or updates a user.
func (r *MemoryUserRepo) Save(u *User) error {
	if u == nil {
		return fmt.Errorf("nil user")
	}
	r.data[u.ID] = u
	return nil
}

// FindByID returns a user or an error if not found.
func (r *MemoryUserRepo) FindByID(id uint) (*User, error) {
	if u, ok := r.data[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("not found")
}
