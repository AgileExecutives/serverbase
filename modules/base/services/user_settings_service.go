package services

import (
	"errors"

	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

// UserSettingsService provides operations for user settings
type UserSettingsService struct {
	db *gorm.DB
}

// NewUserSettingsService creates a new UserSettingsService
func NewUserSettingsService(db *gorm.DB) *UserSettingsService {
	return &UserSettingsService{db: db}
}

// GetOrCreate returns settings for a user or creates defaults
func (s *UserSettingsService) GetOrCreate(userID uint) (models.UserSettings, error) {
	var settings models.UserSettings
	if err := s.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			settings = models.UserSettings{
				UserID:   userID,
				Language: "en",
				Timezone: "UTC",
				Theme:    "light",
				Settings: "{}",
			}
			if err := s.db.Create(&settings).Error; err != nil {
				return models.UserSettings{}, err
			}
			return settings, nil
		}
		return models.UserSettings{}, err
	}
	return settings, nil
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

	if err := s.db.Save(&settings).Error; err != nil {
		return models.UserSettings{}, err
	}
	return settings, nil
}

// Reset resets settings to defaults for the given user
func (s *UserSettingsService) Reset(userID uint) (models.UserSettings, error) {
	var settings models.UserSettings
	if err := s.db.Where("user_id = ?", userID).First(&settings).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserSettings{}, err
		}
		// create defaults
		settings = models.UserSettings{
			UserID:   userID,
			Language: "en",
			Timezone: "UTC",
			Theme:    "light",
			Settings: "{}",
		}
		if err := s.db.Create(&settings).Error; err != nil {
			return models.UserSettings{}, err
		}
		return settings, nil
	}

	settings.Language = "en"
	settings.Timezone = "UTC"
	settings.Theme = "light"
	settings.Settings = "{}"

	if err := s.db.Save(&settings).Error; err != nil {
		return models.UserSettings{}, err
	}
	return settings, nil
}
