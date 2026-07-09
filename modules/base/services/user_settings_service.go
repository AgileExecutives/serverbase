package services

import (
	"errors"

	"github.com/AgileExecutives/serverbase/internal/models"
	"github.com/AgileExecutives/serverbase/modules/base/repo"
)

// UserSettingsService provides operations for user settings
type UserSettingsService struct {
	repo repo.UserSettingsRepo
}

// NewUserSettingsService creates a new UserSettingsService backed by the given repo
func NewUserSettingsService(r repo.UserSettingsRepo) *UserSettingsService {
	return &UserSettingsService{repo: r}
}

// GetOrCreate returns settings for a user or creates defaults
func (s *UserSettingsService) GetOrCreate(userID uint) (models.UserSettings, error) {
	settings, err := s.repo.FindByUserID(userID)
	if err == nil {
		return settings, nil
	}
	// create defaults
	defaults := models.UserSettings{
		UserID:   userID,
		Language: "en",
		Timezone: "UTC",
		Theme:    "light",
		Settings: "{}",
	}
	if err := s.repo.Create(&defaults); err != nil {
		return models.UserSettings{}, err
	}
	return defaults, nil
}

// Update updates or creates settings for a user based on the request
func (s *UserSettingsService) Update(userID uint, req models.UserSettingsUpdateRequest) (models.UserSettings, error) {
	settings, err := s.GetOrCreate(userID)
	if err != nil {
		return models.UserSettings{}, err
	}

	if req.Language != "" {
		settings.Language = req.Language
	}
	if req.Timezone != "" {
		settings.Timezone = req.Timezone
	}
	if req.Theme != "" {
		settings.Theme = req.Theme
	}
	if req.Settings != "" {
		settings.Settings = req.Settings
	}

	if err := s.repo.Save(&settings); err != nil {
		return models.UserSettings{}, err
	}
	return settings, nil
}

// Reset resets settings to defaults for the given user
func (s *UserSettingsService) Reset(userID uint) (models.UserSettings, error) {
	settings, err := s.repo.FindByUserID(userID)
	if err != nil {
		// not found — create defaults
		if !errors.Is(err, errors.New("record not found")) {
			// unknown error from repo; still attempt to create defaults
		}
		defaults := models.UserSettings{
			UserID:   userID,
			Language: "en",
			Timezone: "UTC",
			Theme:    "light",
			Settings: "{}",
		}
		if err := s.repo.Create(&defaults); err != nil {
			return models.UserSettings{}, err
		}
		return defaults, nil
	}

	settings.Language = "en"
	settings.Timezone = "UTC"
	settings.Theme = "light"
	settings.Settings = "{}"

	if err := s.repo.Save(&settings); err != nil {
		return models.UserSettings{}, err
	}
	return settings, nil
}
