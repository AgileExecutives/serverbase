package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

type GormOrganizationRepo struct{ db *gorm.DB }

func NewGormOrganizationRepo(db *gorm.DB) *GormOrganizationRepo { return &GormOrganizationRepo{db: db} }

func (r *GormOrganizationRepo) Create(ctx context.Context, o *models.Organization) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *GormOrganizationRepo) GetByID(ctx context.Context, id, tenantID uint) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&org).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *GormOrganizationRepo) ListByTenant(ctx context.Context, offset, limit int, tenantID uint) ([]models.Organization, int64, error) {
	var organizations []models.Organization
	var total int64
	if err := r.db.WithContext(ctx).Model(&models.Organization{}).Where("tenant_id = ?", tenantID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Offset(offset).Limit(limit).Find(&organizations).Error; err != nil {
		return nil, 0, err
	}
	return organizations, total, nil
}

func (r *GormOrganizationRepo) Update(ctx context.Context, o *models.Organization) error {
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *GormOrganizationRepo) Delete(ctx context.Context, id, tenantID uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&models.Organization{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
