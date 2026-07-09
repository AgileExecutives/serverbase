package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
)

type GormEmailRepo struct{ db *gorm.DB }

func NewGormEmailRepo(db *gorm.DB) *GormEmailRepo { return &GormEmailRepo{db: db} }

func (r *GormEmailRepo) List(ctx context.Context, offset, limit int, status string) ([]models.Email, int64, error) {
	var emails []models.Email
	var total int64
	q := r.db.Model(&models.Email{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&emails).Error; err != nil {
		return nil, 0, err
	}
	return emails, total, nil
}

func (r *GormEmailRepo) FindByID(ctx context.Context, id uint) (*models.Email, error) {
	var e models.Email
	if err := r.db.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *GormEmailRepo) Create(ctx context.Context, e *models.Email) error {
	return r.db.Create(e).Error
}

func (r *GormEmailRepo) UpdateStatus(ctx context.Context, id uint, status, errorMessage string) error {
	return r.db.Model(&models.Email{}).Where("id = ?", id).Updates(models.Email{Status: status, ErrorMessage: errorMessage}).Error
}

func (r *GormEmailRepo) Stats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)
	var total int64
	if err := r.db.Model(&models.Email{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total
	var count int64
	statuses := []string{"pending", "sent", "delivered", "failed"}
	for _, s := range statuses {
		if err := r.db.Model(&models.Email{}).Where("status = ?", s).Count(&count).Error; err != nil {
			return nil, err
		}
		stats[s] = count
	}
	return stats, nil
}
