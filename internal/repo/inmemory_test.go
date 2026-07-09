package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
)

func TestInMemoryCustomerRepo_CRUD(t *testing.T) {
	r := NewInMemoryCustomerRepo()
	ctx := context.Background()

	c := models.Customer{
		Name:     "Acme",
		Email:    "acme@example.com",
		TenantID: 1,
		PlanID:   1,
		Active:   true,
	}

	if err := r.Create(ctx, &c); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if c.ID == 0 {
		t.Fatalf("expected ID to be set")
	}

	got, err := r.GetByID(ctx, c.ID, c.TenantID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Email != c.Email {
		t.Fatalf("unexpected email: %s", got.Email)
	}

	// Update
	got.Phone = "+123"
	if err := r.Update(ctx, got); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	// List
	list, total, err := r.ListByTenant(ctx, 1, 0, 10, nil)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total == 0 || len(list) == 0 {
		t.Fatalf("expected list to contain items")
	}

	// Delete
	if err := r.Delete(ctx, got); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Plan exists
	ok, err := r.PlanExists(ctx, 1)
	if err != nil {
		t.Fatalf("planexists failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected plan to exist")
	}
}
