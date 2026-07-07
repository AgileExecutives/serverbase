package testutils

import "testing"

func TestMemoryUserRepo_SaveAndFind(t *testing.T) {
	repo := NewMemoryUserRepo()
	u := &User{ID: 1, Email: "x@example.com", Name: "X"}
	if err := repo.Save(u); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	got, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if got.Email != "x@example.com" {
		t.Fatalf("unexpected email: %s", got.Email)
	}

	if _, err := repo.FindByID(2); err == nil {
		t.Fatalf("expected error for missing user")
	}
}
