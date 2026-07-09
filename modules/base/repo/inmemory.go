package repo

import (
	"errors"
	"sync"

	"github.com/AgileExecutives/serverbase/internal/models"
)

type inMemoryUserSettingsRepo struct {
	mu   sync.RWMutex
	data map[uint]models.UserSettings
}

func NewInMemoryUserSettingsRepo() UserSettingsRepo {
	return &inMemoryUserSettingsRepo{data: make(map[uint]models.UserSettings)}
}

func (r *inMemoryUserSettingsRepo) FindByUserID(userID uint) (models.UserSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if s, ok := r.data[userID]; ok {
		return s, nil
	}
	return models.UserSettings{}, errors.New("record not found")
}

func (r *inMemoryUserSettingsRepo) Create(settings *models.UserSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[settings.UserID] = *settings
	return nil
}

func (r *inMemoryUserSettingsRepo) Save(settings *models.UserSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[settings.UserID] = *settings
	return nil
}
