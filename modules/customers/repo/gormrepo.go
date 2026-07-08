package repo

import (
	"context"

	"github.com/AgileExecutives/shared-modules/saas-base/models"
	"gorm.io/gorm"
)

type GormCustomerRepo struct{ db *gorm.DB }

func NewGormCustomerRepo(db *gorm.DB) *GormCustomerRepo { return &GormCustomerRepo{db: db} }

func (r *GormCustomerRepo) FindByID(ctx context.Context, id uint) (*models.Customer, error) {
	var c models.Customer
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormCustomerRepo) FindByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var c models.Customer
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *GormCustomerRepo) Save(ctx context.Context, c *models.Customer) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *GormCustomerRepo) FindByTenant(ctx context.Context, tenantID uint) ([]models.Customer, error) {
	var res []models.Customer
	if err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

var _ CustomerRepo = (*GormCustomerRepo)(nil)
