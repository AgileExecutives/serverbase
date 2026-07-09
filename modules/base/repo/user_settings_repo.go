package repo

import "github.com/AgileExecutives/serverbase/internal/models"

// UserSettingsRepo defines persistence operations for user settings
type UserSettingsRepo interface {
	FindByUserID(userID uint) (models.UserSettings, error)
	Create(settings *models.UserSettings) error
	Save(settings *models.UserSettings) error
}
