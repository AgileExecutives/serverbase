package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
)

func TestInMemoryCustomerRepo_SaveAndFind(t *testing.T) {
	r := NewInMemoryCustomerRepo()
	ctx := context.Background()

	c := &models.Customer{Email: "cust@example.com", Name: "Cust Inc"}
	if err := r.Save(ctx, c); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if c.ID == 0 {
		t.Fatalf("expected id assigned")
	}

	got, err := r.FindByEmail(ctx, "cust@example.com")
	if err != nil {
		t.Fatalf("find by email error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected customer, got nil")
	}

	found, err := r.FindByID(ctx, c.ID)
	if err != nil {
		t.Fatalf("find by id error: %v", err)
	}
	if found == nil {
		t.Fatalf("expected customer by id, got nil")
	}
}
