package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

// GormTenantRepo is a GORM-backed implementation of TenantRepo.
type GormTenantRepo struct {
	db *gorm.DB
}

func NewGormTenantRepo(db *gorm.DB) *GormTenantRepo { return &GormTenantRepo{db: db} }

func (r *GormTenantRepo) FindByID(ctx context.Context, id uint) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.WithContext(ctx).First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormTenantRepo) FindBySlug(ctx context.Context, slug string) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&t).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *GormTenantRepo) Save(ctx context.Context, t *models.Tenant) error {
	return r.db.WithContext(ctx).Save(t).Error
}

func (r *GormTenantRepo) List(ctx context.Context) ([]models.Tenant, error) {
	var tenants []models.Tenant
	if err := r.db.WithContext(ctx).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

var _ TenantRepo = (*GormTenantRepo)(nil)
