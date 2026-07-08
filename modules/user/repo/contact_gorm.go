package repo

import (
	"context"
	"time"

	basemodels "github.com/AgileExecutives/serverbase/modules/user/models"
	"github.com/AgileExecutives/serverbase/pkg/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormContactRepo struct{ db *gorm.DB }

func NewGormContactRepo(db *gorm.DB) *GormContactRepo { return &GormContactRepo{db: db} }

func (r *GormContactRepo) ListContacts(ctx context.Context, offset, limit int, active *bool, contactType string) ([]models.Contact, int64, error) {
	var contacts []models.Contact
	var total int64
	q := r.db.Model(&models.Contact{})
	if active != nil {
		q = q.Where("active = ?", *active)
	}
	if contactType != "" {
		q = q.Where("type = ?", contactType)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (r *GormContactRepo) FindByID(ctx context.Context, id uint) (*models.Contact, error) {
	var c models.Contact
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *GormContactRepo) CreateContact(ctx context.Context, c *models.Contact) error {
	return r.db.Create(c).Error
}

func (r *GormContactRepo) UpdateContact(ctx context.Context, c *models.Contact) error {
	return r.db.Save(c).Error
}

func (r *GormContactRepo) DeleteContact(ctx context.Context, c *models.Contact) error {
	return r.db.Delete(c).Error
}

func (r *GormContactRepo) UpsertNewsletter(ctx context.Context, n *basemodels.Newsletter) (bool, error) {
	n.LastContact = time.Now()
	res := r.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "email"}}, DoUpdates: clause.AssignmentColumns([]string{"name", "interest", "source", "last_contact"})}).Create(n)
	if res.Error != nil {
		// fallback behaviour: try to find and update existing
		var existing basemodels.Newsletter
		if r.db.Where("email = ?", n.Email).First(&existing).Error == nil {
			existing.Name = n.Name
			existing.Interest = n.Interest
			existing.Source = n.Source
			existing.LastContact = time.Now()
			if saveErr := r.db.Save(&existing).Error; saveErr == nil {
				return true, nil
			} else {
				return false, saveErr
			}
		}
		return false, res.Error
	}
	return true, nil
}

func (r *GormContactRepo) ListNewsletters(ctx context.Context) ([]basemodels.Newsletter, error) {
	var list []basemodels.Newsletter
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *GormContactRepo) DeleteNewsletterByEmail(ctx context.Context, email string) (int64, error) {
	res := r.db.Unscoped().Where("email = ?", email).Delete(&basemodels.Newsletter{})
	return res.RowsAffected, res.Error
}

// Ensure interface
var _ ContactRepo = (*GormContactRepo)(nil)
