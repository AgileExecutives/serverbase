package services

import (
	"testing"

	"github.com/AgileExecutives/serverbase/internal/models"
	baseRepo "github.com/AgileExecutives/serverbase/modules/base/repo"
)

func TestUserSettings_GetOrCreate_Update_Reset(t *testing.T) {
	r := baseRepo.NewInMemoryUserSettingsRepo()
	svc := NewUserSettingsService(r)

	// GetOrCreate should create defaults
	settings, err := svc.GetOrCreate(42)
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}
	if settings.UserID != 42 || settings.Language != "en" || settings.Theme != "light" {
		t.Fatalf("unexpected defaults: %+v", settings)
	}

	// Update some fields
	req := models.UserSettingsUpdateRequest{Language: "de", Theme: "dark", Settings: `{"a":1}`}
	updated, err := svc.Update(42, req)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Language != "de" || updated.Theme != "dark" {
		t.Fatalf("Update did not persist: %+v", updated)
	}

	// Reset should restore defaults
	reset, err := svc.Reset(42)
	if err != nil {
		t.Fatalf("Reset failed: %v", err)
	}
	if reset.Language != "en" || reset.Theme != "light" || reset.Settings != "{}" {
		t.Fatalf("Reset did not restore defaults: %+v", reset)
	}
}
