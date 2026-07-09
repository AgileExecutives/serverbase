package repo

import (
	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

type gormUserSettingsRepo struct{ db *gorm.DB }

func NewGormUserSettingsRepo(db *gorm.DB) UserSettingsRepo {
	return &gormUserSettingsRepo{db: db}
}

func (r *gormUserSettingsRepo) FindByUserID(userID uint) (models.UserSettings, error) {
	var s models.UserSettings
	if err := r.db.Where("user_id = ?", userID).First(&s).Error; err != nil {
		return models.UserSettings{}, err
	}
	return s, nil
}

func (r *gormUserSettingsRepo) Create(settings *models.UserSettings) error {
	return r.db.Create(settings).Error
}

func (r *gormUserSettingsRepo) Save(settings *models.UserSettings) error {
	return r.db.Save(settings).Error
}
