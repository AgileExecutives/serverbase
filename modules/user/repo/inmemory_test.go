package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
)

func TestInMemoryUserRepo_SaveAndFind(t *testing.T) {
	r := NewInMemoryUserRepo()
	ctx := context.Background()

	user := &models.User{Email: "u1@example.com", FirstName: "A", LastName: "B"}
	if err := r.Save(ctx, user); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if user.ID == 0 {
		t.Fatalf("expected id assigned")
	}

	got, err := r.FindByEmail(ctx, "u1@example.com")
	if err != nil {
		t.Fatalf("find by email error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected user, got nil")
	}

	found, err := r.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("find by id error: %v", err)
	}
	if found == nil {
		t.Fatalf("expected user by id, got nil")
	}
}
