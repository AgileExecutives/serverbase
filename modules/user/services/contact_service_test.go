package services

import (
	"context"
	"errors"
	"testing"

	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/modules/user/repo"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
)

func TestContactService_CRUD(t *testing.T) {
	ctx := context.Background()
	r := repo.NewInMemoryContactRepo()
	s := NewContactServiceWithRepo(r, nil)

	// Create
	c := models.Contact{FirstName: "John", LastName: "Smith", Email: "john@example.com", Active: true}
	if err := s.CreateContact(ctx, &c); err != nil {
		t.Fatalf("CreateContact failed: %v", err)
	}

	// Get
	got, err := s.GetContact(ctx, c.ID)
	if err != nil {
		t.Fatalf("GetContact failed: %v", err)
	}
	if got.Email != "john@example.com" {
		t.Fatalf("unexpected email: %s", got.Email)
	}

	// Update
	got.LastName = "Doe"
	if err := s.UpdateContact(ctx, got); err != nil {
		t.Fatalf("UpdateContact failed: %v", err)
	}
	got2, err := s.GetContact(ctx, got.ID)
	if err != nil {
		t.Fatalf("GetContact after update failed: %v", err)
	}
	if got2.LastName != "Doe" {
		t.Fatalf("update did not persist")
	}

	// List
	list, total, err := s.ListContacts(ctx, 0, 10, nil, "")
	if err != nil {
		t.Fatalf("ListContacts failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("unexpected list results: total=%d len=%d", total, len(list))
	}

	// Delete
	if err := s.DeleteContact(ctx, got2); err != nil {
		t.Fatalf("DeleteContact failed: %v", err)
	}
	if _, err := s.GetContact(ctx, got2.ID); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected not found after delete, got: %v", err)
	}
}

func TestContactService_Newsletter(t *testing.T) {
	ctx := context.Background()
	r := repo.NewInMemoryContactRepo()
	s := NewContactServiceWithRepo(r, nil)

	n := basemodels.Newsletter{Name: "Alice", Email: "alice@example.com", Interest: "news"}
	ok, err := s.UpsertNewsletter(ctx, &n)
	if err != nil {
		t.Fatalf("UpsertNewsletter failed: %v", err)
	}
	if !ok {
		t.Fatalf("UpsertNewsletter returned ok=false")
	}

	list, err := s.ListNewsletters(ctx)
	if err != nil {
		t.Fatalf("ListNewsletters failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 newsletter, got %d", len(list))
	}

	rows, err := s.DeleteNewsletterByEmail(ctx, "alice@example.com")
	if err != nil {
		t.Fatalf("DeleteNewsletterByEmail failed: %v", err)
	}
	if rows != 1 {
		t.Fatalf("expected 1 row deleted, got %d", rows)
	}
}
