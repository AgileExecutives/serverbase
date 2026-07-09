package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

type GormCustomerRepo struct{ db *gorm.DB }

func NewGormCustomerRepo(db *gorm.DB) *GormCustomerRepo { return &GormCustomerRepo{db: db} }

func (r *GormCustomerRepo) ListByTenant(ctx context.Context, tenantID uint, offset, limit int, active *bool) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Customer{}).Where("tenant_id = ?", tenantID)
	if active != nil {
		query = query.Where("active = ?", *active)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *GormCustomerRepo) GetByID(ctx context.Context, id, tenantID uint) (*models.Customer, error) {
	var c models.Customer
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormCustomerRepo) Create(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *GormCustomerRepo) Update(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *GormCustomerRepo) Delete(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Delete(c).Error
}

func (r *GormCustomerRepo) PlanExists(ctx context.Context, planID uint) (bool, error) {
	var p models.Plan
	if err := r.db.WithContext(ctx).First(&p, planID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
