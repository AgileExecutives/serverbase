package repo

import (
	"context"

	"github.com/AgileExecutives/serverbase/internal/models"
	"gorm.io/gorm"
)

type GormUserRepo struct {
	db *gorm.DB
}

func NewGormUserRepo(db *gorm.DB) *GormUserRepo { return &GormUserRepo{db: db} }

func (r *GormUserRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) FindByUsernameOrEmail(ctx context.Context, identifier string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ? OR email = ?", identifier, identifier).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepo) Save(ctx context.Context, u *models.User) error {
	return r.db.Save(u).Error
}

func (r *GormUserRepo) SaveNewsletter(ctx context.Context, n *models.Newsletter) error {
	return r.db.Create(n).Error
}

func (r *GormUserRepo) SaveTokenBlacklist(ctx context.Context, tb *models.TokenBlacklist) error {
	return r.db.Create(tb).Error
}

// Ensure interface compliance
var _ UserRepo = (*GormUserRepo)(nil)
var _ NewsletterRepo = (*GormUserRepo)(nil)
var _ TokenBlacklistRepo = (*GormUserRepo)(nil)
