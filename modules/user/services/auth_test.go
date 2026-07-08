package services

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/pkg/testutils"
)

// FakeUserRepo implements the minimal interfaces.UserRepository for tests.
type FakeUserRepo struct {
	byEmail map[string]*models.User
}

func NewFakeUserRepo() *FakeUserRepo { return &FakeUserRepo{byEmail: map[string]*models.User{}} }

func (r *FakeUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	// simple fake: search in map
	for _, u := range r.byEmail {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (r *FakeUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, nil
}

func (r *FakeUserRepo) Save(ctx context.Context, u *models.User) error {
	if u == nil {
		return nil
	}
	r.byEmail[u.Email] = u
	return nil
}

func TestAuthService_FindByEmailAndSave(t *testing.T) {
	logger := testutils.NewMockLogger()
	repo := NewFakeUserRepo()
	svc := NewAuthServiceWithRepo(repo, logger)
	ctx := context.Background()

	// missing user returns nil
	u, err := svc.FindByEmail(ctx, "nope@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u != nil {
		t.Fatalf("expected nil for missing user, got %+v", u)
	}

	// save and then find
	user := &models.User{ID: 1, Email: "a@b.com", FirstName: "A", LastName: "B"}
	if err := svc.SaveUser(ctx, user); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := svc.FindByEmail(ctx, "a@b.com")
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if got == nil {
		t.Fatalf("expected user, got nil")
	}
	if got.Email != "a@b.com" {
		t.Fatalf("unexpected email: %s", got.Email)
	}
}
