package repo

import (
	"context"
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
)

func TestInMemoryTenantRepo_SaveAndFind(t *testing.T) {
	r := NewInMemoryTenantRepo()
	ctx := context.Background()

	tnt := &models.Tenant{Name: "Acme", Slug: "acme"}
	if err := r.Save(ctx, tnt); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if tnt.ID == 0 {
		t.Fatalf("expected id assigned")
	}

	got, err := r.FindBySlug(ctx, "acme")
	if err != nil {
		t.Fatalf("find by slug error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected tenant, got nil")
	}

	found, err := r.FindByID(ctx, tnt.ID)
	if err != nil {
		t.Fatalf("find by id error: %v", err)
	}
	if found == nil {
		t.Fatalf("expected tenant by id, got nil")
	}
}
